package application

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"strings"
	"sync"
	"time"

	"github.com/dihedron/rafter/logging"
	"github.com/hashicorp/raft"
)

func New(l logging.Logger) *Cache {
	l.Info("creating application...")
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
	c.logger.Debug("log entry (type %T): %s", l, logging.ToJSON(l))
	time.Sleep(50 * time.Millisecond)
	message := &Message{}
	if err := json.Unmarshal(l.Data, message); err != nil {
		c.logger.Error("error unmarshalling message: %v", err)
		return nil
	}
	var result *Message
	switch message.Type {
	case Get:
		c.mtx.RLock()
		value := c.cache[message.Key]
		c.mtx.RUnlock()
		result = &Message{
			Key:   message.Key,
			Value: value,
			Index: l.Index,
		}
	case Set:
		c.mtx.Lock()
		c.cache[message.Key] = message.Value
		c.mtx.Unlock()
		result = &Message{
			Index: l.Index,
		}
	case Remove:
		c.mtx.Lock()
		value := c.cache[message.Key]
		c.cache[message.Key] = message.Value
		c.mtx.Unlock()
		result = &Message{
			Key:   message.Key,
			Value: value,
			Index: l.Index,
		}
	case List:

	case Clear:
	}

	data, err := json.Marshal(result)
	if err != nil {
		c.logger.Error("error marshalling response: %v", err)
		return nil
	}
	return data

	// w := string(l.Data)
	// for i := 0; i < len(c.words); i++ {
	// 	if compareWords(w, c.words[i]) {
	// 		copy(c.words[i+1:], c.words[i:])
	// 		c.words[i] = w
	// 		break
	// 	}
	// }
	// return nil
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
