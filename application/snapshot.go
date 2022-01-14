package application

import (
	"fmt"
	"strings"

	"github.com/hashicorp/raft"
)

type Snapshot struct {
	words []string
}

func (s *Snapshot) Persist(sink raft.SnapshotSink) error {
	_, err := sink.Write([]byte(strings.Join(s.words, "\n")))
	if err != nil {
		sink.Cancel()
		return fmt.Errorf("sink.Write(): %v", err)
	}
	return sink.Close()
}

func (s *Snapshot) Release() {
}
