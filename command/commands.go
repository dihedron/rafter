package command

import (
	"github.com/dihedron/rafter/command/data"
	"github.com/dihedron/rafter/command/run"
)

type Base struct {
	Debug string `short:"D" long:"debug" description:"The debug level of the application." optional:"yes" choice:"off" choice:"debug" choice:"info" choice:"warn" choice:"error"`

	Logger string `short:"L" long:"logger" description:"The logger to use." optional:"yes" choice:"zap" choice:"console" choice:"hcl" choice:"file" choice:"warn" choice:"off"`
}

// Commands is the set of root command groups.
type Commands struct {
	Run run.Run `command:"run" alias:"r" description:"Run the cluster."`

	Data data.Data `command:"data" alias:"d" description:"Manage data in the cluster."`
	// Join Join `command:"join" alias:"j" description:"Join a node to the cluster."`

	// Leave Leave `command:"leave" alias:"l" description:"Leave a node to the cluster."`
}
