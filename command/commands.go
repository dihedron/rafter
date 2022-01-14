package command

// Commands is the set of root command groups.
type Commands struct {
	// Plugin is the namespace of all plugin related commands.
	//Run Run `command:"run" alias:"r" description:"Run the cluster."`

	Run2 Run2 `command:"run" alias:"r" description:"Run the cluster."`

	// Join Join `command:"join" alias:"j" description:"Join a node to the cluster."`

	// Leave Leave `command:"leave" alias:"l" description:"Leave a node to the cluster."`
}
