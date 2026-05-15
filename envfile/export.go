package envfile

import (
	"fmt"
	"io"
	"strings"
)

// ExportFormat defines the output format for environment variable export.
type ExportFormat int

const (
	// FormatShell emits "export KEY=VALUE" lines suitable for shell sourcing.
	FormatShell ExportFormat = iota
	// FormatDocker emits "KEY=VALUE" lines suitable for Docker --env-file.
	FormatDocker
	// FormatJSON emits a JSON object of key/value pairs.
	FormatJSON
)

// Export writes the entries to w in the requested format.
func Export(entries []Entry, format ExportFormat, w io.Writer) error {
	switch format {
	case FormatShell:
		return exportShell(entries, w)
	case FormatDocker:
		return exportDocker(entries, w)
	case FormatJSON:
		return exportJSON(entries, w)
	default:
		return fmt.Errorf("envfile: unknown export format %d", format)
	}
}

func exportShell(entries []Entry, w io.Writer) error {
	for _, e := range entries {
		if e.Comment != "" {
			if _, err := fmt.Fprintf(w, "# %s\n", e.Comment); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintf(w, "export %s=%s\n", e.Key, shellQuote(e.Value)); err != nil {
			return err
		}
	}
	return nil
}

func exportDocker(entries []Entry, w io.Writer) error {
	for _, e := range entries {
		if _, err := fmt.Fprintf(w, "%s=%s\n", e.Key, e.Value); err != nil {
			return err
		}
	}
	return nil
}

func exportJSON(entries []Entry, w io.Writer) error {
	var sb strings.Builder
	sb.WriteString("{\n")
	for i, e := range entries {
		comma := ","
		if i == len(entries)-1 {
			comma = ""
		}
		sb.WriteString(fmt.Sprintf("  %q: %q%s\n", e.Key, e.Value, comma))
	}
	sb.WriteString("}\n")
	_, err := io.WriteString(w, sb.String())
	return err
}

// shellQuote wraps value in single quotes if it contains special characters.
func shellQuote(v string) string {
	specials := " \t\n\r\"'\\$`!#&;|<>(){}"
	for _, c := range specials {
		if strings.ContainsRune(v, c) {
			return "'" + strings.ReplaceAll(v, "'", "'\\'''") + "'"
		}
	}
	return v
}
