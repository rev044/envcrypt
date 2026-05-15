package envfile

import (
	"testing"
)

func TestToMap(t *testing.T) {
	entries := []Entry{
		{Key: "A", Value: "1"},
		{Key: "B", Value: "2"},
		{Key: "A", Value: "3"}, // duplicate — last wins
	}
	m := ToMap(entries)
	if m["A"] != "3" {
		t.Errorf("expected A=3, got %s", m["A"])
	}
	if m["B"] != "2" {
		t.Errorf("expected B=2, got %s", m["B"])
	}
	if len(m) != 2 {
		t.Errorf("expected 2 keys, got %d", len(m))
	}
}

func TestToMapEmpty(t *testing.T) {
	m := ToMap(nil)
	if len(m) != 0 {
		t.Errorf("expected empty map, got %d entries", len(m))
	}
}

func TestEntryFields(t *testing.T) {
	e := Entry{Key: "FOO", Value: "bar", Comment: "a comment"}
	if e.Key != "FOO" {
		t.Errorf("Key mismatch")
	}
	if e.Value != "bar" {
		t.Errorf("Value mismatch")
	}
	if e.Comment != "a comment" {
		t.Errorf("Comment mismatch")
	}
}
