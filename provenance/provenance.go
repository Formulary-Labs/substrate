// Package provenance implements the provenance writer for all Formulary tools.
// All tool writes must log provenance by calling Write. The log is an
// append-only JSONL file at the path specified in the Entry.
//
// This re-implements the pattern from scripts/provenance_log.py so that Go
// tools produce provenance entries compatible with the Python agent layer.
//
// Entry format mirrors the Python provenance log schema so entries from both
// layers are readable by the same tooling.
//
// Hash-chain integrity: each entry carries a SHA-256 Digest of its own
// canonical JSON (with the Digest field zeroed) and a PreviousDigest linking
// it to the prior entry. Call Verify to validate the chain before relying
// on the log as evidence.
package provenance

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
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

//nolint:revive // QualityGate constants are self-documenting string identifiers.
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

	// Hash-chain fields. PreviousDigest links this entry to its predecessor.
	// Digest is the SHA-256 hex of the canonical JSON of this entry with
	// the Digest field zeroed out.
	PreviousDigest string `json:"previous_digest,omitempty"`
	Digest         string `json:"digest,omitempty"`
}

// Write appends a provenance entry to the JSONL log at logPath.
// If logPath is empty, it defaults to "logs/provenance.jsonl" relative to
// the current working directory.
//
// The directory is created if it does not exist. The file is created if it
// does not exist. Writes are appended, never truncated.
//
// Each entry is cryptographically chained: PreviousDigest is set to the
// stored Digest of the prior entry, and Digest is computed as the SHA-256
// of this entry's canonical JSON.
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

	// Link to the prior entry for chain integrity.
	e.PreviousDigest = lastDigest(logPath)

	// Compute this entry's digest using canonical JSON with Digest zeroed.
	e.Digest = ""
	canonical, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("provenance: computing digest: %w", err)
	}
	sum := sha256.Sum256(canonical)
	e.Digest = hex.EncodeToString(sum[:])

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

// lastDigest reads the last line of the JSONL file and returns its "digest"
// field. Returns "" if the file does not exist or has no chained entries.
func lastDigest(logPath string) string {
	f, err := os.Open(logPath)
	if err != nil {
		return ""
	}
	defer f.Close() //nolint:errcheck

	var last string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if line := scanner.Text(); line != "" {
			last = line
		}
	}
	if last == "" {
		return ""
	}
	var e struct {
		Digest string `json:"digest"`
	}
	if err := json.Unmarshal([]byte(last), &e); err != nil {
		return ""
	}
	return e.Digest
}

// MustWrite calls Write and panics if it returns an error. Use only in tests
// or where a provenance failure is a hard error.
func MustWrite(logPath string, e Entry) {
	if err := Write(logPath, e); err != nil {
		panic(err)
	}
}

// VerificationResult holds the integrity check outcome for a single log entry.
type VerificationResult struct {
	// EntryID is the UUID of the checked entry.
	EntryID string `json:"entry_id"`
	// EntryIndex is the zero-based position of the entry in the log.
	EntryIndex int `json:"entry_index"`
	// Valid is true when the stored Digest matches the recomputed SHA-256.
	Valid bool `json:"valid"`
	// ChainValid is true when PreviousDigest matches the prior entry's Digest.
	ChainValid bool `json:"chain_valid"`
	// FailureReason describes the first integrity failure found, if any.
	FailureReason string `json:"failure_reason,omitempty"`
}

// Verify reads the provenance JSONL at logPath and validates the hash chain.
// Each entry's Digest is recomputed and compared to its stored value.
// PreviousDigest linkage between consecutive entries is also checked.
//
// Returns nil results (not an error) for a non-existent or empty log.
// Legacy entries that predate hash-chain support (Digest == "") are flagged
// as invalid so callers can decide how to handle mixed logs.
func Verify(logPath string) ([]VerificationResult, error) {
	data, err := os.ReadFile(logPath)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("provenance: reading log for verification: %w", err)
	}

	var results []VerificationResult
	var prevDigest string

	scanner := bufio.NewScanner(bytes.NewReader(data))
	idx := 0
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var e Entry
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			results = append(results, VerificationResult{
				EntryIndex:    idx,
				Valid:         false,
				ChainValid:    false,
				FailureReason: fmt.Sprintf("JSON parse error: %v", err),
			})
			idx++
			continue
		}

		res := VerificationResult{
			EntryID:    e.ID,
			EntryIndex: idx,
			Valid:      true,
			ChainValid: true,
		}

		storedDigest := e.Digest

		if storedDigest == "" {
			res.Valid = false
			res.FailureReason = "missing digest field (legacy entry predating hash-chain)"
		} else {
			// Recompute digest with Digest zeroed.
			e.Digest = ""
			canonical, _ := json.Marshal(e)
			sum := sha256.Sum256(canonical)
			computed := hex.EncodeToString(sum[:])
			if computed != storedDigest {
				res.Valid = false
				res.FailureReason = fmt.Sprintf("digest mismatch: stored=%s computed=%s", storedDigest, computed)
			}
		}

		// Verify chain linkage.
		if idx == 0 {
			if e.PreviousDigest != "" {
				res.ChainValid = false
				res.FailureReason = joinReasons(res.FailureReason, "first entry has non-empty previous_digest")
			}
		} else if e.PreviousDigest != prevDigest {
			res.ChainValid = false
			res.FailureReason = joinReasons(res.FailureReason,
				fmt.Sprintf("chain break: expected previous_digest=%s got=%s", prevDigest, e.PreviousDigest))
		}

		prevDigest = storedDigest
		results = append(results, res)
		idx++
	}
	return results, nil
}

// joinReasons concatenates two failure reason strings with a semicolon separator.
func joinReasons(a, b string) string {
	if a == "" {
		return b
	}
	return a + "; " + b
}
