package envfile

import (
	"fmt"
	"os"
	"path/filepath"
)

// Profile represents a named environment configuration (e.g. "dev", "staging", "prod").
type Profile struct {
	Name string
	Dir  string
}

// NewProfile creates a Profile rooted at dir with the given name.
func NewProfile(dir, name string) *Profile {
	return &Profile{Name: name, Dir: dir}
}

// Path returns the canonical file path for this profile's encrypted env file.
func (p *Profile) Path() string {
	return filepath.Join(p.Dir, fmt.Sprintf(".env.%s.enc", p.Name))
}

// Exists reports whether the profile file exists on disk.
func (p *Profile) Exists() bool {
	_, err := os.Stat(p.Path())
	return err == nil
}

// ListProfiles scans dir and returns all profile names whose encrypted env
// files are present (files matching the pattern .env.<name>.enc).
func ListProfiles(dir string) ([]string, error) {
	matches, err := filepath.Glob(filepath.Join(dir, ".env.*.enc"))
	if err != nil {
		return nil, fmt.Errorf("profile: glob %s: %w", dir, err)
	}
	names := make([]string, 0, len(matches))
	for _, m := range matches {
		base := filepath.Base(m)
		// base is ".env.<name>.enc"
		var name string
		if _, err := fmt.Sscanf(base, ".env.%s", &name); err != nil {
			continue
		}
		// strip trailing ".enc"
		if len(name) > 4 && name[len(name)-4:] == ".enc" {
			name = name[:len(name)-4]
		}
		if name != "" {
			names = append(names, name)
		}
	}
	return names, nil
}
