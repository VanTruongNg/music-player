package envelope

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/bwmarrin/snowflake"
)

var snowflakeNode *snowflake.Node

func init() {
	node, err := snowflake.NewNode(2) // Node ID = 2 for auth-service
	if err != nil {
		panic(fmt.Sprintf("failed to create snowflake node: %v", err))
	}
	snowflakeNode = node
}

func generateID() string {
	return snowflakeNode.Generate().String()
}

func (e *Envelope) Marshal() ([]byte, error) {
	data, err := json.Marshal(e)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal envelope: %w", err)
	}
	return data, nil
}

func Unmarshal(data []byte) (*Envelope, error) {
	var envelope Envelope
	if err := json.Unmarshal(data, &envelope); err != nil {
		return nil, fmt.Errorf("failed to unmarshal envelope: %w", err)
	}
	return &envelope, nil
}

func NewEnvelope(source string, priority Priority, data interface{}) (*Envelope, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	envelope := &Envelope{
		MessageID: generateID(),
		Version:   "1.0",
		Timestamp: time.Now().UTC(),
		Source:    source,
		Priority:  priority,
		Data:      dataBytes,
	}

	return envelope, nil
}

func (e *Envelope) GetData(v interface{}) error {
	if err := json.Unmarshal(e.Data, v); err != nil {
		return fmt.Errorf("failed to unmarshal envelope data: %w", err)
	}
	return nil
}
