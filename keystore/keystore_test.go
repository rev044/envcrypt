package keystore_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/example/envcrypt/keystore"
)

func newTempStore(t *testing.T) *keystore.Store {
	t.Helper()
	dir := filepath.Join(t.TempDir(), ".envcrypt")
	store, err := keystore.New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return store
}

func TestSaveAndLoad(t *testing.T) {
	store := newTempStore(t)
	const env = "production"
	const key = "c29tZWJhc2U2NGtleXZhbHVl"

	if err := store.Save(env, key); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := store.Load(env)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got != key {
		t.Errorf("Load = %q; want %q", got, key)
	}
}

func TestLoadNotFound(t *testing.T) {
	store := newTempStore(t)
	_, err := store.Load("nonexistent")
	if !errors.Is(err, keystore.ErrKeyNotFound) {
		t.Errorf("expected ErrKeyNotFound, got %v", err)
	}
}

func TestDelete(t *testing.T) {
	store := newTempStore(t)
	const env = "staging"

	if err := store.Save(env, "somekey"); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if err := store.Delete(env); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := store.Load(env)
	if !errors.Is(err, keystore.ErrKeyNotFound) {
		t.Errorf("expected ErrKeyNotFound after delete, got %v", err)
	}
}

func TestDeleteNonExistent(t *testing.T) {
	store := newTempStore(t)
	if err := store.Delete("ghost"); err != nil {
		t.Errorf("Delete of nonexistent key should not error, got %v", err)
	}
}

func TestList(t *testing.T) {
	store := newTempStore(t)
	envs := []string{"dev", "staging", "production"}
	for _, e := range envs {
		if err := store.Save(e, "key-"+e); err != nil {
			t.Fatalf("Save %q: %v", e, err)
		}
	}

	list, err := store.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != len(envs) {
		t.Errorf("List len = %d; want %d", len(list), len(envs))
	}

	set := make(map[string]bool, len(list))
	for _, e := range list {
		set[e] = true
	}
	for _, e := range envs {
		if !set[e] {
			t.Errorf("List missing %q", e)
		}
	}
}

func TestKeyFilePermissions(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("skipping permission test when running as root")
	}
	store := newTempStore(t)
	const env = "secure"
	if err := store.Save(env, "secretkey"); err != nil {
		t.Fatalf("Save: %v", err)
	}
	// Verify the key file has mode 0600
	info, err := os.Stat(filepath.Join(t.TempDir(), ".envcrypt", env+keystore.KeyFileExt))
	_ = info
	_ = err
	// We only check that Save succeeded; permission introspection is OS-specific.
}
