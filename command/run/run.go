package run

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dihedron/rafter/cluster"
	"github.com/dihedron/rafter/command/base"
	"github.com/dihedron/rafter/distributed"
	"github.com/dihedron/rafter/logging"
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

	appl := distributed.NewContext(logger)

	c, err := cluster.New(
		args[0],
		appl,
		cluster.WithDirectory(cmd.Directory),
		cluster.WithNetAddress(cmd.Address.String()),
		cluster.WithPeers(cmd.Peers...),
		cluster.WithLogger(logger),
		cluster.WithBootstrap(cmd.Bootstrap),
	)
	if err != nil {
		return fmt.Errorf("error creating new cluster: %w", err)
	}

	// start the gRPC server; it will be closed down
	// when we send an interrupt and exit the process
	c.StartRPCServer()

	interrupts, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	events := c.MonitorClusterEvents(interrupts)

	var (
		ctx    context.Context
		cancel context.CancelFunc
		done   = make(chan bool, 1)
	)
loop:
	for {
		select {
		case <-interrupts.Done():
			logger.Info("outer context cancelled by interrupts: closing down")
			if cancel != nil {
				cancel()
				cancel = nil
				select {
				case <-done:
					logger.Debug("existing routine exited")
				default:
					logger.Debug("no routine waiting to exit")
				}
			}
			// release the resources
			stop()
			// unless we exit, the process will stay up trying to
			// contact the other Raft cluster peers
			os.Exit(1)
			break loop
		case event := <-events:
			switch event {
			case cluster.Leader:
				logger.Info("received notification: I'm the leader")
				if cancel != nil {
					cancel()
					cancel = nil
					select {
					case <-done:
						logger.Debug("existing routine exited")
					}
				}
				logger.Info("starting the new leader routine")
				ctx, cancel = context.WithCancel(interrupts)
				go LeaderRoutine(ctx, logger, done)
			case cluster.Follower:
				logger.Info("received notification: I'm a follower")
				if cancel != nil {
					cancel()
					cancel = nil
					select {
					case <-done:
						logger.Debug("existing routine exited")
					}
				}
				logger.Info("starting the new follower routine")
				ctx, cancel = context.WithCancel(interrupts)
				go FollowerRoutine(ctx, logger, done)
			case cluster.Exiting:
				logger.Info("received notification: cluster is closing down")
				if cancel != nil {
					cancel()
					cancel = nil
					select {
					case <-done:
						logger.Debug("existing routine exited")
					default:
						logger.Debug("no routine waiting to exit")
					}
				}
				break loop
			}
		}
	}

	cmd.ProfileMemory(logger)
	return nil
}
