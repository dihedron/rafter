package cache

import (
	"encoding/json"

	"github.com/hashicorp/raft"
)

type Snapshot struct {
	snapshot map[string]string
}

func (s *Snapshot) Persist(sink raft.SnapshotSink) error {
	err := func() error {
		b, err := json.Marshal(s.snapshot)
		if err != nil {
			return err
		}
		if _, err := sink.Write(b); err != nil {
			return err
		}
		return sink.Close()
	}()

	if err != nil {
		sink.Cancel()
	}

	return err
}

func (s *Snapshot) Release() {

}
