// Package envfile provides utilities for reading, writing, and manipulating
// .env files used by envcrypt.
//
// A .env file consists of lines in the form:
//
//	KEY=VALUE
//
// Lines beginning with '#' are treated as comments and are preserved during
// round-trip operations. Blank lines are skipped on parse.
//
// Typical usage:
//
//	ef, err := envfile.Parse(".env")
//	if err != nil { ... }
//	m := ef.ToMap()
//	// modify m ...
//	if err := envfile.Write(".env.enc", envfile.FromMap(m)); err != nil { ... }
package envfile
