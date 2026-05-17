package envfile

import (
	"fmt"
	"regexp"
	"strings"
)

// SchemaField defines the expected shape of a single env variable.
type SchemaField struct {
	Key      string
	Required bool
	Pattern  *regexp.Regexp // optional value pattern
}

// Schema holds a collection of field definitions used to validate an env file.
type Schema struct {
	Fields []SchemaField
}

// SchemaViolation describes a single schema validation failure.
type SchemaViolation struct {
	Key     string
	Message string
}

func (v SchemaViolation) Error() string {
	return fmt.Sprintf("schema violation for %q: %s", v.Key, v.Message)
}

// ValidateSchema checks that entries satisfy the schema.
// It returns all violations found; a nil slice means the entries are valid.
func ValidateSchema(entries []Entry, s Schema) []SchemaViolation {
	present := make(map[string]string, len(entries))
	for _, e := range entries {
		if e.Key != "" {
			present[e.Key] = e.Value
		}
	}

	var violations []SchemaViolation

	for _, field := range s.Fields {
		val, exists := present[field.Key]
		if !exists || strings.TrimSpace(val) == "" {
			if field.Required {
				violations = append(violations, SchemaViolation{
					Key:     field.Key,
					Message: "required key is missing or empty",
				})
			}
			continue
		}

		if field.Pattern != nil && !field.Pattern.MatchString(val) {
			violations = append(violations, SchemaViolation{
				Key:     field.Key,
				Message: fmt.Sprintf("value %q does not match required pattern %s", val, field.Pattern.String()),
			})
		}
	}

	return violations
}
