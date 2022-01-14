package application

import (
	"io"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/hashicorp/raft"
)

func New() *Cache {
	return &Cache{
		cache: map[string]string{},
	}
}

// Cache keeps track of the three longest words it ever saw.
type Cache struct {
	mtx   sync.RWMutex
	words [3]string
	cache map[string]string
}

var _ raft.FSM = &Cache{}

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

func (f *Cache) Apply(l *raft.Log) interface{} {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	w := string(l.Data)
	for i := 0; i < len(f.words); i++ {
		if compareWords(w, f.words[i]) {
			copy(f.words[i+1:], f.words[i:])
			f.words[i] = w
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
