package envfile

import (
	"testing"
)

func TestDiffEntriesAdded(t *testing.T) {
	base := []Entry{{Key: "FOO", Value: "bar"}}
	next := []Entry{{Key: "FOO", Value: "bar"}, {Key: "NEW", Value: "val"}}

	changes := DiffEntries(base, next)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Kind != ChangeAdded || changes[0].Key != "NEW" || changes[0].NewValue != "val" {
		t.Errorf("unexpected change: %+v", changes[0])
	}
}

func TestDiffEntriesModified(t *testing.T) {
	base := []Entry{{Key: "FOO", Value: "bar"}}
	next := []Entry{{Key: "FOO", Value: "baz"}}

	changes := DiffEntries(base, next)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	c := changes[0]
	if c.Kind != ChangeModified || c.Key != "FOO" || c.OldValue != "bar" || c.NewValue != "baz" {
		t.Errorf("unexpected change: %+v", c)
	}
}

func TestDiffEntriesRemoved(t *testing.T) {
	base := []Entry{{Key: "FOO", Value: "bar"}, {Key: "OLD", Value: "gone"}}
	next := []Entry{{Key: "FOO", Value: "bar"}}

	changes := DiffEntries(base, next)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Kind != ChangeRemoved || changes[0].Key != "OLD" {
		t.Errorf("unexpected change: %+v", changes[0])
	}
}

func TestDiffEntriesNoChange(t *testing.T) {
	base := []Entry{{Key: "A", Value: "1"}, {Key: "B", Value: "2"}}
	changes := DiffEntries(base, base)
	if len(changes) != 0 {
		t.Errorf("expected no changes, got %d", len(changes))
	}
}

func TestPatchAppliesChanges(t *testing.T) {
	base := []Entry{{Key: "FOO", Value: "bar"}, {Key: "OLD", Value: "gone"}}
	next := []Entry{{Key: "FOO", Value: "baz"}, {Key: "NEW", Value: "val"}}

	changes := DiffEntries(base, next)
	result, err := Patch(base, changes)
	if err != nil {
		t.Fatalf("Patch returned error: %v", err)
	}

	m := make(map[string]string)
	for _, e := range result {
		m[e.Key] = e.Value
	}

	if m["FOO"] != "baz" {
		t.Errorf("expected FOO=baz, got %s", m["FOO"])
	}
	if m["NEW"] != "val" {
		t.Errorf("expected NEW=val, got %s", m["NEW"])
	}
	if _, ok := m["OLD"]; ok {
		t.Error("expected OLD to be removed")
	}
}

func TestPatchDoesNotMutateBase(t *testing.T) {
	base := []Entry{{Key: "X", Value: "1"}}
	changes := []Change{{Kind: ChangeModified, Key: "X", OldValue: "1", NewValue: "2"}}

	_, err := Patch(base, changes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if base[0].Value != "1" {
		t.Error("Patch mutated the base slice")
	}
}
