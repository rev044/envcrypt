// Package audit implements an append-only, newline-delimited JSON audit log
// for the envcrypt tool.
//
// Each operation that modifies or accesses encrypted environment data — such as
// key generation, encryption, decryption, rotation, and deletion — should be
// recorded via a Logger so that operators can review a tamper-evident history
// of activity per environment.
//
// Usage:
//
//	l := audit.New(".envcrypt/audit.log")
//	l.Log(audit.EventEncrypt, "production", "key-abc123", "encrypted .env")
//
//	events, err := l.ReadAll()
//
The log file is created with 0600 permissions and events are written one per
// line as JSON objects containing a UTC timestamp, event kind, environment
// name, optional key ID, and an optional human-readable message.
package audit
