package main

import (
	"fmt"
	"os"

	"github.com/dihedron/rafter/command"
	"github.com/dihedron/rafter/logging"
	"github.com/jessevdk/go-flags"
)

func main() {

	logger := logging.NewConsoleLogger(logging.StdOut)
	logger.Trace("hallo")
	logger.Debug("hallo")
	logger.Info("hallo")
	logger.Warn("hallo")
	logger.Error("hallo")

	return

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
