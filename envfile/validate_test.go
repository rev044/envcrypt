package envfile

import (
	"testing"
)

func TestValidateClean(t *testing.T) {
	entries := []Entry{
		{Key: "APP_ENV", Value: "production"},
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "_PRIVATE", Value: "secret"},
	}
	result := Validate(entries)
	if !result.Valid() {
		t.Fatalf("expected valid, got errors: %s", result.Error())
	}
}

func TestValidateEmptyKey(t *testing.T) {
	entries := []Entry{
		{Key: "", Value: "oops"},
	}
	result := Validate(entries)
	if result.Valid() {
		t.Fatal("expected invalid due to empty key")
	}
	if len(result.Errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(result.Errors))
	}
}

func TestValidateInvalidKeyChars(t *testing.T) {
	entries := []Entry{
		{Key: "123BAD", Value: "value"},
		{Key: "has-hyphen", Value: "value"},
	}
	result := Validate(entries)
	if result.Valid() {
		t.Fatal("expected invalid due to bad key characters")
	}
	if len(result.Errors) != 2 {
		t.Fatalf("expected 2 errors, got %d", len(result.Errors))
	}
}

func TestValidateDuplicateKeys(t *testing.T) {
	entries := []Entry{
		{Key: "FOO", Value: "first"},
		{Key: "BAR", Value: "other"},
		{Key: "FOO", Value: "second"},
	}
	result := Validate(entries)
	if result.Valid() {
		t.Fatal("expected invalid due to duplicate key")
	}
	if len(result.Errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(result.Errors))
	}
	if result.Errors[0].Key != "FOO" {
		t.Errorf("expected error on key FOO, got %q", result.Errors[0].Key)
	}
}

func TestValidateMultipleIssues(t *testing.T) {
	entries := []Entry{
		{Key: "", Value: "empty"},
		{Key: "bad-key", Value: "x"},
		{Key: "GOOD", Value: "ok"},
		{Key: "GOOD", Value: "dup"},
	}
	result := Validate(entries)
	if result.Valid() {
		t.Fatal("expected invalid")
	}
	if len(result.Errors) != 3 {
		t.Fatalf("expected 3 errors, got %d: %s", len(result.Errors), result.Error())
	}
}
