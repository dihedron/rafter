package command

/*
type Run struct {
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

	logger := logging.NewConsoleLogger(os.Stdout)

	fsm := cache.New()

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

	return nil
}
*/
