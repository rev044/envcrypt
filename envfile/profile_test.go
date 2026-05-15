package envfile

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestProfilePath(t *testing.T) {
	p := NewProfile("/tmp/myproject", "staging")
	want := "/tmp/myproject/.env.staging.enc"
	if p.Path() != want {
		t.Errorf("Path() = %q, want %q", p.Path(), want)
	}
}

func TestProfileExistsAndMissing(t *testing.T) {
	dir := t.TempDir()
	p := NewProfile(dir, "dev")

	if p.Exists() {
		t.Fatal("expected profile to not exist yet")
	}

	// Create the file
	if err := os.WriteFile(p.Path(), []byte("data"), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	if !p.Exists() {
		t.Fatal("expected profile to exist after creation")
	}
}

func TestListProfiles(t *testing.T) {
	dir := t.TempDir()

	for _, name := range []string{"dev", "staging", "prod"} {
		path := filepath.Join(dir, ".env."+name+".enc")
		if err := os.WriteFile(path, []byte("x"), 0o600); err != nil {
			t.Fatalf("WriteFile: %v", err)
		}
	}

	// Add a non-matching file that should be ignored.
	_ = os.WriteFile(filepath.Join(dir, ".env"), []byte("x"), 0o600)
	_ = os.WriteFile(filepath.Join(dir, "notes.txt"), []byte("x"), 0o600)

	profiles, err := ListProfiles(dir)
	if err != nil {
		t.Fatalf("ListProfiles: %v", err)
	}

	sort.Strings(profiles)
	want := []string{"dev", "prod", "staging"}
	if len(profiles) != len(want) {
		t.Fatalf("got %v, want %v", profiles, want)
	}
	for i, w := range want {
		if profiles[i] != w {
			t.Errorf("profiles[%d] = %q, want %q", i, profiles[i], w)
		}
	}
}

func TestListProfilesEmptyDir(t *testing.T) {
	dir := t.TempDir()
	profiles, err := ListProfiles(dir)
	if err != nil {
		t.Fatalf("ListProfiles: %v", err)
	}
	if len(profiles) != 0 {
		t.Errorf("expected empty list, got %v", profiles)
	}
}
