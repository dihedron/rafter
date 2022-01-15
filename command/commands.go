package command

import (
	"github.com/dihedron/rafter/command/log"
	"github.com/dihedron/rafter/command/run"
)

// Commands is the set of root command groups.
type Commands struct {
	Run run.Run `command:"run" alias:"r" description:"Run the cluster."`

	Log log.Log `command:"log" alias:"l" description:"Handle values in the cluster distributed log."`
	// Join Join `command:"join" alias:"j" description:"Join a node to the cluster."`

	// Leave Leave `command:"leave" alias:"l" description:"Leave a node to the cluster."`
}
