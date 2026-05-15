// Package keystore manages encryption keys for different environments,
// supporting storage, retrieval, and rotation of keys on disk.
package keystore

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	// DefaultDir is the default directory for storing keys.
	DefaultDir = ".envcrypt"
	// KeyFileExt is the file extension for key files.
	KeyFileExt = ".key"
)

// ErrKeyNotFound is returned when a key for the given environment does not exist.
var ErrKeyNotFound = errors.New("key not found")

// Store manages keys stored in a directory on disk.
type Store struct {
	dir string
}

// New creates a new Store rooted at dir, creating the directory if necessary.
func New(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("keystore: create directory: %w", err)
	}
	return &Store{dir: dir}, nil
}

// keyPath returns the file path for the given environment key.
func (s *Store) keyPath(env string) string {
	return filepath.Join(s.dir, env+KeyFileExt)
}

// Save writes the base64-encoded key to disk for the given environment.
// The file is created with restrictive permissions (0600).
func (s *Store) Save(env, encodedKey string) error {
	path := s.keyPath(env)
	if err := os.WriteFile(path, []byte(strings.TrimSpace(encodedKey)+"\n"), 0600); err != nil {
		return fmt.Errorf("keystore: save key for %q: %w", env, err)
	}
	return nil
}

// Load reads and returns the base64-encoded key for the given environment.
// Returns ErrKeyNotFound if no key file exists for that environment.
func (s *Store) Load(env string) (string, error) {
	path := s.keyPath(env)
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return "", ErrKeyNotFound
	}
	if err != nil {
		return "", fmt.Errorf("keystore: load key for %q: %w", env, err)
	}
	return strings.TrimSpace(string(data)), nil
}

// Delete removes the key file for the given environment.
func (s *Store) Delete(env string) error {
	path := s.keyPath(env)
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("keystore: delete key for %q: %w", env, err)
	}
	return nil
}

// List returns all environment names that have a stored key.
func (s *Store) List() ([]string, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, fmt.Errorf("keystore: list keys: %w", err)
	}
	var envs []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), KeyFileExt) {
			envs = append(envs, strings.TrimSuffix(e.Name(), KeyFileExt))
		}
	}
	return envs, nil
}
