package data

import (
	"github.com/dihedron/rafter/cluster"
	"github.com/dihedron/rafter/command/base"
)

type Base struct {
	base.Base

	Peers []cluster.Peer `short:"p" long:"peer" description:"The address of a peer node in the cluster to join" required:"yes"`
}

// Log is the set of distributed log related commands.
type Data struct {
	Set Set `command:"set" alias:"s" description:"Set a value in the distributed log."`

	Get Get `command:"get" alias:"g" description:"Get a value from a distributed log."`

	Benchmark Benchmark `command:"benchmark" alias:"b" description:"Benchmark the speed of the distributed log."`

	// Join Join `command:"join" alias:"j" description:"Join a node to the cluster."`

	// Leave Leave `command:"leave" alias:"l" description:"Leave a node to the cluster."`
}
