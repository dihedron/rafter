package distributed

import (
	"fmt"

	"github.com/hashicorp/raft"
)

type Snapshot struct {
	data []byte
}

func (s *Snapshot) Persist(sink raft.SnapshotSink) error {
	_, err := sink.Write(s.data)
	if err != nil {
		sink.Cancel()
		return fmt.Errorf("error writing snapshot to sink: %v", err)
	}
	return sink.Close()
}

func (s *Snapshot) Release() {
}
