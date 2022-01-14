package cluster

import (
	"fmt"

	"github.com/dihedron/rafter/unmarshal"
)

type Peer struct {
	ID      string  `json:"id,omitempty" yaml:"id,omitempty"`
	Address Address `json:"address,omitempty" yaml:"address,omitempty"`
}

func (p *Peer) UnmarshalFlag(value string) error {
	return unmarshal.FromFlag(value, p)
}

func (p Peer) String() string {
	return fmt.Sprintf("%s@%+v", p.ID, p.Address)
}
