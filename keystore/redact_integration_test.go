package keystore_test

import (
	"testing"

	"github.com/yourorg/envcrypt/envfile"
)

// TestRedactIntegrationWithKeystore verifies that after loading entries from
// a simulated decrypted env, sensitive fields are properly masked when
// passed through the redact pipeline.
func TestRedactIntegrationWithKeystore(t *testing.T) {
	entries := []envfile.Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PASSWORD", Value: "hunter2"},
		{Key: "API_KEY", Value: "abc-def-ghi"},
		{Key: "APP_ENV", Value: "production"},
		{Key: "AUTH_TOKEN", Value: "tok_live_xyz"},
	}

	redacted := envfile.Redact(entries, envfile.RedactOptions{Mask: "****"})

	expectClear := map[string]string{
		"DB_HOST": "localhost",
		"APP_ENV": "production",
	}
	expectMasked := []string{"DB_PASSWORD", "API_KEY", "AUTH_TOKEN"}

	for _, e := range redacted {
		if want, ok := expectClear[e.Key]; ok {
			if e.Value != want {
				t.Errorf("key %q: expected %q, got %q", e.Key, want, e.Value)
			}
		}
	}

	for _, e := range redacted {
		for _, k := range expectMasked {
			if e.Key == k && e.Value != "****" {
				t.Errorf("key %q should be masked, got %q", e.Key, e.Value)
			}
		}
	}
}

func TestRedactShowPrefixIntegration(t *testing.T) {
	entries := []envfile.Entry{
		{Key: "DB_PASSWORD", Value: "hunter2"},
	}
	redacted := envfile.Redact(entries, envfile.RedactOptions{ShowPrefix: 2})
	got := redacted[0].Value
	if got[:2] != "hu" {
		t.Errorf("expected prefix 'hu', got %q", got)
	}
	if envfile.IsSensitive("DB_PASSWORD") == false {
		t.Error("DB_PASSWORD should be sensitive")
	}
}
