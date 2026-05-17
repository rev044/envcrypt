// Package envfile provides utilities for parsing, writing, and managing
// .env files across multiple environments.
//
// # Redaction
//
// The redact module identifies and masks sensitive environment variable
// values before they are displayed, logged, or exported.
//
// Keys are considered sensitive if they match common naming patterns such as:
//   - PASSWORD, PASSWD, PWD
//   - SECRET, TOKEN, API_KEY, APIKEY
//   - PRIVATE_KEY, PRIVKEY
//   - AUTH, CREDENTIAL
//
// Example usage:
//
//	entries := []Entry{
//	    {Key: "DB_PASSWORD", Value: "s3cr3t"},
//	    {Key: "APP_NAME",    Value: "myapp"},
//	}
//	safe := envfile.Redact(entries, envfile.RedactOptions{Mask: "[hidden]"})
//	// safe[0].Value == "[hidden]"
//	// safe[1].Value == "myapp"
package envfile
