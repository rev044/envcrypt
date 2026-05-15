package envfile

import (
	"sort"
	"testing"
)

func TestMergePreferOverride(t *testing.T) {
	base := Env{"HOST": "localhost", "PORT": "5432"}
	override := Env{"PORT": "9999", "DEBUG": "true"}

	result, err := Merge(base, override, PreferOverride)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %q", result["HOST"])
	}
	if result["PORT"] != "9999" {
		t.Errorf("expected PORT=9999, got %q", result["PORT"])
	}
	if result["DEBUG"] != "true" {
		t.Errorf("expected DEBUG=true, got %q", result["DEBUG"])
	}
}

func TestMergePreferBase(t *testing.T) {
	base := Env{"HOST": "localhost", "PORT": "5432"}
	override := Env{"PORT": "9999", "DEBUG": "true"}

	result, err := Merge(base, override, PreferBase)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result["PORT"] != "5432" {
		t.Errorf("expected PORT=5432 (base wins), got %q", result["PORT"])
	}
	if result["DEBUG"] != "true" {
		t.Errorf("expected DEBUG=true (new key), got %q", result["DEBUG"])
	}
}

func TestMergeErrorOnConflict(t *testing.T) {
	base := Env{"HOST": "localhost"}
	override := Env{"HOST": "remote"}

	_, err := Merge(base, override, ErrorOnConflict)
	if err == nil {
		t.Fatal("expected error on conflict, got nil")
	}
}

func TestMergeNoConflict(t *testing.T) {
	base := Env{"A": "1"}
	override := Env{"B": "2"}

	result, err := Merge(base, override, ErrorOnConflict)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
}

func TestDiff(t *testing.T) {
	a := Env{"HOST": "localhost", "PORT": "5432", "ONLY_A": "x"}
	b := Env{"HOST": "localhost", "PORT": "9999", "ONLY_B": "y"}

	diffs := Diff(a, b)
	sort.Strings(diffs)

	expected := []string{"ONLY_A", "ONLY_B", "PORT"}
	if len(diffs) != len(expected) {
		t.Fatalf("expected diffs %v, got %v", expected, diffs)
	}
	for i, k := range expected {
		if diffs[i] != k {
			t.Errorf("expected diff[%d]=%q, got %q", i, k, diffs[i])
		}
	}
}

func TestDiffNoDifferences(t *testing.T) {
	a := Env{"X": "1"}
	b := Env{"X": "1"}

	if diffs := Diff(a, b); len(diffs) != 0 {
		t.Errorf("expected no diffs, got %v", diffs)
	}
}
