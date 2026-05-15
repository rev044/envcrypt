package envfile

import (
	"fmt"
	"time"
)

// QueryOptions filters history entries when calling Query.
type QueryOptions struct {
	Action  string    // if non-empty, only entries with this action
	Profile string    // if non-empty, only entries for this profile
	Since   time.Time // if non-zero, only entries at or after this time
}

// Query returns history entries matching all non-zero fields in opts.
func (h *History) Query(opts QueryOptions) ([]HistoryEntry, error) {
	all, err := h.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("history query: %w", err)
	}
	var result []HistoryEntry
	for _, e := range all {
		if opts.Action != "" && e.Action != opts.Action {
			continue
		}
		if opts.Profile != "" && e.Profile != opts.Profile {
			continue
		}
		if !opts.Since.IsZero() && e.Timestamp.Before(opts.Since) {
			continue
		}
		result = append(result, e)
	}
	return result, nil
}

// Latest returns the most recent history entry, or an error if history is empty.
func (h *History) Latest() (HistoryEntry, error) {
	all, err := h.ReadAll()
	if err != nil {
		return HistoryEntry{}, fmt.Errorf("history latest: %w", err)
	}
	if len(all) == 0 {
		return HistoryEntry{}, fmt.Errorf("history latest: no entries")
	}
	return all[len(all)-1], nil
}
