package artifact

import (
	"bytes"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Format is the on-disk interchange language of an artifact file.
type Format string

const (
	// FormatGemara is a gemara document (root metadata.type).
	FormatGemara Format = "gemara"
	// FormatOSCAL is a native OSCAL document (single model root).
	FormatOSCAL Format = "oscal"
	// FormatUnknown means the file parsed but matched neither language.
	FormatUnknown Format = "unknown"
)

// oscalModelRoots are the OSCAL complete-schema model names. A native OSCAL
// document has exactly one of these as its root key.
var oscalModelRoots = map[string]struct{}{
	"catalog":                       {},
	"profile":                       {},
	"component-definition":          {},
	"system-security-plan":          {},
	"assessment-plan":               {},
	"assessment-results":            {},
	"plan-of-action-and-milestones": {},
	"mapping-collection":            {},
}

// DetectFormat reports whether path is gemara, OSCAL, or neither.
// It does not schema-validate. Callers that need a pass/fail check should
// load or validate after routing.
func DetectFormat(path string) (Format, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return FormatUnknown, fmt.Errorf("reading artifact at %q: %w", path, err)
	}
	format, err := DetectFormatBytes(data)
	if err != nil {
		return FormatUnknown, fmt.Errorf("detecting format in %q: %w", path, err)
	}
	return format, nil
}

// DetectFormatBytes is DetectFormat for an already-read document.
func DetectFormatBytes(data []byte) (Format, error) {
	if len(bytes.TrimSpace(data)) == 0 {
		return FormatUnknown, fmt.Errorf("artifact is empty")
	}
	var root map[string]any
	if err := yaml.Unmarshal(data, &root); err != nil {
		return FormatUnknown, fmt.Errorf("parsing artifact: %w", err)
	}
	if root == nil {
		return FormatUnknown, fmt.Errorf("artifact is empty")
	}
	// Gemara wins when both shapes could be argued: its type lives on the
	// root metadata object, which OSCAL never uses as a sibling of the model.
	if gemaraRoot(root) {
		return FormatGemara, nil
	}
	if oscalRoot(root) {
		return FormatOSCAL, nil
	}
	return FormatUnknown, nil
}

func gemaraRoot(root map[string]any) bool {
	meta, ok := root["metadata"].(map[string]any)
	if !ok {
		return false
	}
	_, ok = meta["type"]
	return ok
}

func oscalRoot(root map[string]any) bool {
	if len(root) != 1 {
		return false
	}
	for key := range root {
		_, ok := oscalModelRoots[key]
		return ok
	}
	return false
}
