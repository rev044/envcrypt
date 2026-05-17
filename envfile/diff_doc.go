// Package envfile provides utilities for parsing, writing, and managing
// .env files across environments.
//
// # Diff and Patch
//
// The diff/patch subsystem allows you to compute a structured changeset
// between two slices of Entry values and apply that changeset to produce
// a new slice.
//
//	base := []Entry{{Key: "FOO", Value: "bar"}}
//	next := []Entry{{Key: "FOO", Value: "baz"}, {Key: "NEW", Value: "val"}}
//
//	changes := DiffEntries(base, next)
//	patched, err := Patch(base, changes)
//
// ChangeKind values:
//
//	ChangeAdded    – key is present in next but not in base
//	ChangeModified – key exists in both but values differ
//	ChangeRemoved  – key is present in base but not in next
package envfile
