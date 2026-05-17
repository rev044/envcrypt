package envfile

import (
	"regexp"
	"strings"
)

// RedactOptions controls how sensitive values are masked.
type RedactOptions struct {
	// ShowPrefix reveals the first N characters of the value before masking.
	ShowPrefix int
	// Mask is the string used to replace sensitive content. Defaults to "****".
	Mask string
}

var sensitiveKeyPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(password|passwd|pwd)`),
	regexp.MustCompile(`(?i)(secret|token|apikey|api_key)`),
	regexp.MustCompile(`(?i)(private_key|privkey)`),
	regexp.MustCompile(`(?i)(auth|credential)`),
}

// IsSensitive reports whether the given key name looks like it holds
// a sensitive value based on common naming conventions.
func IsSensitive(key string) bool {
	for _, re := range sensitiveKeyPatterns {
		if re.MatchString(key) {
			return true
		}
	}
	return false
}

// Redact returns a copy of entries where values whose keys match sensitive
// patterns are replaced with a mask string.
func Redact(entries []Entry, opts RedactOptions) []Entry {
	if opts.Mask == "" {
		opts.Mask = "****"
	}
	out := make([]Entry, len(entries))
	for i, e := range entries {
		if IsSensitive(e.Key) {
			masked := opts.Mask
			if opts.ShowPrefix > 0 && len(e.Value) > opts.ShowPrefix {
				masked = e.Value[:opts.ShowPrefix] + strings.Repeat("*", len(e.Value)-opts.ShowPrefix)
			}
			out[i] = Entry{Key: e.Key, Value: masked, Comment: e.Comment}
		} else {
			out[i] = e
		}
	}
	return out
}
