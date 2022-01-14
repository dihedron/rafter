package cluster

import (
	"fmt"
	"net"
	"strconv"
)

type Address struct {
	Host string `json:"host,omitempty" yaml:"host,omitempty"`
	Port int    `json:"port,omitempty" yaml:"port,omitempty"`
}

func (a Address) String() string {
	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}

func (a *Address) UnmarshalFlag(value string) error {
	host, port, err := net.SplitHostPort(value)
	if err != nil {
		return fmt.Errorf("invalid format for address '%s': %w", value, err)
	}
	a.Host = host
	if a.Port, err = strconv.Atoi(port); err != nil {
		return fmt.Errorf("invalid format for port '%s': %w", port, err)
	}
	return nil
}
