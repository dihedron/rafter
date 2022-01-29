package application

type Type int8

const (
	Get Type = iota
	Set
	Remove
	List
	Clear
)

type Message struct {
	Type   Type     `json:"type"`
	Key    string   `json:"key,omitempty"`
	Value  []byte   `json:"value,omitempty"`
	Filter string   `json:"filter,omitempty"`
	Keys   []string `json:"keys,omitempty"`
	Index  uint64   `json:"index,omitempty"`
}
