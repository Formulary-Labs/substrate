package artifact_test

import (
	"os"
	"path/filepath"
	"strings"
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

// TestSupportedTypes_coversEveryPublishedType is the drift guard for this
// boundary. The package doc promises tools import only substrate, never
// go-gemara. That promise breaks silently the moment a published artifact type
// has an exported constant but no loader. If gemara adds a fourteenth artifact
// type, this test keeps passing until someone adds the constant — so the guard
// is: every constant that exists must be loadable.
func TestSupportedTypes_coversEveryPublishedType(t *testing.T) {
	published := map[artifact.ArtifactType]string{
		artifact.AuditLogArtifact:          "AuditLog",
		artifact.CapabilityCatalogArtifact: "CapabilityCatalog",
		artifact.ControlCatalogArtifact:    "ControlCatalog",
		artifact.EnforcementLogArtifact:    "EnforcementLog",
		artifact.EvaluationLogArtifact:     "EvaluationLog",
		artifact.GuidanceCatalogArtifact:   "GuidanceCatalog",
		artifact.LexiconArtifact:           "Lexicon",
		artifact.MappingDocumentArtifact:   "MappingDocument",
		artifact.PolicyArtifact:            "Policy",
		artifact.PrincipleCatalogArtifact:  "PrincipleCatalog",
		artifact.RiskCatalogArtifact:       "RiskCatalog",
		artifact.ThreatCatalogArtifact:     "ThreatCatalog",
		artifact.VectorCatalogArtifact:     "VectorCatalog",
	}

	supported := make(map[artifact.ArtifactType]bool, len(published))
	for _, got := range artifact.SupportedTypes() {
		supported[got] = true
		if want, ok := published[got]; ok && got.String() != want {
			t.Errorf("artifact type String() = %q, want %q", got.String(), want)
		}
	}

	for got, name := range published {
		if !supported[got] {
			t.Errorf("artifact type %q has a constant but no loader: SupportedTypes omits it", name)
		}
	}
}

func TestSupportedTypes_count(t *testing.T) {
	// The schema defines thirteen artifact types across seven layers.
	if got := len(artifact.SupportedTypes()); got != 13 {
		t.Fatalf("SupportedTypes() = %d types, want 13", got)
	}
}

// TestLoadByType_missingFile routes every supported type through the dispatcher
// with a path that does not exist. It proves each entry in the dispatch table is
// wired to a loader that fails cleanly rather than panicking, without needing a
// fixture per type.
func TestLoadByType_missingFile(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "absent.yaml")
	for _, at := range artifact.SupportedTypes() {
		t.Run(at.String(), func(t *testing.T) {
			got, err := artifact.LoadByType(missing, at)
			if err == nil {
				t.Fatalf("LoadByType(%q) = %v, want error for missing file", at.String(), got)
			}
			if got != nil {
				t.Errorf("LoadByType(%q) returned a value alongside an error", at.String())
			}
		})
	}
}

func TestLoadByType_unsupported(t *testing.T) {
	_, err := artifact.LoadByType("nowhere.yaml", artifact.InvalidArtifact)
	if err == nil {
		t.Fatal("expected an error for an unsupported artifact type")
	}
	// The error must name what is supported, so a caller can fix the call.
	for _, want := range []string{"ControlCatalog", "VectorCatalog"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q should name supported type %q", err, want)
		}
	}
}

// TestLoadDistribution asserts the fixtures borrowed from probe still load
// through the convenience loaders, and skips when the sibling repo is absent.
func TestLoadLexicon_valid(t *testing.T) {
	path := filepath.Join(testdataDir, lexiconFixture)
	if _, err := os.Stat(path); err != nil {
		t.Skipf("test fixture not found: %v", err)
	}
	if _, err := artifact.LoadLexicon(path); err != nil {
		t.Fatalf("LoadLexicon error: %v", err)
	}
}
