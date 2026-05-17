package envfile

import (
	"regexp"
	"testing"
)

func schemaEntries() []Entry {
	return []Entry{
		{Key: "APP_ENV", Value: "production"},
		{Key: "PORT", Value: "8080"},
		{Key: "DATABASE_URL", Value: "postgres://localhost/db"},
	}
}

func TestValidateSchemaAllPresent(t *testing.T) {
	s := Schema{Fields: []SchemaField{
		{Key: "APP_ENV", Required: true},
		{Key: "PORT", Required: true},
	}}
	violations := ValidateSchema(schemaEntries(), s)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %v", violations)
	}
}

func TestValidateSchemaMissingRequired(t *testing.T) {
	s := Schema{Fields: []SchemaField{
		{Key: "APP_ENV", Required: true},
		{Key: "SECRET_KEY", Required: true},
	}}
	violations := ValidateSchema(schemaEntries(), s)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d: %v", len(violations), violations)
	}
	if violations[0].Key != "SECRET_KEY" {
		t.Errorf("expected violation for SECRET_KEY, got %q", violations[0].Key)
	}
}

func TestValidateSchemaPatternMatch(t *testing.T) {
	portPattern := regexp.MustCompile(`^\d+$`)
	s := Schema{Fields: []SchemaField{
		{Key: "PORT", Required: true, Pattern: portPattern},
	}}
	violations := ValidateSchema(schemaEntries(), s)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %v", violations)
	}
}

func TestValidateSchemaPatternMismatch(t *testing.T) {
	portPattern := regexp.MustCompile(`^\d+$`)
	entries := []Entry{
		{Key: "PORT", Value: "not-a-number"},
	}
	s := Schema{Fields: []SchemaField{
		{Key: "PORT", Required: true, Pattern: portPattern},
	}}
	violations := ValidateSchema(entries, s)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestValidateSchemaOptionalMissing(t *testing.T) {
	s := Schema{Fields: []SchemaField{
		{Key: "OPTIONAL_KEY", Required: false},
	}}
	violations := ValidateSchema(schemaEntries(), s)
	if len(violations) != 0 {
		t.Fatalf("optional missing key should not produce violations, got %v", violations)
	}
}

func TestValidateSchemaEmptyValueTreatedAsMissing(t *testing.T) {
	entries := []Entry{
		{Key: "APP_ENV", Value: "   "},
	}
	s := Schema{Fields: []SchemaField{
		{Key: "APP_ENV", Required: true},
	}}
	violations := ValidateSchema(entries, s)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation for blank value, got %d", len(violations))
	}
}
