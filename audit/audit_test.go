package audit_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envcrypt/audit"
)

func tempLogPath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "audit.log")
}

func TestLogAndReadAll(t *testing.T) {
	path := tempLogPath(t)
	l := audit.New(path)

	if err := l.Log(audit.EventKeygen, "production", "key-abc", "generated new key"); err != nil {
		t.Fatalf("Log: %v", err)
	}
	if err := l.Log(audit.EventEncrypt, "production", "key-abc", ""); err != nil {
		t.Fatalf("Log: %v", err)
	}

	events, err := l.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
	if events[0].Kind != audit.EventKeygen {
		t.Errorf("expected keygen, got %s", events[0].Kind)
	}
	if events[1].Env != "production" {
		t.Errorf("expected env production, got %s", events[1].Env)
	}
	if events[0].KeyID != "key-abc" {
		t.Errorf("expected key-abc, got %s", events[0].KeyID)
	}
}

func TestReadAllEmptyFile(t *testing.T) {
	path := tempLogPath(t)
	l := audit.New(path)

	events, err := l.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll on missing file: %v", err)
	}
	if events != nil {
		t.Errorf("expected nil events for missing file, got %v", events)
	}
}

func TestLogAppendsAcrossInstances(t *testing.T) {
	path := tempLogPath(t)

	audit.New(path).Log(audit.EventRotate, "staging", "key-1", "old key")
	audit.New(path).Log(audit.EventRotate, "staging", "key-2", "new key")

	events, err := audit.New(path).ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
}

func TestLogFilePermissions(t *testing.T) {
	path := tempLogPath(t)
	l := audit.New(path)
	l.Log(audit.EventDecrypt, "dev", "", "")

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Errorf("expected 0600 permissions, got %o", perm)
	}
}
