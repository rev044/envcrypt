// Package envfile provides utilities for parsing, writing, merging,
// validating, exporting, snapshotting, and tracking the history of
// environment variable files.
//
// # History
//
// The History type records timestamped change events for a named profile.
// Each event captures the action performed (e.g. "encrypt", "rotate"),
// the profile name, the affected keys, and an optional human-readable note.
//
// History entries are stored as a JSON array in a file on disk, with
// permissions restricted to 0600 so that secrets are not inadvertently
// exposed through log files.
//
// Example usage:
//
//	h := envfile.NewHistory(".envcrypt/history/production.json")
//	err := h.Record("rotate", "production", []string{"API_KEY"}, "quarterly rotation")
//	entries, err := h.ReadAll()
package envfile
