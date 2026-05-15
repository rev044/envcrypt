package envfile

// Entry represents a single key/value pair parsed from a .env file.
// An optional inline comment may be associated with the entry.
type Entry struct {
	// Key is the environment variable name.
	Key string
	// Value is the raw (unquoted) value of the variable.
	Value string
	// Comment holds any inline or preceding comment text stripped of the
	// leading "#" and surrounding whitespace.
	Comment string
}

// ToMap converts a slice of Entry values into a plain string map.
// Duplicate keys are resolved by keeping the last occurrence.
func ToMap(entries []Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return m
}
