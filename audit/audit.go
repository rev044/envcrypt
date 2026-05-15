// Package audit provides a simple append-only audit log for tracking
// key operations such as encryption, decryption, and key rotation.
package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// EventKind describes the type of auditable operation.
type EventKind string

const (
	EventEncrypt EventKind = "encrypt"
	EventDecrypt EventKind = "decrypt"
	EventRotate  EventKind = "rotate"
	EventKeygen  EventKind = "keygen"
	EventDelete  EventKind = "delete"
)

// Event represents a single audit log entry.
type Event struct {
	Timestamp time.Time `json:"timestamp"`
	Kind      EventKind `json:"kind"`
	Env       string    `json:"env"`
	KeyID     string    `json:"key_id,omitempty"`
	Message   string    `json:"message,omitempty"`
}

// Logger writes audit events to a newline-delimited JSON file.
type Logger struct {
	path string
}

// New creates a Logger that appends events to the file at path.
func New(path string) *Logger {
	return &Logger{path: path}
}

// Log appends an Event to the audit log file.
func (l *Logger) Log(kind EventKind, env, keyID, message string) error {
	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("audit: open log file: %w", err)
	}
	defer f.Close()

	e := Event{
		Timestamp: time.Now().UTC(),
		Kind:      kind,
		Env:       env,
		KeyID:     keyID,
		Message:   message,
	}
	if err := json.NewEncoder(f).Encode(e); err != nil {
		return fmt.Errorf("audit: encode event: %w", err)
	}
	return nil
}

// ReadAll reads and returns all events from the audit log file.
func (l *Logger) ReadAll() ([]Event, error) {
	f, err := os.Open(l.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("audit: open log file: %w", err)
	}
	defer f.Close()

	var events []Event
	dec := json.NewDecoder(f)
	for dec.More() {
		var e Event
		if err := dec.Decode(&e); err != nil {
			return nil, fmt.Errorf("audit: decode event: %w", err)
		}
		events = append(events, e)
	}
	return events, nil
}
