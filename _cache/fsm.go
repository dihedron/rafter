package cache

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"github.com/hashicorp/raft"
)

type Cache struct {
	lock  sync.RWMutex
	cache map[string]string
}

func New() *Cache {
	return &Cache{cache: map[string]string{}}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	value, ok := c.cache[key]
	return value, ok
}

func (c *Cache) Set(key string, value string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.cache[key] = value
}

func (c *Cache) Del(key string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.cache, key)
}

// Apply applies a Raft log entry to the key-value store.
func (c *Cache) Apply(l *raft.Log) interface{} {
	var m Message
	if err := json.Unmarshal(l.Data, &m); err != nil {
		panic(fmt.Sprintf("failed to unmarshal message: %s", err.Error()))
	}
	c.Set(m.ID, m.Payload)
	return nil
}

// Snapshot returns a snapshot of the key-value store.
func (c *Cache) Snapshot() (raft.FSMSnapshot, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	clone := make(map[string]string)
	for k, v := range c.cache {
		clone[k] = v
	}
	return &Snapshot{snapshot: clone}, nil
}

// Restore stores the key-value store to a previous state.
func (s *Cache) Restore(rc io.ReadCloser) error {
	clone := make(map[string]string)
	if err := json.NewDecoder(rc).Decode(&clone); err != nil {
		return err
	}
	// Set the state from the snapshot, no lock required according to
	// Hashicorp docs.
	s.cache = clone
	return nil
}
