package provenance_test

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/Formulary-Labs/substrate/provenance"
)

func TestWrite_creates_file(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "logs", "provenance.jsonl")

	e := provenance.Entry{
		Spec:        "functions/test-spec.md",
		Output:      "data/test/out.json",
		OutputType:  "other",
		Program:     "test-program",
		Purpose:     "Test provenance write",
		Reusability: provenance.Instance,
		QualityGate: provenance.Pass,
		Tool:        "test-tool",
	}

	if err := provenance.Write(logPath, e); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}

	// Verify file exists and contains one valid JSON line.
	f, err := os.Open(logPath)
	if err != nil {
		t.Fatalf("could not open log file: %v", err)
	}
	defer f.Close() //nolint:errcheck // test read-only file, close error is harmless

	scanner := bufio.NewScanner(f)
	var lines []provenance.Entry
	for scanner.Scan() {
		var entry provenance.Entry
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			t.Fatalf("invalid JSON in log: %v", err)
		}
		lines = append(lines, entry)
	}

	if len(lines) != 1 {
		t.Fatalf("expected 1 line in log, got %d", len(lines))
	}
	if lines[0].Program != "test-program" {
		t.Errorf("program = %q, want %q", lines[0].Program, "test-program")
	}
	if lines[0].ID == "" {
		t.Error("ID must not be empty")
	}
	if lines[0].Timestamp == "" {
		t.Error("Timestamp must not be empty")
	}
}

func TestWrite_appends(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "provenance.jsonl")

	e := provenance.Entry{
		Tool:        "test-tool",
		Program:     "prog",
		Purpose:     "first entry",
		Reusability: provenance.Instance,
		QualityGate: provenance.Pass,
	}
	for range 3 {
		if err := provenance.Write(logPath, e); err != nil {
			t.Fatalf("Write error: %v", err)
		}
	}

	f, err := os.Open(logPath)
	if err != nil {
		t.Fatalf("could not reopen log file: %v", err)
	}
	defer f.Close() //nolint:errcheck // test read-only file, close error is harmless
	var count int
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		count++
	}
	if count != 3 {
		t.Errorf("expected 3 log lines, got %d", count)
	}
}
