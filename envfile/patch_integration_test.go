package envfile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envcrypt/envfile"
)

// TestPatchRoundTripViaFiles verifies that DiffEntries + Patch produces an
// env file identical to the "next" file when starting from the "base" file.
func TestPatchRoundTripViaFiles(t *testing.T) {
	dir := t.TempDir()

	basePath := filepath.Join(dir, "base.env")
	nextPath := filepath.Join(dir, "next.env")

	_ = os.WriteFile(basePath, []byte("FOO=bar\nOLD=gone\n"), 0600)
	_ = os.WriteFile(nextPath, []byte("FOO=baz\nNEW=val\n"), 0600)

	base, err := envfile.Parse(basePath)
	if err != nil {
		t.Fatalf("parse base: %v", err)
	}
	next, err := envfile.Parse(nextPath)
	if err != nil {
		t.Fatalf("parse next: %v", err)
	}

	changes := envfile.DiffEntries(base, next)
	patched, err := envfile.Patch(base, changes)
	if err != nil {
		t.Fatalf("patch: %v", err)
	}

	nextMap := envfile.ToMap(next)
	patchedMap := envfile.ToMap(patched)

	for k, v := range nextMap {
		if patchedMap[k] != v {
			t.Errorf("key %s: want %q, got %q", k, v, patchedMap[k])
		}
	}
	if len(patchedMap) != len(nextMap) {
		t.Errorf("length mismatch: patched=%d next=%d", len(patchedMap), len(nextMap))
	}
}
