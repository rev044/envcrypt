// Package envfile provides utilities for parsing, writing, validating,
// merging, exporting, and linting .env files.
//
// # Lint
//
// The Lint function performs style and hygiene checks on a slice of Entry
// values without modifying them. It returns a (possibly empty) slice of
// LintIssue, each carrying a line number, the offending key, a human-readable
// message, and a severity level of either LintWarn or LintError.
//
// Current checks:
//
//   - Empty value                  → warn
//   - Key is not fully upper-case  → warn
//   - Value has leading/trailing whitespace → error
//   - Key starts with a digit      → warn
//
// Example:
//
//	entries, _ := envfile.Parse(r)
//	issues := envfile.Lint(entries)
//	for _, issue := range issues {
//		fmt.Println(issue)
//	}
package envfile
