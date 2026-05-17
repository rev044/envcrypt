package envfile

import (
	"testing"
)

func TestLintClean(t *testing.T) {
	entries := []Entry{
		{Key: "DATABASE_URL", Value: "postgres://localhost/db"},
		{Key: "PORT", Value: "8080"},
	}
	issues := Lint(entries)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d: %v", len(issues), issues)
	}
}

func TestLintEmptyValue(t *testing.T) {
	entries := []Entry{
		{Key: "SECRET", Value: ""},
	}
	issues := Lint(entries)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Severity != LintWarn {
		t.Errorf("expected warn severity, got %s", issues[0].Severity)
	}
	if issues[0].Key != "SECRET" {
		t.Errorf("expected key SECRET, got %s", issues[0].Key)
	}
}

func TestLintLowercaseKey(t *testing.T) {
	entries := []Entry{
		{Key: "db_host", Value: "localhost"},
	}
	issues := Lint(entries)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Severity != LintWarn {
		t.Errorf("expected warn, got %s", issues[0].Severity)
	}
}

func TestLintLeadingTrailingWhitespace(t *testing.T) {
	entries := []Entry{
		{Key: "API_KEY", Value: " secret "},
	}
	issues := Lint(entries)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Severity != LintError {
		t.Errorf("expected error severity, got %s", issues[0].Severity)
	}
}

func TestLintKeyStartsWithDigit(t *testing.T) {
	entries := []Entry{
		{Key: "1SECRET", Value: "value"},
	}
	issues := Lint(entries)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Severity != LintWarn {
		t.Errorf("expected warn, got %s", issues[0].Severity)
	}
}

func TestLintMultipleIssues(t *testing.T) {
	entries := []Entry{
		{Key: "bad_key", Value: " "},
	}
	issues := Lint(entries)
	// lowercase key (warn) + whitespace-only value counts as leading/trailing whitespace (error)
	if len(issues) < 2 {
		t.Fatalf("expected at least 2 issues, got %d: %v", len(issues), issues)
	}
}

func TestLintIssueString(t *testing.T) {
	issue := LintIssue{Line: 3, Key: "FOO", Message: "value is empty", Severity: LintWarn}
	s := issue.String()
	if s == "" {
		t.Error("expected non-empty string representation")
	}
}
