package cluster

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	grpctx "github.com/Jille/raft-grpc-transport"
	"github.com/dihedron/rafter/logging"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	"google.golang.org/grpc"
)

const (
	retainSnapshotCount = 2
	raftTimeout         = 10 * time.Second
)

type Cluster struct {
	id        string
	directory string
	address   string
	peers     []Peer
	bootstrap bool
	raft      *raft.Raft
	logger    logging.Logger
}

// New creates a new Cluster, applying all the provided functional options.
func New(id string, fsm raft.FSM, options ...Option) (*Cluster, error) {
	c := &Cluster{
		id:     id,
		peers:  []Peer{},
		logger: &logging.NoOpLogger{},
	}
	for _, option := range options {
		option(c)
	}

	err := os.MkdirAll(c.directory, 0700)
	if err != nil {
		return nil, fmt.Errorf("error creating raft base directory '%s': %w", c.directory, err)
	}

	// create the snapshot store; this allows the Raft to truncate the log
	snapshots, err := raft.NewFileSnapshotStore(c.directory, retainSnapshotCount, os.Stderr)
	if err != nil {
		return nil, fmt.Errorf("error creating file snapshot store: %w", err)
	}

	// create the BoltDB instance for both log store and stable store
	boltDB, err := raftboltdb.NewBoltStore(filepath.Join(c.directory, "raft.db"))
	if err != nil {
		return nil, fmt.Errorf("error creating new Bolt store: %w", err)
	}

	// setup Raft configuration
	configuration := raft.DefaultConfig()
	configuration.LocalID = raft.ServerID(c.id)

	// // setup Raft communication
	// addr, err := net.ResolveTCPAddr("tcp", c.address)
	// if err != nil {
	// 	return nil, fmt.Errorf("error resolving bind address '%s': %w", c.address, err)
	// }
	// addr.String()
	// transport, err := raft.NewTCPTransport(c.address, addr, 3, 10*time.Second, os.Stderr)
	// if err != nil {
	// 	return nil, fmt.Errorf("error creating TCP transport: %w", err)
	// }

	///////////////////////

	tm := grpctx.New(raft.ServerAddress(c.address), []grpc.DialOption{grpc.WithInsecure()})
	s := grpc.NewServer()
	tm.Register(s)
	//r, err := raft.NewRaft(..., tm.Transport())

	//////////////////////

	// instantiate the Raft system
	//c.raft, err = raft.NewRaft(configuration, fsm, boltDB, boltDB, snapshots, transport)
	c.raft, err = raft.NewRaft(configuration, fsm, boltDB, boltDB, snapshots, tm.Transport())
	if err != nil {
		return nil, fmt.Errorf("error creating new raft cluster: %w", err)
	}

	c.logger.Debug("raft cluster created")

	if c.bootstrap {
		servers := []raft.Server{
			{
				ID: raft.ServerID(c.id),
				//Address: transport.LocalAddr(),
				Address: tm.Transport().LocalAddr(),
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
		if err := c.raft.BootstrapCluster(cluster).Error(); err != nil {
			// maybe already bootstrapped??
			c.logger.Warn("cluster already bootstrapped?: %v", err)
		}
	}
	return c, err
}

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
