// Package envfile provides utilities for parsing, writing, merging,
// diffing, and validating .env files.
//
// # Parsing
//
// Parse reads KEY=VALUE lines from an io.Reader, skipping blank lines
// and comments (lines starting with '#'). It returns a slice of Entry
// values that preserve the original order.
//
// # Writing
//
// Write serialises a slice of Entry values back to KEY=VALUE format.
//
// # Merging and Diffing
//
// Merge combines two sets of entries with configurable conflict
// resolution (prefer base, prefer override, or error on conflict).
// Diff returns the keys that differ between two entry sets.
//
// # Validation
//
// Validate inspects a slice of Entry values and reports empty keys,
// keys with invalid characters, and duplicate keys as a ValidationResult.
package envfile
