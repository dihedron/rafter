package cluster

import (
	"github.com/dihedron/rafter/logging"
)

// Option is the type for functional options.
type Option func(*Cluster)

func WithBootstrap(value bool) Option {
	return func(c *Cluster) {
		c.bootstrap = value
	}
}

// WithDirectory specifies the directory where the Raft cluster
// state is stored.
func WithDirectory(dir string) Option {
	return func(c *Cluster) {
		c.directory = dir // path.Join(dir, c.id)
	}
}

// WithNetAddress specifies the address used in multiplexing mode
// both for intra-cluster communications and to expose the gRPC
// services.
func WithNetAddress(address string) Option {
	return func(c *Cluster) {
		if address != "" {
			c.address = address
		}
	}
}

// WithPeer specifies a peer to contact to join the cluster.
func WithPeer(peer Peer) Option {
	return func(c *Cluster) {
		c.peers = append(c.peers, peer)
	}
}

// WithPeers specifies the peers to contact to join the cluster.
func WithPeers(peers ...Peer) Option {
	return func(c *Cluster) {
		c.peers = append(c.peers, peers...)
	}
}

// WithLogger specifies a logger.
func WithLogger(logger logging.Logger) Option {
	return func(c *Cluster) {
		c.logger = logger
	}
}
