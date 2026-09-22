package artifact_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Formulary-Labs/substrate/artifact"
)

// These tests use the real gemara test fixtures from the probe tool's testdata
// directory, which is the nearest available set of valid gemara artifacts.
const (
	testdataDir    = "../../probe/testdata"
	catalogFixture = "good-control-catalog.yaml"
	evalLogFixture = "good-evaluation-log.yaml"
	lexiconFixture = "good-lexicon.yaml"
	badFixture     = "bad.yaml"
)

func TestLoadControlCatalog_valid(t *testing.T) {
	path := filepath.Join(testdataDir, catalogFixture)
	if _, err := os.Stat(path); err != nil {
		t.Skipf("test fixture not found: %v", err)
	}

	cat, err := artifact.LoadControlCatalog(path)
	if err != nil {
		t.Fatalf("LoadControlCatalog error: %v", err)
	}
	if cat == nil {
		t.Fatal("expected non-nil catalog")
	}
	if cat.Metadata.Id == "" {
		t.Error("catalog metadata.id must not be empty")
	}
	if len(cat.Controls) == 0 {
		t.Error("expected at least one control in catalog")
	}
}

func TestLoadEvaluationLog_valid(t *testing.T) {
	path := filepath.Join(testdataDir, evalLogFixture)
	if _, err := os.Stat(path); err != nil {
		t.Skipf("test fixture not found: %v", err)
	}

	log, err := artifact.LoadEvaluationLog(path)
	if err != nil {
		t.Fatalf("LoadEvaluationLog error: %v", err)
	}
	if log == nil {
		t.Fatal("expected non-nil evaluation log")
	}
	if log.Metadata.Id == "" {
		t.Error("evaluation log metadata.id must not be empty")
	}
}

func TestDetectType_controlCatalog(t *testing.T) {
	path := filepath.Join(testdataDir, catalogFixture)
	if _, err := os.Stat(path); err != nil {
		t.Skipf("test fixture not found: %v", err)
	}

	typ, err := artifact.DetectType(path)
	if err != nil {
		t.Fatalf("DetectType error: %v", err)
	}
	if typ == artifact.InvalidArtifact {
		t.Error("expected valid artifact type, got InvalidArtifact")
	}
}

func TestDetectType_evaluationLog(t *testing.T) {
	path := filepath.Join(testdataDir, evalLogFixture)
	if _, err := os.Stat(path); err != nil {
		t.Skipf("test fixture not found: %v", err)
	}

	typ, err := artifact.DetectType(path)
	if err != nil {
		t.Fatalf("DetectType error: %v", err)
	}
	if typ == artifact.InvalidArtifact {
		t.Error("expected valid artifact type for evaluation log")
	}
}

func TestDetectType_missingFile(t *testing.T) {
	_, err := artifact.DetectType("/no/such/file.yaml")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadControlCatalog_missingFile(t *testing.T) {
	_, err := artifact.LoadControlCatalog("/no/such/catalog.yaml")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadEvaluationLog_missingFile(t *testing.T) {
	_, err := artifact.LoadEvaluationLog("/no/such/eval.yaml")
	if err == nil {
		t.Error("expected error for missing file")
	}
}
