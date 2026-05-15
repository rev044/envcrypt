package envfile

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidationError holds details about a failed validation.
type ValidationError struct {
	Key     string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error for key %q: %s", e.Key, e.Message)
}

// ValidationResult aggregates all validation errors found in an env file.
type ValidationResult struct {
	Errors []ValidationError
}

func (r *ValidationResult) Valid() bool {
	return len(r.Errors) == 0
}

func (r *ValidationResult) Error() string {
	var sb strings.Builder
	for _, e := range r.Errors {
		sb.WriteString(e.Error())
		sb.WriteString("\n")
	}
	return strings.TrimRight(sb.String(), "\n")
}

var validKeyRe = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// Validate checks a slice of Entry values for common issues:
//   - empty keys
//   - keys containing invalid characters
//   - duplicate keys
func Validate(entries []Entry) *ValidationResult {
	result := &ValidationResult{}
	seen := make(map[string]int)

	for i, e := range entries {
		if e.Key == "" {
			result.Errors = append(result.Errors, ValidationError{
				Key:     fmt.Sprintf("<entry %d>", i),
				Message: "key must not be empty",
			})
			continue
		}
		if !validKeyRe.MatchString(e.Key) {
			result.Errors = append(result.Errors, ValidationError{
				Key:     e.Key,
				Message: "key contains invalid characters (must match [A-Za-z_][A-Za-z0-9_]*)",
			})
		}
		if prev, ok := seen[e.Key]; ok {
			result.Errors = append(result.Errors, ValidationError{
				Key:     e.Key,
				Message: fmt.Sprintf("duplicate key (first seen at entry %d)", prev),
			})
		} else {
			seen[e.Key] = i
		}
	}
	return result
}
