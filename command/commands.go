package command

import (
	"github.com/dihedron/rafter/command/administration"
	"github.com/dihedron/rafter/command/data"
	"github.com/dihedron/rafter/command/run"
)

// Commands is the set of root command groups.
type Commands struct {
	Run run.Run `command:"run" alias:"r" description:"Run the cluster."`

	Data data.Data `command:"data" alias:"d" description:"Manage data in the cluster."`

	Administration administration.Administration `command:"administration" alias:"admin" alias:"a" description:"Run command against the cluster."`
}
