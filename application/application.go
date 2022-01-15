package application

import (
	"io"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/dihedron/rafter/logging"
	"github.com/hashicorp/raft"
)

func New(l logging.Logger) *Cache {
	return &Cache{
		cache:  map[string]string{},
		logger: l,
	}
}

// Cache keeps track of the three longest words it ever saw.
type Cache struct {
	mtx    sync.RWMutex
	words  [3]string
	cache  map[string]string
	logger logging.Logger
}

var _ raft.FSM = &Cache{}

func (c *Cache) Apply(l *raft.Log) interface{} {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.logger.Debug("log entry is of type %T", l)
	w := string(l.Data)
	for i := 0; i < len(c.words); i++ {
		if compareWords(w, c.words[i]) {
			copy(c.words[i+1:], c.words[i:])
			c.words[i] = w
			break
		}
	}
	return nil
}

func (f *Cache) Snapshot() (raft.FSMSnapshot, error) {
	// Make sure that any future calls to f.Apply() don't change the snapshot.
	return &Snapshot{cloneWords(f.words)}, nil
}

func (f *Cache) Restore(r io.ReadCloser) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	words := strings.Split(string(b), "\n")
	copy(f.words[:], words)
	return nil
}

// compareWords returns true if a is longer (lexicography breaking ties).
func compareWords(a, b string) bool {
	if len(a) == len(b) {
		return a < b
	}
	return len(a) > len(b)
}

func cloneWords(words [3]string) []string {
	var ret [3]string
	copy(ret[:], words[:])
	return ret[:]
}
