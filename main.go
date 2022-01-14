package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dihedron/rafter/command"
	"github.com/jessevdk/go-flags"
)

var (
	myAddr = flag.String("address", "localhost:50051", "TCP host+port for this node")
	raftId = flag.String("raft_id", "", "Node id used by Raft")

	raftDir       = flag.String("raft_data_dir", "data/", "Raft data dir")
	raftBootstrap = flag.Bool("raft_bootstrap", false, "Whether to bootstrap the Raft cluster")
)

func main() {
	options := command.Commands{}
	if _, err := flags.NewParser(&options, flags.Default).Parse(); err != nil {
		switch flagsErr := err.(type) {
		case flags.ErrorType:
			if flagsErr == flags.ErrHelp {
				os.Exit(0)
			}
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		default:
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}
}

func main2() {

	flag.Parse()
}
