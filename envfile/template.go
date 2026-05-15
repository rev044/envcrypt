package envfile

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// templateVarRe matches ${VAR_NAME} or $VAR_NAME style placeholders.
var templateVarRe = regexp.MustCompile(`\$\{([A-Z_][A-Z0-9_]*)\}|\$([A-Z_][A-Z0-9_]*)`)

// RenderOptions controls how template rendering behaves.
type RenderOptions struct {
	// AllowMissing suppresses errors for unresolved variables.
	AllowMissing bool
	// FallbackToEnv uses the process environment for missing variables.
	FallbackToEnv bool
}

// Render resolves variable references within entry values using the provided
// vars map as the source of substitutions. Entries are processed in order;
// earlier entries may be referenced by later ones.
//
// Example: given BASE=/app and PATH=${BASE}/bin, PATH resolves to /app/bin.
func Render(entries []Entry, opts RenderOptions) ([]Entry, error) {
	vars := make(map[string]string, len(entries))
	for _, e := range entries {
		vars[e.Key] = e.Value
	}

	out := make([]Entry, 0, len(entries))
	for _, e := range entries {
		resolved, err := renderValue(e.Value, vars, opts)
		if err != nil {
			return nil, fmt.Errorf("envfile: render %q: %w", e.Key, err)
		}
		// Update vars so subsequent entries can reference the resolved value.
		vars[e.Key] = resolved
		out = append(out, Entry{Key: e.Key, Value: resolved, Comment: e.Comment})
	}
	return out, nil
}

func renderValue(value string, vars map[string]string, opts RenderOptions) (string, error) {
	var renderErr error
	result := templateVarRe.ReplaceAllStringFunc(value, func(match string) string {
		if renderErr != nil {
			return match
		}
		sub := templateVarRe.FindStringSubmatch(match)
		name := sub[1]
		if name == "" {
			name = sub[2]
		}
		if v, ok := vars[name]; ok {
			return v
		}
		if opts.FallbackToEnv {
			if v, ok := os.LookupEnv(name); ok {
				return v
			}
		}
		if opts.AllowMissing {
			return match
		}
		renderErr = fmt.Errorf("undefined variable %q", name)
		return match
	})
	if renderErr != nil {
		return "", renderErr
	}
	return strings.TrimSpace(result), nil
}
