// Package provenance implements the provenance writer for all Formulary tools.
// All tool writes must log provenance by calling Write. The log is an
// append-only JSONL file at the path specified in the Entry.
//
// This re-implements the pattern from scripts/provenance_log.py so that Go
// tools produce provenance entries compatible with the Python agent layer.
//
// Entry format mirrors the Python provenance log schema so entries from both
// layers are readable by the same tooling.
package provenance

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

// Reusability indicates how reusable the produced artifact is.
type Reusability string

const (
	// Template means the artifact is a reusable template.
	Template Reusability = "template"
	// Reference means the artifact is a reference implementation.
	Reference Reusability = "reference"
	// Instance means the artifact is a one-off instance.
	Instance Reusability = "instance"
	// Artifact means the artifact is a primary deliverable.
	Artifact Reusability = "artifact"
)

// QualityGate indicates the quality gate result for the produced artifact.
type QualityGate string

const (
	Pass                QualityGate = "pass"
	FailedOnceCorrected QualityGate = "failed_once_corrected"
	Escalated           QualityGate = "escalated"
	Skipped             QualityGate = "skipped"
	NotApplicable       QualityGate = "not_applicable"
)

// Entry is a single provenance log record.
type Entry struct {
	// ID is a unique identifier for this log entry.
	ID string `json:"id"`
	// Timestamp is the ISO 8601 UTC timestamp of the log write.
	Timestamp string `json:"timestamp"`
	// Spec is the path to the spec that governed this output.
	Spec string `json:"spec"`
	// Output is the path to the produced artifact.
	Output string `json:"output"`
	// OutputType is the category of the produced artifact.
	OutputType string `json:"output_type"`
	// Program is the program slug associated with this output.
	Program string `json:"program"`
	// Purpose is a one-sentence description of why this output was produced.
	Purpose string `json:"purpose"`
	// Reusability indicates how reusable the produced artifact is.
	Reusability Reusability `json:"reusability"`
	// QualityGate indicates the quality gate result.
	QualityGate QualityGate `json:"quality_gate"`
	// RunID is an optional pipeline run identifier.
	RunID string `json:"run_id,omitempty"`
	// Tool identifies the Formulary tool that produced this entry.
	Tool string `json:"tool"`
	// ToolVersion is the version of the tool that produced this entry.
	ToolVersion string `json:"tool_version,omitempty"`
}

// Write appends a provenance entry to the JSONL log at logPath.
// If logPath is empty, it defaults to "logs/provenance.jsonl" relative to
// the current working directory.
//
// The directory is created if it does not exist. The file is created if it
// does not exist. Writes are appended, never truncated.
func Write(logPath string, e Entry) error {
	if logPath == "" {
		logPath = "logs/provenance.jsonl"
	}

	if e.ID == "" {
		e.ID = uuid.New().String()
	}
	if e.Timestamp == "" {
		e.Timestamp = time.Now().UTC().Format(time.RFC3339)
	}

	if err := os.MkdirAll(filepath.Dir(logPath), 0o750); err != nil {
		return fmt.Errorf("provenance: creating log directory: %w", err)
	}

	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("provenance: opening log file: %w", err)
	}
	defer f.Close() //nolint:errcheck

	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("provenance: marshaling entry: %w", err)
	}
	if _, err := fmt.Fprintf(f, "%s\n", data); err != nil {
		return fmt.Errorf("provenance: writing entry: %w", err)
	}
	return nil
}

// MustWrite calls Write and panics if it returns an error. Use only in tests
// or where a provenance failure is a hard error.
func MustWrite(logPath string, e Entry) {
	if err := Write(logPath, e); err != nil {
		panic(err)
	}
}
