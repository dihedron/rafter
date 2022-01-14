package cache

type Message struct {
	ID      string `json:"id,omitempty"`
	Payload string `json:"payload,omitempty"`
}
