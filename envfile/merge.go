package envfile

import "fmt"

// MergeStrategy defines how conflicts are resolved when merging two Env maps.
type MergeStrategy int

const (
	// PreferBase keeps the value from the base map on conflict.
	PreferBase MergeStrategy = iota
	// PreferOverride replaces the base value with the override value on conflict.
	PreferOverride
	// ErrorOnConflict returns an error if the same key exists in both maps.
	ErrorOnConflict
)

// Merge combines base and override Env maps according to the given strategy.
// Keys present only in override are always added to the result.
func Merge(base, override Env, strategy MergeStrategy) (Env, error) {
	result := make(Env, len(base))
	for k, v := range base {
		result[k] = v
	}

	for k, v := range override {
		if _, exists := result[k]; exists {
			switch strategy {
			case PreferBase:
				// keep existing value — do nothing
			case PreferOverride:
				result[k] = v
			case ErrorOnConflict:
				return nil, fmt.Errorf("envfile: merge conflict on key %q", k)
			}
		} else {
			result[k] = v
		}
	}

	return result, nil
}

// Diff returns the keys that differ between a and b.
// A key is considered different if it is missing from one map or has a
// different value.
func Diff(a, b Env) []string {
	seen := make(map[string]struct{})
	var diffs []string

	for k, av := range a {
		seen[k] = struct{}{}
		if bv, ok := b[k]; !ok || av != bv {
			diffs = append(diffs, k)
		}
	}

	for k := range b {
		if _, ok := seen[k]; !ok {
			diffs = append(diffs, k)
		}
	}

	return diffs
}
