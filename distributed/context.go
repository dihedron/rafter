package distributed

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"sync"

	"github.com/dihedron/rafter/logging"
	"github.com/hashicorp/raft"
)

func NewContext(l logging.Logger) *Context {
	l.Info("creating new distributed context...")
	return &Context{
		values: map[string][]byte{},
		logger: l,
	}
}

// Context is a cluster-wide shared, distributed context.
type Context struct {
	mtx    sync.RWMutex
	values map[string][]byte
	logger logging.Logger
}

var _ raft.FSM = &Context{}

func (c *Context) Apply(l *raft.Log) interface{} {
	var err error
	c.logger.Trace("applying log entry: %s", logging.ToJSON(l))
	message := &Message{}
	if err = json.Unmarshal(l.Data, message); err != nil {
		c.logger.Error("error unmarshalling message: %v", err)
		return fmt.Errorf("error unmarshalling input message: %w", err)
	}
	var result *Message
	switch message.Type {
	case Get:
		c.mtx.RLock()
		value := c.values[message.Key]
		c.mtx.RUnlock()
		result = &Message{
			Key:   message.Key,
			Value: []byte(value),
			Index: l.Index,
		}
	case Set:
		c.mtx.Lock()
		c.values[message.Key] = message.Value
		c.mtx.Unlock()
		result = &Message{
			Index: l.Index,
		}
	case Remove:
		c.mtx.Lock()
		value := c.values[message.Key]
		delete(c.values, message.Key)
		c.mtx.Unlock()
		result = &Message{
			Key:   message.Key,
			Value: value,
			Index: l.Index,
		}
	case List:
		var re *regexp.Regexp
		if message.Filter != "" {
			if re, err = regexp.Compile(message.Filter); err != nil {
				c.logger.Error("error compiling regular expression '%s': %v", message.Filter, err)
				return fmt.Errorf("error compiling regular expression '%s': %w", message.Filter, err)
			}
		}
		c.mtx.RLock()
		keys := []string{}
		for k := range c.values {
			if re == nil || re.Match([]byte(k)) {
				keys = append(keys, k)
			}
		}
		c.mtx.RUnlock()
		result = &Message{
			Keys:  keys,
			Index: l.Index,
		}
	case Clear:
		c.mtx.Lock()
		value := c.values[message.Key]
		c.values[message.Key] = message.Value
		c.mtx.Unlock()
		result = &Message{
			Key:   message.Key,
			Value: value,
			Index: l.Index,
		}
	}

	data, err := json.Marshal(result)
	if err != nil {
		c.logger.Error("error marshalling response: %v", err)
		return nil
	}
	return data
}

func (c *Context) Snapshot() (raft.FSMSnapshot, error) {
	// Make sure that any future calls to f.Apply() don't change the snapshot.
	data, err := json.Marshal(c.values)
	if err != nil {
		return nil, fmt.Errorf("error marshalling snapshot content to JSON: %w", err)
	}
	return &Snapshot{data: data}, nil
}

func (c *Context) Restore(r io.ReadCloser) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	cache := map[string][]byte{}
	if err := json.Unmarshal(data, &cache); err != nil {
		return fmt.Errorf("error unmarshalling snapshot content from JSON: %w", err)
	}
	c.values = cache
	return nil
}
