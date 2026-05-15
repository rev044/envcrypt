// Package envfile provides primitives for parsing, writing, validating,
// merging, exporting, snapshotting, and templating .env files.
//
// # Profiles
//
// A Profile represents a named environment tier (e.g. "dev", "staging",
// "prod") and maps to a corresponding encrypted file on disk following the
// naming convention:
//
//	.env.<name>.enc
//
// Use NewProfile to construct a profile for a specific directory and name,
// then call Path to obtain the resolved file path or Exists to check whether
// the encrypted file is already present.
//
// ListProfiles scans a directory and returns the names of all profiles whose
// encrypted files are found, making it easy to enumerate available
// environments without hard-coding them.
//
// Example:
//
//	profiles, err := envfile.ListProfiles("./secrets")
//	if err != nil { ... }
//	for _, name := range profiles {
//		p := envfile.NewProfile("./secrets", name)
//		fmt.Println(p.Path())
//	}
package envfile
