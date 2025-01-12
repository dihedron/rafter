package cluster

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/Jille/raft-grpc-leader-rpc/leaderhealth"
	transport "github.com/Jille/raft-grpc-transport"
	"github.com/Jille/raftadmin"
	"github.com/dihedron/rafter/distributed"
	proto "github.com/dihedron/rafter/distributed/proto"
	"github.com/dihedron/rafter/logging"
	"github.com/dihedron/rafter/logging/noop"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	RetainSnapshotCount = 2
	RaftTimeout         = 10 * time.Second
)

type Cluster struct {
	id        string
	directory string
	address   Address
	peers     []Peer
	bootstrap bool
	context   *distributed.Context
	raft      *raft.Raft
	transport *transport.Manager
	server    *grpc.Server
	logger    logging.Logger
}

func New(id string, context *distributed.Context, options ...Option) (*Cluster, error) {

	c := &Cluster{
		id:      id,
		peers:   []Peer{},
		logger:  &noop.Logger{},
		context: context,
	}
	for _, option := range options {
		option(c)
	}

	// initialise the Raft cluster
	if err := os.MkdirAll(c.directory, 0700); err != nil {
		c.logger.Error("error creating raft base directory a '%s': %v", c.directory, err)
		return nil, fmt.Errorf("error creating raft base directory '%s': %w", c.directory, err)
	}

	// create the snapshot store; this allows the Raft to truncate the log
	snapshots, err := raft.NewFileSnapshotStore(c.directory, RetainSnapshotCount, os.Stderr)
	if err != nil {
		c.logger.Error("error creating file snapshot store: %v", err)
		return nil, fmt.Errorf("error creating file snapshot store: %w", err)
	}

	// create the BoltDB instance for both log store and stable store
	boltDB, err := raftboltdb.NewBoltStore(filepath.Join(c.directory, "raft.db"))
	if err != nil {
		c.logger.Error("error creating BoltDB store: %v", err)
		return nil, fmt.Errorf("error creating new BoltDB store: %w", err)
	}

	c.transport = transport.New(raft.ServerAddress(c.address.String()), []grpc.DialOption{grpc.WithInsecure()})

	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(c.id)
	config.SnapshotThreshold = 64
	c.raft, err = raft.NewRaft(config, c.context, boltDB, boltDB, snapshots, c.transport.Transport())
	if err != nil {
		c.logger.Error("error creating new raft cluster: %v", err)
		return nil, fmt.Errorf("error creating new Raft cluster: %w", err)
	}

	if c.bootstrap {
		servers := []raft.Server{
			{
				ID:       raft.ServerID(c.id),
				Suffrage: raft.Voter,
				Address:  c.transport.Transport().LocalAddr(),
			},
		}
		if len(c.peers) > 0 {
			for _, peer := range c.peers {
				servers = append(servers, raft.Server{
					ID:      raft.ServerID(peer.ID),
					Address: raft.ServerAddress(peer.Address.String()),
				})
			}
		}
		cluster := raft.Configuration{
			Servers: servers,
		}

		f := c.raft.BootstrapCluster(cluster)
		if err := f.Error(); err != nil {
			// maybe it's only because the cluster was already bootstrapped??
			c.logger.Warn("error boostrapping raft clutser: %v", err)
			//return nil, fmt.Errorf("raft.Raft.BootstrapCluster: %v", err)
		}
	}

	return c, nil
}

func (c *Cluster) StartRPCServer() error {
	// check that we can listen on the given address
	socket, err := net.Listen("tcp", c.address.String())
	if err != nil {
		c.logger.Error("failed to listen: %v", err)
		return fmt.Errorf("failed to listen on '%s': %w", c.address.String(), err)
	}
	c.logger.Debug("TCP address %s available", c.address.String())
	// start the gRPC server
	c.server = grpc.NewServer()
	proto.RegisterContextServer(c.server, distributed.NewRPCInterface(c.context, c.raft, c.logger))
	c.transport.Register(c.server)
	leaderhealth.Setup(c.raft, c.server, []string{"quis.RaftLeader"})
	raftadmin.Register(c.server, c.raft)
	reflection.Register(c.server)

	c.logger.Info("starting gRPC server")

	go func() error {
		if err := c.server.Serve(socket); err != nil {
			c.logger.Error("failed to serve gRPC interface: %w")
			return fmt.Errorf("error starting gRPC interface: %w", err)
		}
		return nil
	}()
	return nil
}

