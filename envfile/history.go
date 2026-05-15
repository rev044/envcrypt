package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// HistoryEntry records a single change event for an env file.
type HistoryEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Action    string    `json:"action"`
	Profile   string    `json:"profile"`
	Keys      []string  `json:"keys"`
	Note      string    `json:"note,omitempty"`
}

// History holds a sequence of HistoryEntry records for a profile.
type History struct {
	path string
}

// NewHistory returns a History backed by the given file path.
func NewHistory(path string) *History {
	return &History{path: path}
}

// Record appends a new entry to the history file.
func (h *History) Record(action, profile string, keys []string, note string) error {
	entries, err := h.ReadAll()
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("history: read: %w", err)
	}
	entries = append(entries, HistoryEntry{
		Timestamp: time.Now().UTC(),
		Action:    action,
		Profile:   profile,
		Keys:      keys,
		Note:      note,
	})
	return h.write(entries)
}

// ReadAll returns all recorded history entries.
func (h *History) ReadAll() ([]HistoryEntry, error) {
	data, err := os.ReadFile(h.path)
	if err != nil {
		return nil, err
	}
	var entries []HistoryEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("history: parse: %w", err)
	}
	return entries, nil
}

func (h *History) write(entries []HistoryEntry) error {
	if err := os.MkdirAll(filepath.Dir(h.path), 0o700); err != nil {
		return fmt.Errorf("history: mkdir: %w", err)
	}
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("history: marshal: %w", err)
	}
	return os.WriteFile(h.path, data, 0o600)
}
