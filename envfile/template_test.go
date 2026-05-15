package envfile

import (
	"os"
	"testing"
)

func entries(pairs ...string) []Entry {
	var out []Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestRenderSimpleSubstitution(t *testing.T) {
	in := entries("BASE", "/app", "DATA", "${BASE}/data")
	out, err := Render(in, RenderOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[1].Value != "/app/data" {
		t.Errorf("got %q, want %q", out[1].Value, "/app/data")
	}
}

func TestRenderDollarSyntax(t *testing.T) {
	in := entries("HOME", "/home/user", "CONF", "$HOME/.config")
	out, err := Render(in, RenderOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[1].Value != "/home/user/.config" {
		t.Errorf("got %q, want %q", out[1].Value, "/home/user/.config")
	}
}

func TestRenderChained(t *testing.T) {
	in := entries("A", "hello", "B", "${A} world", "C", "${B}!")
	out, err := Render(in, RenderOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[2].Value != "hello world!" {
		t.Errorf("got %q, want %q", out[2].Value, "hello world!")
	}
}

func TestRenderMissingVarError(t *testing.T) {
	in := entries("X", "${UNDEFINED}")
	_, err := Render(in, RenderOptions{})
	if err == nil {
		t.Fatal("expected error for undefined variable, got nil")
	}
}

func TestRenderAllowMissing(t *testing.T) {
	in := entries("X", "${UNDEFINED}")
	out, err := Render(in, RenderOptions{AllowMissing: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Value != "${UNDEFINED}" {
		t.Errorf("got %q, want %q", out[0].Value, "${UNDEFINED}")
	}
}

func TestRenderFallbackToEnv(t *testing.T) {
	t.Setenv("MY_HOST", "localhost")
	_ = os.Setenv("MY_HOST", "localhost")
	in := entries("ADDR", "${MY_HOST}:8080")
	out, err := Render(in, RenderOptions{FallbackToEnv: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Value != "localhost:8080" {
		t.Errorf("got %q, want %q", out[0].Value, "localhost:8080")
	}
}

func TestRenderPreservesComments(t *testing.T) {
	in := []Entry{{Key: "FOO", Value: "bar", Comment: "# a comment"}}
	out, err := Render(in, RenderOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Comment != "# a comment" {
		t.Errorf("comment not preserved: got %q", out[0].Comment)
	}
}
