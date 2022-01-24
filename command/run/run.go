package run

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	transport "github.com/Jille/raft-grpc-transport"
	"github.com/dihedron/rafter/application"
	"github.com/dihedron/rafter/cluster"
	"github.com/dihedron/rafter/command/base"
	"github.com/dihedron/rafter/logging"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	"google.golang.org/grpc"
)

type Run struct {
	base.Base
	// Bootstrap starts the cluster in bootstrap mode.
	Bootstrap bool `short:"b" long:"bootstrap" description:"Whether to boostrap the cluster." optional:"yes"`
	// Address is the intra-cluster bind address for Raft communications.
	Address cluster.Address `short:"a" long:"address" description:"The network address for Raft and exposed services." optional:"yes" default:"localhost:7001"`
	// Join specified whether the node should join a cluster.
	Peers []cluster.Peer `short:"p" long:"peer" description:"The address of a peer node in the cluster to join" optional:"yes"`
	// State is the directory for Raft cluster state storage.
	Directory string `short:"d" long:"directory" description:"The base directory where Raft cluster state and snapshots are stored." optional:"yes" default:"./state"`
}

func (cmd *Run) Execute(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("no node id specified: (%v)", args)
	}
	fmt.Printf("starting a node at '%s' (base directory '%s'), with peers %+v\n", cmd.Address, cmd.Directory, cmd.Peers)

	logger := logging.NewConsoleLogger(logging.StdOut)
	// logger := logging.NewConsoleLogger(logging.StdOut)
	// logger := logging.NewLogLogger("rafter")

	defer cmd.ProfileCPU(logger).Close()

	fsm := application.New(logger)

	c, err := cluster.New(
		args[0],
		fsm,
		cluster.WithDirectory(cmd.Directory),
		cluster.WithNetAddress(cmd.Address.String()),
		cluster.WithPeers(cmd.Peers...),
		cluster.WithLogger(logger),
		cluster.WithBootstrap(cmd.Bootstrap),
	)
	if err != nil {
		return fmt.Errorf("error creating new cluster: %w", err)
	}
	c.Test()
	cmd.ProfileMemory(logger)
	return nil
}

func (cmd *Run) NewRaft(ctx context.Context, myID, myAddress string, fsm raft.FSM) (*raft.Raft, *transport.Manager, error) {
	c := raft.DefaultConfig()
	c.LocalID = raft.ServerID(myID)

	err := os.MkdirAll(cmd.Directory, 0700)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating raft base directory '%s': %w", cmd.Directory, err)
	}

	// create the snapshot store; this allows the Raft to truncate the log
	snapshots, err := raft.NewFileSnapshotStore(cmd.Directory, 10, os.Stderr)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating file snapshot store: %w", err)
	}

	// create the BoltDB instance for both log store and stable store
	boltDB, err := raftboltdb.NewBoltStore(filepath.Join(cmd.Directory, "raft.db"))
	if err != nil {
		return nil, nil, fmt.Errorf("error creating new Bolt store: %w", err)
	}

	tm := transport.New(raft.ServerAddress(myAddress), []grpc.DialOption{grpc.WithInsecure()})

	r, err := raft.NewRaft(c, fsm, boltDB, boltDB, snapshots, tm.Transport())
	if err != nil {
		return nil, nil, fmt.Errorf("raft.NewRaft: %v", err)
	}

	if cmd.Bootstrap {
		servers := []raft.Server{
			{
				ID:       raft.ServerID(myID),
				Suffrage: raft.Voter,
				Address:  tm.Transport().LocalAddr(),
			},
		}
		if len(cmd.Peers) > 0 {
			for _, peer := range cmd.Peers {
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
