package envfile

import (
	"fmt"
	"strings"
)

// LintSeverity indicates how severe a lint issue is.
type LintSeverity string

const (
	LintWarn  LintSeverity = "warn"
	LintError LintSeverity = "error"
)

// LintIssue represents a single linting finding.
type LintIssue struct {
	Line     int
	Key      string
	Message  string
	Severity LintSeverity
}

func (i LintIssue) String() string {
	return fmt.Sprintf("[%s] line %d (%s): %s", i.Severity, i.Line, i.Key, i.Message)
}

// Lint runs style and hygiene checks on a slice of Entry values and returns
// any findings. It does not modify the entries.
func Lint(entries []Entry) []LintIssue {
	var issues []LintIssue

	for idx, e := range entries {
		lineNum := idx + 1

		// Warn on empty values
		if e.Value == "" {
			issues = append(issues, LintIssue{
				Line:     lineNum,
				Key:      e.Key,
				Message:  "value is empty",
				Severity: LintWarn,
			})
		}

		// Warn on keys that are not upper-case
		if e.Key != strings.ToUpper(e.Key) {
			issues = append(issues, LintIssue{
				Line:     lineNum,
				Key:      e.Key,
				Message:  "key should be upper-case",
				Severity: LintWarn,
			})
		}

		// Error on values that look like they contain unquoted whitespace-only strings
		if strings.TrimSpace(e.Value) != e.Value && e.Value != "" {
			issues = append(issues, LintIssue{
				Line:     lineNum,
				Key:      e.Key,
				Message:  "value has leading or trailing whitespace",
				Severity: LintError,
			})
		}

		// Warn on keys that start with a digit (unusual but technically valid)
		if len(e.Key) > 0 && e.Key[0] >= '0' && e.Key[0] <= '9' {
			issues = append(issues, LintIssue{
				Line:     lineNum,
				Key:      e.Key,
				Message:  "key starts with a digit",
				Severity: LintWarn,
			})
		}
	}

	return issues
}
