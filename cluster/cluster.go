package cluster

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/Jille/raft-grpc-leader-rpc/leaderhealth"
	transport "github.com/Jille/raft-grpc-transport"
	"github.com/Jille/raftadmin"
	"github.com/dihedron/rafter/application"
	"github.com/dihedron/rafter/logging"
	pb "github.com/dihedron/rafter/proto"
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
	raft      *raft.Raft
	transport *transport.Manager
	server    *grpc.Server
	logger    logging.Logger
}

func New(id string, fsm raft.FSM, options ...Option) (*Cluster, error) {

	c := &Cluster{
		id:     id,
		peers:  []Peer{},
		logger: &logging.NoOpLogger{},
	}
	for _, option := range options {
		option(c)
	}

	// check that we can listen on the given address
	socket, err := net.Listen("tcp", c.address.String())
	if err != nil {
		c.logger.Error("failed to listen: %v", err)
		return nil, fmt.Errorf("failed to listen on '%s': %w", c.address.String(), err)
	}

	c.logger.Info("TCP address %s available", c.address.String())

	cache := application.New(c.logger)

	// initialise the Raft cluster
	if err := os.MkdirAll(c.directory, 0700); err != nil {
		// TODO: logger.Error
		return nil, fmt.Errorf("error creating raft base directory '%s': %w", c.directory, err)
	}

	// create the snapshot store; this allows the Raft to truncate the log
	snapshots, err := raft.NewFileSnapshotStore(c.directory, RetainSnapshotCount, os.Stderr)
	if err != nil {
		// TODO: logger.Error
		return nil, fmt.Errorf("error creating file snapshot store: %w", err)
	}

	// create the BoltDB instance for both log store and stable store
	boltDB, err := raftboltdb.NewBoltStore(filepath.Join(c.directory, "raft.db"))
	if err != nil {
		// TODO: logger.Error
		return nil, fmt.Errorf("error creating new BoltDB store: %w", err)
	}

	c.transport = transport.New(raft.ServerAddress(c.address.String()), []grpc.DialOption{grpc.WithInsecure()})

	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(c.id)
	c.raft, err = raft.NewRaft(config, fsm, boltDB, boltDB, snapshots, c.transport.Transport())
	if err != nil {
		// TODO: logger.Error
		return nil, fmt.Errorf("error cereating new Raft cluster: %w", err)
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

	// ctx := context.Background()
	// r, tm, err := c.initRaft(ctx, id, c.address.String(), cache)
	// if err != nil {
	// 	log.Fatalf("failed to start raft: %v", err)
	// }
	c.server = grpc.NewServer()
	pb.RegisterLogServer(c.server, application.NewRPCInterface(cache, c.raft, c.logger))
	c.transport.Register(c.server)
	leaderhealth.Setup(c.raft, c.server, []string{"Log"})
	raftadmin.Register(c.server, c.raft)
	reflection.Register(c.server)
	if err := c.server.Serve(socket); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	return c, nil
}

/*
func (cmd *Cluster) initRaft(ctx context.Context, nodeId string, nodeAddress string, fsm raft.FSM) (*raft.Raft, *transport.Manager, error) {
	c := raft.DefaultConfig()
	c.LocalID = raft.ServerID(nodeId)

	// err := os.MkdirAll(cmd.directory, 0700)
	// if err != nil {
	// 	return nil, nil, fmt.Errorf("error creating raft base directory '%s': %w", cmd.directory, err)
	// }

	// create the snapshot store; this allows the Raft to truncate the log
	snapshots, err := raft.NewFileSnapshotStore(cmd.directory, 10, os.Stderr)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating file snapshot store: %w", err)
	}

	// create the BoltDB instance for both log store and stable store
	boltDB, err := raftboltdb.NewBoltStore(filepath.Join(cmd.directory, "raft.db"))
	if err != nil {
		return nil, nil, fmt.Errorf("error creating new Bolt store: %w", err)
	}

	tm := transport.New(raft.ServerAddress(nodeAddress), []grpc.DialOption{grpc.WithInsecure()})

	r, err := raft.NewRaft(c, fsm, boltDB, boltDB, snapshots, tm.Transport())
	if err != nil {
		return nil, nil, fmt.Errorf("raft.NewRaft: %v", err)
	}

	if cmd.bootstrap {
		servers := []raft.Server{
			{
				ID:       raft.ServerID(nodeId),
				Suffrage: raft.Voter,
				Address:  tm.Transport().LocalAddr(),
			},
		}
		if len(cmd.peers) > 0 {
			for _, peer := range cmd.peers {
				servers = append(servers, raft.Server{
					ID:      raft.ServerID(peer.ID),
					Address: raft.ServerAddress(peer.Address.String()),
				})
			}
		}
		cluster := raft.Configuration{
			Servers: servers,
		}

		f := r.BootstrapCluster(cluster)
		if err := f.Error(); err != nil {
			return nil, nil, fmt.Errorf("raft.Raft.BootstrapCluster: %v", err)
		}
	}

	return r, tm, nil
}
*/

const (
	LeadershipPollInterval = time.Duration(500 * time.Millisecond)
)

func (c *Cluster) Test() {
	// handle interrupts
	interrupts := make(chan os.Signal, 1)
	signal.Notify(interrupts, os.Interrupt, syscall.SIGTERM)
	defer close(interrupts)

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

	// leader := false

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
	// election_loop:
	for {
		select {
		case election := <-elections:
			c.logger.Info("cluster leadership changed (leader: %t)", election)
			// leader = election
			// break election_loop
		case observation := <-observations:
			c.logger.Debug("received observation: %T", observation.Data)
			switch observation := observation.Data.(type) {
			case raft.PeerObservation:
				c.logger.Debug("received peer observation (id: %s, address: %s)", observation.Peer.ID, observation.Peer.Address)
			case raft.LeaderObservation:
				c.logger.Debug("received leader observation (leader: %s)", observation.Leader)
			case raft.RequestVoteRequest:
				c.logger.Debug("received request vote request observation (leadership transfer: %t, term: %d)", observation.LeadershipTransfer, observation.Term)
			case raft.RaftState:
				c.logger.Debug("received raft state observation: %s", observation)
			default:
				c.logger.Warn("unhandled observation type: %T", observation)
			}
		case interrupt := <-interrupts:
			c.logger.Info("received interrupt: %d", interrupt)
			os.Exit(1)
		case <-ticker.C:
			switch c.raft.State() {
			case raft.Leader:
				c.logger.Info("this node is the leader")
				// leader = true
				// break election_loop
			case raft.Follower:
				c.logger.Info("this node is a follower")
				// leader = false
				// break election_loop
			case raft.Candidate:
				c.logger.Info("this node is a candidate")
			case raft.Shutdown:
				c.logger.Info("raft cluster is shut down")
			}
		}
	}
	// fmt.Printf("leader: %t\n", leader)

	// select {
	// case <-interrupts:
	// 	c.logger.Info("closing down...")
	// 	os.Exit(1)
	// }
}
