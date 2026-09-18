package format_test

import (
	"testing"

	"github.com/Formulary-Labs/substrate/format"
)

func TestParse_valid(t *testing.T) {
	for _, s := range []string{"json", "yaml", "md", "csv", "html"} {
		f, err := format.Parse(s)
		if err != nil {
			t.Errorf("Parse(%q) unexpected error: %v", s, err)
		}
		if string(f) != s {
			t.Errorf("Parse(%q) = %q, want %q", s, f, s)
		}
	}
}

func TestParse_invalid(t *testing.T) {
	_, err := format.Parse("pdf")
	if err == nil {
		t.Error("Parse(\"pdf\") expected error, got nil")
	}
}

func TestMIMEType(t *testing.T) {
	cases := map[format.Format]string{
		format.JSON: "application/json",
		format.YAML: "application/yaml",
		format.MD:   "text/markdown",
		format.CSV:  "text/csv",
		format.HTML: "text/html",
	}
	for f, want := range cases {
		if got := f.MIMEType(); got != want {
			t.Errorf("%q.MIMEType() = %q, want %q", f, got, want)
		}
	}
}

func TestExtension(t *testing.T) {
	cases := map[format.Format]string{
		format.JSON: "json",
		format.YAML: "yaml",
		format.MD:   "md",
		format.CSV:  "csv",
		format.HTML: "html",
	}
	for f, want := range cases {
		if got := f.Extension(); got != want {
			t.Errorf("%q.Extension() = %q, want %q", f, got, want)
		}
	}
}
