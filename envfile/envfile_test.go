package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	return p
}

func TestParse(t *testing.T) {
	content := "# comment\nDB_HOST=localhost\nDB_PORT=5432\n"
	p := writeTempFile(t, content)

	ef, err := Parse(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ef.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(ef.Entries))
	}
	if ef.Entries[1].Key != "DB_HOST" || ef.Entries[1].Value != "localhost" {
		t.Errorf("unexpected entry: %+v", ef.Entries[1])
	}
}

func TestParseInvalidLine(t *testing.T) {
	p := writeTempFile(t, "BADLINE\n")
	_, err := Parse(p)
	if err == nil {
		t.Fatal("expected error for invalid line")
	}
}

func TestWriteRoundTrip(t *testing.T) {
	original := "# header\nFOO=bar\nBAZ=qux\n"
	p := writeTempFile(t, original)

	ef, err := Parse(p)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	out := filepath.Join(t.TempDir(), ".env.out")
	if err := Write(out, ef); err != nil {
		t.Fatalf("write: %v", err)
	}

	ef2, err := Parse(out)
	if err != nil {
		t.Fatalf("re-parse: %v", err)
	}
	if len(ef.Entries) != len(ef2.Entries) {
		t.Errorf("entry count mismatch: %d vs %d", len(ef.Entries), len(ef2.Entries))
	}
}

func TestToMapFromMap(t *testing.T) {
	m := map[string]string{"KEY1": "val1", "KEY2": "val2"}
	ef := FromMap(m)
	got := ef.ToMap()
	for k, v := range m {
		if got[k] != v {
			t.Errorf("key %s: expected %q got %q", k, v, got[k])
		}
	}
}
