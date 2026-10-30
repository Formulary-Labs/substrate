package artifact_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Formulary-Labs/substrate/artifact"
)

func TestDetectFormat_oscalCatalog(t *testing.T) {
	format, err := artifact.DetectFormat(filepath.Join("testdata", "good-catalog.yaml"))
	if err != nil {
		t.Fatalf("DetectFormat: %v", err)
	}
	if format != artifact.FormatOSCAL {
		t.Fatalf("format = %q, want oscal", format)
	}
}

func TestDetectFormat_gemaraCatalog(t *testing.T) {
	path := filepath.Join("..", "..", "probe", "testdata", "good-control-catalog.yaml")
	if _, err := os.Stat(path); err != nil {
		t.Skipf("gemara fixture not found: %v", err)
	}
	format, err := artifact.DetectFormat(path)
	if err != nil {
		t.Fatalf("DetectFormat: %v", err)
	}
	if format != artifact.FormatGemara {
		t.Fatalf("format = %q, want gemara", format)
	}
}

func TestLoadOSCALCatalog_projectsControls(t *testing.T) {
	path := filepath.Join("testdata", "good-catalog.yaml")
	cat, err := artifact.LoadOSCALCatalog(path)
	if err != nil {
		t.Fatalf("LoadOSCALCatalog: %v", err)
	}
	if cat.Metadata.Title != "Example Access Control Catalog" {
		t.Fatalf("title = %q", cat.Metadata.Title)
	}

	unified, err := artifact.LoadCatalog(path)
	if err != nil {
		t.Fatalf("LoadCatalog: %v", err)
	}
	if unified.Format != artifact.FormatOSCAL {
		t.Fatalf("format = %q", unified.Format)
	}
	if unified.ID != "nist-800-53" {
		t.Fatalf("id = %q, want nist-800-53", unified.ID)
	}
	if len(unified.Controls) != 1 || unified.Controls[0].ID != "ac-1" {
		t.Fatalf("controls = %+v", unified.Controls)
	}
	if unified.Controls[0].Group != "ac" {
		t.Fatalf("group = %q", unified.Controls[0].Group)
	}
	if unified.Controls[0].Objective == "" {
		t.Fatal("expected statement prose as objective")
	}
	projected := unified.AsControlCatalog()
	if projected == nil || len(projected.Controls) != 1 {
		t.Fatalf("projected catalog = %+v", projected)
	}
	if projected.Metadata.Id != "nist-800-53" {
		t.Fatalf("projected id = %q", projected.Metadata.Id)
	}
	if projected.Controls[0].Id != "ac-1" {
		t.Fatalf("projected control = %+v", projected.Controls[0])
	}
}

func TestLoadCatalog_gemaraPreservesOriginal(t *testing.T) {
	path := filepath.Join("..", "..", "probe", "testdata", "good-control-catalog.yaml")
	if _, err := os.Stat(path); err != nil {
		t.Skipf("gemara fixture not found: %v", err)
	}
	unified, err := artifact.LoadCatalog(path)
	if err != nil {
		t.Fatalf("LoadCatalog: %v", err)
	}
	if unified.Format != artifact.FormatGemara || unified.Gemara == nil {
		t.Fatalf("expected gemara catalog, got format %q", unified.Format)
	}
	if unified.AsControlCatalog() != unified.Gemara {
		t.Fatal("AsControlCatalog should return the original gemara catalog")
	}
	if len(unified.Controls) == 0 {
		t.Fatal("expected controls copied onto the unified catalog")
	}
}

func TestLoadOSCAL_jsonCatalog(t *testing.T) {
	doc, err := artifact.LoadOSCAL(filepath.Join("testdata", "good-catalog.json"))
	if err != nil {
		t.Fatalf("LoadOSCAL: %v", err)
	}
	if doc.Catalog == nil {
		t.Fatal("expected catalog model")
	}
}

func TestLoadOSCALAssessmentResults(t *testing.T) {
	path := filepath.Join("testdata", "good-assessment-results.yaml")
	results, err := artifact.LoadOSCALAssessmentResults(path)
	if err != nil {
		t.Fatalf("LoadOSCALAssessmentResults: %v", err)
	}
	if results.Metadata.Title == "" {
		t.Fatal("expected a title")
	}
	if _, err := artifact.LoadOSCALCatalog(path); err == nil {
		t.Fatal("expected catalog loader to reject assessment-results")
	}
}

func TestValidateOSCAL_invalidCatalog(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "bad-catalog.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	modelType, issues, err := artifact.ValidateOSCAL(data)
	if err != nil {
		t.Fatalf("ValidateOSCAL error: %v", err)
	}
	if modelType != artifact.OSCALModelCatalog {
		t.Fatalf("model = %q", modelType)
	}
	if len(issues) == 0 {
		t.Fatal("expected schema issues for a catalog missing metadata.title")
	}
}

func TestLoadOSCAL_rejectsInvalid(t *testing.T) {
	_, err := artifact.LoadOSCAL(filepath.Join("testdata", "bad-catalog.yaml"))
	if err == nil {
		t.Fatal("expected invalid catalog to fail LoadOSCAL")
	}
}
