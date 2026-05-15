package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTakeSnapshotAndVerify(t *testing.T) {
	entries := []Entry{
		{Key: "APP_ENV", Value: "production"},
		{Key: "DB_HOST", Value: "localhost"},
	}
	snap := TakeSnapshot("prod", entries)

	if snap.Environment != "prod" {
		t.Errorf("expected environment 'prod', got %q", snap.Environment)
	}
	if len(snap.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(snap.Entries))
	}
	if snap.Checksum == "" {
		t.Error("expected non-empty checksum")
	}
	if !snap.Verify() {
		t.Error("expected Verify() to return true for fresh snapshot")
	}
}

func TestSnapshotVerifyTampered(t *testing.T) {
	entries := []Entry{
		{Key: "SECRET", Value: "abc123"},
	}
	snap := TakeSnapshot("staging", entries)

	// Tamper with an entry value after snapshot is taken.
	snap.Entries[0].Value = "hacked"

	if snap.Verify() {
		t.Error("expected Verify() to return false after tampering")
	}
}

func TestSnapshotVerifyTamperedKey(t *testing.T) {
	entries := []Entry{
		{Key: "SECRET", Value: "abc123"},
	}
	snap := TakeSnapshot("staging", entries)

	// Tamper with an entry key after snapshot is taken.
	snap.Entries[0].Key = "HACKED_KEY"

	if snap.Verify() {
		t.Error("expected Verify() to return false after tampering with key")
	}
}

func TestSaveAndLoadSnapshot(t *testing.T) {
	entries := []Entry{
		{Key: "PORT", Value: "8080"},
		{Key: "LOG_LEVEL", Value: "debug"},
	}
	snap := TakeSnapshot("dev", entries)

	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	if err := SaveSnapshot(path, snap); err != nil {
		t.Fatalf("SaveSnapshot: %v", err)
	}

	loaded, err := LoadSnapshot(path)
	if err != nil {
		t.Fatalf("LoadSnapshot: %v", err)
	}

	if loaded.Environment != snap.Environment {
		t.Errorf("environment mismatch: got %q, want %q", loaded.Environment, snap.Environment)
	}
	if loaded.Checksum != snap.Checksum {
		t.Errorf("checksum mismatch: got %q, want %q", loaded.Checksum, snap.Checksum)
	}
	if !loaded.Verify() {
		t.Error("loaded snapshot failed Verify()")
	}
	if len(loaded.Entries) != len(snap.Entries) {
		t.Errorf("entry count mismatch: got %d, want %d", len(loaded.Entries), len(snap.Entries))
	}
}

func TestLoadSnapshotFilePermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	snap := TakeSnapshot("test", []Entry{{Key: "X", Value: "1"}})
	if err := SaveSnapshot(path, snap); err != nil {
		t.Fatalf("SaveSnapshot: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("expected file permission 0600, got %04o", perm)
	}
}

func TestLoadSnapshotNotFound(t *testing.T) {
	_, err := LoadSnapshot("/nonexistent/path/snap.json")
	if err == nil {
		t.Error("expected error loading non-existent snapshot, got nil")
	}
}
