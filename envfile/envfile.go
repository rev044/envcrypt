package envfile

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Entry represents a single key-value pair in an env file.
type Entry struct {
	Key     string
	Value   string
	Comment string
}

// EnvFile holds all entries parsed from a .env file.
type EnvFile struct {
	Entries []Entry
}

// Parse reads and parses a .env file from the given path.
func Parse(path string) (*EnvFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening env file: %w", err)
	}
	defer f.Close()

	var ef EnvFile
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") {
			ef.Entries = append(ef.Entries, Entry{Comment: line})
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid line: %q", line)
		}
		ef.Entries = append(ef.Entries, Entry{
			Key:   strings.TrimSpace(parts[0]),
			Value: strings.TrimSpace(parts[1]),
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning env file: %w", err)
	}
	return &ef, nil
}

// Write serialises the EnvFile to the given path.
func Write(path string, ef *EnvFile) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating env file: %w", err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, e := range ef.Entries {
		if e.Comment != "" {
			fmt.Fprintln(w, e.Comment)
			continue
		}
		fmt.Fprintf(w, "%s=%s\n", e.Key, e.Value)
	}
	return w.Flush()
}

// ToMap converts the EnvFile entries into a plain map, skipping comments.
func (ef *EnvFile) ToMap() map[string]string {
	m := make(map[string]string, len(ef.Entries))
	for _, e := range ef.Entries {
		if e.Key != "" {
			m[e.Key] = e.Value
		}
	}
	return m
}

// FromMap creates an EnvFile from a plain map (order is not guaranteed).
func FromMap(m map[string]string) *EnvFile {
	ef := &EnvFile{}
	for k, v := range m {
		ef.Entries = append(ef.Entries, Entry{Key: k, Value: v})
	}
	return ef
}
