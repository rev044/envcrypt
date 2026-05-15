// Package envfile provides utilities for parsing, writing, validating,
// merging, snapshotting, and exporting .env files.
//
// # Export
//
// The Export function serialises a slice of Entry values into one of three
// standard formats:
//
//   - FormatShell  — "export KEY=VALUE" lines, safe to source in bash/zsh.
//     Values that contain shell-special characters are automatically
//     wrapped in single quotes.
//
//   - FormatDocker — plain "KEY=VALUE" lines accepted by Docker's
//     --env-file flag and docker-compose env_file directive.
//
//   - FormatJSON   — a JSON object mapping each key to its string value,
//     useful for piping into other tooling (e.g. AWS Secrets Manager).
//
// Example:
//
//	entries, _ := envfile.Parse(r)
//	envfile.Export(entries, envfile.FormatShell, os.Stdout)
package envfile