func (c *Cluster) StopRPCServer() {
	c.logger.Info("stopping gRPC server")
	c.server.GracefulStop()
}

type NodeState uint8

const (
	Initial NodeState = iota
	Leader
	Follower
	Exiting
)

const (
	LeadershipPollInterval = time.Duration(500 * time.Millisecond)
)

func (c *Cluster) MonitorClusterEvents(ctx context.Context) <-chan NodeState {

	events := make(chan NodeState, 1)

	go func(ctx context.Context, events chan<- NodeState) {

		// start a ticker so that we're woken up every X milliseconds, regardless
		ticker := time.NewTicker(LeadershipPollInterval)
		c.logger.Debug("background checker started ticking every %+v ms", LeadershipPollInterval)
		defer ticker.Stop()

		// open the channel to get leader elections
		elections := c.raft.LeaderCh()

		// also register for cluster-related observations
		observations := make(chan raft.Observation, 1)
		observer := raft.NewObserver(observations, true, nil)
		c.raft.RegisterObserver(observer)

		state := Initial

		// at the very beginning, only the leader receives a ledership election
		// notification via the elections channel; the followers know nothing
		// about their state so they have to resort to checking the state from
		// the Raft cluster; starting up the cluster takes some time: from the
		// logs we see that while the leader knows it is the leader immediately
		// via the Raft.LeaderCh() channel, the followers only know that a
		// new leader is being elected via the observer events because they are
		// requested to vote for a candidate; after the leader has been elected, it
		// takes a while for the followers to get up to date: they have to apply
		// all the outstanding log entries to their current state before starting
		// to receive new entries and this usually takes a few seconds (depending
		// on how old the snapshot is); this initial loop is only needed to check
		// whether we're leaders of followers, therefore we'll be spending very
		// little time inside of it; after having bootstrapped the cluster, leader
		// elections, demotions and changes will flow into the leader and follower
		// loops and will be handled there
	loop:
		for {
			select {
			case <-ctx.Done():
				c.logger.Info("context has been cancelled, closing down")
				events <- Exiting
				break loop
			case elected := <-elections:
				c.logger.Info("cluster leadership changed (leader: %t)", elected)
				if state == Initial || (state == Follower && elected) || (state == Leader && !elected) {
					if elected && c.raft.State() == raft.Leader {
						c.logger.Info("I'm the new leader")
						state = Leader
						events <- Leader
					} else if c.raft.State() == raft.Follower {
						c.logger.Info("I'm a follower now")
						state = Follower
						events <- Follower
					}
				}
			case observation := <-observations:
				c.logger.Debug("received observation: %T", observation.Data)
				switch observation := observation.Data.(type) {
				case raft.PeerObservation:
					c.logger.Debug("received peer observation (id: %s, address: %s)", observation.Peer.ID, observation.Peer.Address)
				case raft.LeaderObservation:
					c.logger.Info("received leader observation (leader: %s)", observation.Leader)
				case raft.RequestVoteRequest:
					c.logger.Debug("received request vote request observation (leadership transfer: %t, term: %d)", observation.LeadershipTransfer, observation.Term)
				case raft.RaftState:
					c.logger.Debug("received raft state observation: %s", observation)
				default:
					c.logger.Warn("unhandled observation type: %T", observation)
				}
			// case interrupt := <-interrupts:
			// 	c.logger.Info("received interrupt: %d", interrupt)
			// 	events <- Exit
			// os.Exit(1)
			case <-ticker.C:
				switch c.raft.State() {
				case raft.Leader:
					if state != Leader {
						c.logger.Info("this node has become the new leader")
						state = Leader
						events <- Leader
					}
				case raft.Follower:
					if state != Follower {
						c.logger.Info("this node has become a follower")
						state = Follower
						events <- Follower
					}
				case raft.Candidate:
					c.logger.Debug("this node is a candidate")
				case raft.Shutdown:
					c.logger.Info("raft cluster is shutting down")
					events <- Exiting
				}
			}
		}
	}(ctx, events)
	return events
}
