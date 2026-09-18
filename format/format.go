// Package format defines the output format types and helpers for all
// Formulary tools.
//
// Every tool must support at minimum Format JSON and Format MD. Tools that
// produce tabular data should also support CSV. Tools that produce documents
// support HTML. The default format is JSON — machine-readable output first.
package format

import "fmt"

// Format is an output format identifier.
type Format string

const (
	// JSON is the default machine-readable output format.
	// All tool JSON output must be parseable by jq and by the next tool in
	// a pipeline.
	JSON Format = "json"

	// YAML is an alternative machine-readable output format.
	YAML Format = "yaml"

	// MD is a human-readable markdown output format.
	MD Format = "md"

	// CSV is a tabular output format for spreadsheet compatibility.
	CSV Format = "csv"

	// HTML is a static HTML output format for browser or PDF rendering.
	// Used by exhibit and vital.
	HTML Format = "html"
)

// Parse parses a format string, returning an error for unrecognised values.
func Parse(s string) (Format, error) {
	switch Format(s) {
	case JSON, YAML, MD, CSV, HTML:
		return Format(s), nil
	default:
		return "", fmt.Errorf("unsupported format %q: must be one of json, yaml, md, csv, html", s)
	}
}

// Valid returns true if the format string is a known Format value.
func Valid(s string) bool {
	_, err := Parse(s)
	return err == nil
}

// MIMEType returns the MIME type for the format, for use in HTTP headers and
// content type metadata.
func (f Format) MIMEType() string {
	switch f {
	case JSON:
		return "application/json"
	case YAML:
		return "application/yaml"
	case MD:
		return "text/markdown"
	case CSV:
		return "text/csv"
	case HTML:
		return "text/html"
	default:
		return "application/octet-stream"
	}
}

// Extension returns the file extension for the format, without a leading dot.
func (f Format) Extension() string {
	switch f {
	case JSON:
		return "json"
	case YAML:
		return "yaml"
	case MD:
		return "md"
	case CSV:
		return "csv"
	case HTML:
		return "html"
	default:
		return "bin"
	}
}
