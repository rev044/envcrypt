// Package envfile provides utilities for parsing, writing, merging,
// validating, and snapshotting .env files.
//
// # Snapshots
//
// A Snapshot captures the complete state of an environment's entries at a
// specific point in time, together with a FNV-1a checksum that can be used
// to detect tampering or unintended modifications.
//
// Basic usage:
//
//	entries, _ := envfile.Parse(r)
//	snap := envfile.TakeSnapshot("production", entries)
//
//	// Persist to disk (mode 0600).
//	_ = envfile.SaveSnapshot("/var/envcrypt/snapshots/prod.json", snap)
//
//	// Later, reload and verify integrity.
//	loaded, _ := envfile.LoadSnapshot("/var/envcrypt/snapshots/prod.json")
//	if !loaded.Verify() {
//		log.Fatal("snapshot integrity check failed")
//	}
package envfile
