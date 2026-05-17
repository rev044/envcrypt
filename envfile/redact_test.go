package envfile

import (
	"testing"
)

func TestIsSensitive(t *testing.T) {
	cases := []struct {
		key       string
		expected  bool
	}{
		{"DB_PASSWORD", true},
		{"API_KEY", true},
		{"AUTH_TOKEN", true},
		{"SECRET", true},
		{"PRIVATE_KEY", true},
		{"APP_CREDENTIAL", true},
		{"DATABASE_URL", false},
		{"PORT", false},
		{"APP_NAME", false},
	}
	for _, tc := range cases {
		t.Run(tc.key, func(t *testing.T) {
			if got := IsSensitive(tc.key); got != tc.expected {
				t.Errorf("IsSensitive(%q) = %v, want %v", tc.key, got, tc.expected)
			}
		})
	}
}

func TestRedactMasksSensitiveValues(t *testing.T) {
	entries := []Entry{
		{Key: "DB_PASSWORD", Value: "supersecret"},
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "API_KEY", Value: "abc123"},
	}
	redacted := Redact(entries, RedactOptions{})
	if redacted[0].Value != "****" {
		t.Errorf("expected DB_PASSWORD to be masked, got %q", redacted[0].Value)
	}
	if redacted[1].Value != "myapp" {
		t.Errorf("expected APP_NAME unchanged, got %q", redacted[1].Value)
	}
	if redacted[2].Value != "****" {
		t.Errorf("expected API_KEY to be masked, got %q", redacted[2].Value)
	}
}

func TestRedactShowPrefix(t *testing.T) {
	entries := []Entry{
		{Key: "DB_PASSWORD", Value: "supersecret"},
	}
	redacted := Redact(entries, RedactOptions{ShowPrefix: 3})
	got := redacted[0].Value
	if got[:3] != "sup" {
		t.Errorf("expected prefix 'sup', got %q", got)
	}
	if len(got) != len("supersecret") {
		t.Errorf("expected same length as original, got %d", len(got))
	}
}

func TestRedactCustomMask(t *testing.T) {
	entries := []Entry{
		{Key: "SECRET", Value: "topsecret"},
	}
	redacted := Redact(entries, RedactOptions{Mask: "[REDACTED]"})
	if redacted[0].Value != "[REDACTED]" {
		t.Errorf("expected custom mask, got %q", redacted[0].Value)
	}
}

func TestRedactPreservesComments(t *testing.T) {
	entries := []Entry{
		{Key: "DB_PASSWORD", Value: "secret", Comment: "# db creds"},
	}
	redacted := Redact(entries, RedactOptions{})
	if redacted[0].Comment != "# db creds" {
		t.Errorf("expected comment preserved, got %q", redacted[0].Comment)
	}
}
