package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot captures the state of an env file at a point in time.
type Snapshot struct {
	Timestamp   time.Time `json:"timestamp"`
	Environment string    `json:"environment"`
	Entries     []Entry   `json:"entries"`
	Checksum    string    `json:"checksum"`
}

// TakeSnapshot creates a Snapshot from a slice of entries.
func TakeSnapshot(environment string, entries []Entry) Snapshot {
	return Snapshot{
		Timestamp:   time.Now().UTC(),
		Environment: environment,
		Entries:     entries,
		Checksum:    checksumEntries(entries),
	}
}

// SaveSnapshot writes a Snapshot to a JSON file at the given path.
func SaveSnapshot(path string, snap Snapshot) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("snapshot: open %q: %w", path, err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(snap); err != nil {
		return fmt.Errorf("snapshot: encode: %w", err)
	}
	return nil
}

// LoadSnapshot reads a Snapshot from a JSON file at the given path.
func LoadSnapshot(path string) (Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: open %q: %w", path, err)
	}
	defer f.Close()
	var snap Snapshot
	if err := json.NewDecoder(f).Decode(&snap); err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: decode: %w", err)
	}
	return snap, nil
}

// Verify checks whether the snapshot checksum matches its entries.
func (s Snapshot) Verify() bool {
	return checksumEntries(s.Entries) == s.Checksum
}

// checksumEntries produces a deterministic FNV-1a hex checksum over entry key=value pairs.
func checksumEntries(entries []Entry) string {
	var h uint64 = 14695981039346656037
	for _, e := range entries {
		for _, c := range e.Key + "=" + e.Value {
			h ^= uint64(c)
			h *= 1099511628211
		}
		h ^= '\n'
		h *= 1099511628211
	}
	return fmt.Sprintf("%016x", h)
}
