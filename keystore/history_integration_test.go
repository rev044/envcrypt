package keystore_test

import (
	"path/filepath"
	"testing"

	"github.com/user/envcrypt/envfile"
	"github.com/user/envcrypt/keystore"
)

// TestKeyRotationRecordedInHistory verifies that a simulated key rotation
// workflow records the correct history entries.
func TestKeyRotationRecordedInHistory(t *testing.T) {
	dir := t.TempDir()
	store := keystore.New(filepath.Join(dir, "keys"))
	history := envfile.NewHistory(filepath.Join(dir, "history.json"))

	// Save an initial key.
	if err := store.Save("production", []byte("aaaabbbbccccddddeeeeffffgggghhhh")); err != nil {
		t.Fatalf("Save initial key: %v", err)
	}
	if err := history.Record("keygen", "production", nil, "initial key"); err != nil {
		t.Fatalf("Record keygen: %v", err)
	}

	// Simulate rotation: save new key, record event.
	if err := store.Save("production", []byte("11112222333344445555666677778888")); err != nil {
		t.Fatalf("Save rotated key: %v", err)
	}
	if err := history.Record("rotate", "production", []string{"DB_URL", "API_KEY"}, "manual rotation"); err != nil {
		t.Fatalf("Record rotate: %v", err)
	}

	entries, err := history.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 history entries, got %d", len(entries))
	}
	if entries[1].Action != "rotate" {
		t.Errorf("expected second action 'rotate', got %q", entries[1].Action)
	}
	if len(entries[1].Keys) != 2 {
		t.Errorf("expected 2 keys in rotate entry, got %d", len(entries[1].Keys))
	}
}
