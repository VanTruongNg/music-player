package envelope

import (
	"fmt"
	"time"
)

func (e *Envelope) Validate() error {
	if e.MessageID == "" {
		return fmt.Errorf("message_id is required")
	}

	if e.Source == "" {
		return fmt.Errorf("source is required")
	}

	if e.Timestamp.IsZero() {
		return fmt.Errorf("timestamp is required")
	}

	if e.Timestamp.After(time.Now().UTC().Add(1 * time.Minute)) {
		return fmt.Errorf("timestamp cannot be in the future")
	}

	if !e.Priority.IsValid() {
		return fmt.Errorf("invalid priority: %s", e.Priority)
	}

	if len(e.Data) == 0 {
		return fmt.Errorf("data is required")
	}

	return nil
}
