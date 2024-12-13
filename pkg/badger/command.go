package badger

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

// Command represents a Raft command.
type Command struct {
	Op    string
	Key   string
	Value []byte
}

func (c *Command) Serialize() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(c)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize command: %v", err)
	}
	return buf.Bytes(), nil
}

func (c *Command) Deserialize(data []byte) error {
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	err := decoder.Decode(c)
	if err != nil {
		return fmt.Errorf("failed to deserialize command: %v", err)
	}
	return nil
}
