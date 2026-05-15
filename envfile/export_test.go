package envfile

import (
	"strings"
	"testing"
)

func TestExportShell(t *testing.T) {
	entries := []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_URL", Value: "postgres://localhost/db"},
	}
	var sb strings.Builder
	if err := Export(entries, FormatShell, &sb); err != nil {
		t.Fatalf("Export shell: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "export APP_NAME=myapp") {
		t.Errorf("expected export APP_NAME=myapp in output, got:\n%s", out)
	}
	if !strings.Contains(out, "export DB_URL=") {
		t.Errorf("expected export DB_URL= in output, got:\n%s", out)
	}
}

func TestExportShellQuotesSpecial(t *testing.T) {
	entries := []Entry{
		{Key: "MSG", Value: "hello world"},
	}
	var sb strings.Builder
	if err := Export(entries, FormatShell, &sb); err != nil {
		t.Fatalf("Export shell: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "export MSG='hello world'") {
		t.Errorf("expected quoted value, got:\n%s", out)
	}
}

func TestExportDocker(t *testing.T) {
	entries := []Entry{
		{Key: "PORT", Value: "8080"},
		{Key: "ENV", Value: "production"},
	}
	var sb strings.Builder
	if err := Export(entries, FormatDocker, &sb); err != nil {
		t.Fatalf("Export docker: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "PORT=8080") {
		t.Errorf("expected PORT=8080, got:\n%s", out)
	}
	if strings.Contains(out, "export ") {
		t.Errorf("docker format should not contain 'export', got:\n%s", out)
	}
}

func TestExportJSON(t *testing.T) {
	entries := []Entry{
		{Key: "FOO", Value: "bar"},
	}
	var sb strings.Builder
	if err := Export(entries, FormatJSON, &sb); err != nil {
		t.Fatalf("Export JSON: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, `"FOO": "bar"`) {
		t.Errorf("expected JSON key/value, got:\n%s", out)
	}
}

func TestExportUnknownFormat(t *testing.T) {
	var sb strings.Builder
	err := Export(nil, ExportFormat(99), &sb)
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func TestExportShellComment(t *testing.T) {
	entries := []Entry{
		{Key: "API_KEY", Value: "secret", Comment: "API credentials"},
	}
	var sb strings.Builder
	if err := Export(entries, FormatShell, &sb); err != nil {
		t.Fatalf("Export shell: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "# API credentials") {
		t.Errorf("expected comment in output, got:\n%s", out)
	}
}
