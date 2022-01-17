package command

import (
	"github.com/dihedron/rafter/command/data"
	"github.com/dihedron/rafter/command/run"
)

// Commands is the set of root command groups.
type Commands struct {
	Run run.Run `command:"run" alias:"r" description:"Run the cluster."`

	Data data.Data `command:"data" alias:"d" description:"Manage data in the cluster."`
	// Join Join `command:"join" alias:"j" description:"Join a node to the cluster."`

	// Leave Leave `command:"leave" alias:"l" description:"Leave a node to the cluster."`
}
