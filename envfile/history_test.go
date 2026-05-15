package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func tempHistoryPath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "history.json")
}

func TestHistoryRecordAndReadAll(t *testing.T) {
	h := NewHistory(tempHistoryPath(t))

	if err := h.Record("encrypt", "production", []string{"DB_URL", "API_KEY"}, ""); err != nil {
		t.Fatalf("Record: %v", err)
	}
	if err := h.Record("rotate", "production", []string{"API_KEY"}, "scheduled rotation"); err != nil {
		t.Fatalf("Record: %v", err)
	}

	entries, err := h.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Action != "encrypt" {
		t.Errorf("expected action 'encrypt', got %q", entries[0].Action)
	}
	if entries[1].Note != "scheduled rotation" {
		t.Errorf("expected note 'scheduled rotation', got %q", entries[1].Note)
	}
}

func TestHistoryReadAllEmpty(t *testing.T) {
	h := NewHistory(tempHistoryPath(t))
	_, err := h.ReadAll()
	if !os.IsNotExist(err) {
		t.Errorf("expected not-exist error, got %v", err)
	}
}

func TestHistoryAppendsAcrossInstances(t *testing.T) {
	path := tempHistoryPath(t)
	h1 := NewHistory(path)
	if err := h1.Record("encrypt", "staging", []string{"SECRET"}, ""); err != nil {
		t.Fatalf("h1.Record: %v", err)
	}

	h2 := NewHistory(path)
	if err := h2.Record("decrypt", "staging", []string{"SECRET"}, ""); err != nil {
		t.Fatalf("h2.Record: %v", err)
	}

	entries, err := h2.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestHistoryFilePermissions(t *testing.T) {
	path := tempHistoryPath(t)
	h := NewHistory(path)
	if err := h.Record("keygen", "dev", nil, ""); err != nil {
		t.Fatalf("Record: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Errorf("expected 0600, got %04o", perm)
	}
}
