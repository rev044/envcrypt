package envfile

// ChangeKind describes the type of difference between two env files.
type ChangeKind string

const (
	ChangeAdded    ChangeKind = "added"
	ChangeModified ChangeKind = "modified"
	ChangeRemoved  ChangeKind = "removed"
)

// Change represents a single key-level difference.
type Change struct {
	Kind     ChangeKind
	Key      string
	OldValue string // empty for ChangeAdded
	NewValue string // empty for ChangeRemoved
}

// DiffEntries returns the ordered list of Changes needed to transform
// base into next.
func DiffEntries(base, next []Entry) []Change {
	baseMap := make(map[string]string, len(base))
	for _, e := range base {
		baseMap[e.Key] = e.Value
	}

	nextMap := make(map[string]string, len(next))
	for _, e := range next {
		nextMap[e.Key] = e.Value
	}

	var changes []Change

	// Added or modified
	for _, e := range next {
		if old, ok := baseMap[e.Key]; !ok {
			changes = append(changes, Change{Kind: ChangeAdded, Key: e.Key, NewValue: e.Value})
		} else if old != e.Value {
			changes = append(changes, Change{Kind: ChangeModified, Key: e.Key, OldValue: old, NewValue: e.Value})
		}
	}

	// Removed
	for _, e := range base {
		if _, ok := nextMap[e.Key]; !ok {
			changes = append(changes, Change{Kind: ChangeRemoved, Key: e.Key, OldValue: e.Value})
		}
	}

	return changes
}

// Patch applies the given Changes to base and returns the resulting entries.
// The original slice is not modified.
func Patch(base []Entry, changes []Change) ([]Entry, error) {
	result := make([]Entry, len(base))
	copy(result, base)

	for _, c := range changes {
		switch c.Kind {
		case ChangeAdded:
			result = append(result, Entry{Key: c.Key, Value: c.NewValue})
		case ChangeModified:
			for i, e := range result {
				if e.Key == c.Key {
					result[i].Value = c.NewValue
					break
				}
			}
		case ChangeRemoved:
			filtered := result[:0]
			for _, e := range result {
				if e.Key != c.Key {
					filtered = append(filtered, e)
				}
			}
			result = filtered
		}
	}

	return result, nil
}
