// Package artifact provides helpers for loading and working with gemara
// artifacts. It wraps github.com/gemaraproj/go-gemara to provide a
// simplified, Formulary-specific interface.
//
// All gemara type aliases live here so that Formulary tools import only this
// package — not go-gemara directly. This insulates tools from upstream API
// changes.
package artifact

import (
	"context"
	"fmt"
	"os"

	gemara "github.com/gemaraproj/go-gemara"
	"github.com/gemaraproj/go-gemara/fetcher"
)

// Re-export key gemara types so tools can use artifact.ControlCatalog instead
// of importing go-gemara directly.
type (
	// ControlCatalog re-exports gemara.ControlCatalog for Formulary tools.
	ControlCatalog = gemara.ControlCatalog
	// GuidanceCatalog re-exports gemara.GuidanceCatalog for Formulary tools.
	GuidanceCatalog = gemara.GuidanceCatalog
	// AuditLog re-exports gemara.AuditLog for Formulary tools.
	AuditLog = gemara.AuditLog
	// EvaluationLog re-exports gemara.EvaluationLog for Formulary tools.
	EvaluationLog = gemara.EvaluationLog
	// Metadata re-exports gemara.Metadata for Formulary tools.
	Metadata = gemara.Metadata
	// ArtifactType re-exports gemara.ArtifactType for Formulary tools.
	ArtifactType = gemara.ArtifactType //nolint:revive // stutter is intentional
	// MappingDocument re-exports gemara.MappingDocument for Formulary tools.
	MappingDocument = gemara.MappingDocument
	// Mapping re-exports gemara.Mapping for Formulary tools.
	Mapping = gemara.Mapping
	// MappingTarget re-exports gemara.MappingTarget for Formulary tools.
	MappingTarget = gemara.MappingTarget
	// RiskCatalog re-exports gemara.RiskCatalog for Formulary tools.
	RiskCatalog = gemara.RiskCatalog
	// Risk re-exports gemara.Risk for Formulary tools.
	Risk = gemara.Risk
	// Severity re-exports gemara.Severity for Formulary tools.
	Severity = gemara.Severity
	// Policy re-exports gemara.Policy for Formulary tools.
	Policy = gemara.Policy
	// Scope re-exports gemara.Scope for Formulary tools.
	Scope = gemara.Scope
	// RACI re-exports gemara.RACI for Formulary tools.
	RACI = gemara.RACI
	// Contact re-exports gemara.Contact for Formulary tools.
	Contact = gemara.Contact
	// Control re-exports gemara.Control for Formulary tools.
	Control = gemara.Control
	// Group re-exports gemara.Group for Formulary tools.
	Group = gemara.Group
	// AssessmentRequirement re-exports gemara.AssessmentRequirement for Formulary tools.
	AssessmentRequirement = gemara.AssessmentRequirement
)

// InvalidArtifact re-exports the gemara invalid artifact sentinel.
var InvalidArtifact = gemara.InvalidArtifact //nolint:gochecknoglobals

// Re-export ArtifactType constants so tools do not import go-gemara directly.
const (
	AuditLogArtifact          = gemara.AuditLogArtifact
	CapabilityCatalogArtifact = gemara.CapabilityCatalogArtifact
	ControlCatalogArtifact    = gemara.ControlCatalogArtifact
	EnforcementLogArtifact    = gemara.EnforcementLogArtifact
	EvaluationLogArtifact     = gemara.EvaluationLogArtifact
	GuidanceCatalogArtifact   = gemara.GuidanceCatalogArtifact
	LexiconArtifact           = gemara.LexiconArtifact
	MappingDocumentArtifact   = gemara.MappingDocumentArtifact
	PolicyArtifact            = gemara.PolicyArtifact
	PrincipleCatalogArtifact  = gemara.PrincipleCatalogArtifact
	RiskCatalogArtifact       = gemara.RiskCatalogArtifact
	ThreatCatalogArtifact     = gemara.ThreatCatalogArtifact
	VectorCatalogArtifact     = gemara.VectorCatalogArtifact
)

// LoadControlCatalog loads a gemara ControlCatalog from the given file path.
// Supports .yaml, .yml, and .json extensions.
func LoadControlCatalog(path string) (*ControlCatalog, error) {
	f := &fetcher.File{}
	catalog, err := gemara.Load[ControlCatalog](context.Background(), f, path)
	if err != nil {
		return nil, fmt.Errorf("loading control catalog from %q: %w", path, err)
	}
	return catalog, nil
}

// LoadGuidanceCatalog loads a gemara GuidanceCatalog from the given file path.
func LoadGuidanceCatalog(path string) (*GuidanceCatalog, error) {
	f := &fetcher.File{}
	catalog, err := gemara.Load[GuidanceCatalog](context.Background(), f, path)
	if err != nil {
		return nil, fmt.Errorf("loading guidance catalog from %q: %w", path, err)
	}
	return catalog, nil
}

// LoadAuditLog loads a gemara AuditLog from the given file path.
func LoadAuditLog(path string) (*AuditLog, error) {
	f := &fetcher.File{}
	log, err := gemara.Load[AuditLog](context.Background(), f, path)
	if err != nil {
		return nil, fmt.Errorf("loading audit log from %q: %w", path, err)
	}
	return log, nil
}

// LoadEvaluationLog loads a gemara EvaluationLog from the given file path.
func LoadEvaluationLog(path string) (*EvaluationLog, error) {
	f := &fetcher.File{}
	log, err := gemara.Load[EvaluationLog](context.Background(), f, path)
	if err != nil {
		return nil, fmt.Errorf("loading evaluation log from %q: %w", path, err)
	}
	return log, nil
}

// LoadMappingDocument loads a gemara MappingDocument from the given file path.
func LoadMappingDocument(path string) (*MappingDocument, error) {
	f := &fetcher.File{}
	doc, err := gemara.Load[MappingDocument](context.Background(), f, path)
	if err != nil {
		return nil, fmt.Errorf("loading mapping document from %q: %w", path, err)
	}
	return doc, nil
}

// LoadRiskCatalog loads a gemara RiskCatalog from the given file path.
// Supports .yaml, .yml, and .json extensions.
func LoadRiskCatalog(path string) (*RiskCatalog, error) {
	f := &fetcher.File{}
	catalog, err := gemara.Load[RiskCatalog](context.Background(), f, path)
	if err != nil {
		return nil, fmt.Errorf("loading risk catalog from %q: %w", path, err)
	}
	return catalog, nil
}

// LoadPolicy loads a gemara Policy from the given file path.
// Supports .yaml, .yml, and .json extensions.
func LoadPolicy(path string) (*Policy, error) {
	f := &fetcher.File{}
	policy, err := gemara.Load[Policy](context.Background(), f, path)
	if err != nil {
		return nil, fmt.Errorf("loading policy from %q: %w", path, err)
	}
	return policy, nil
}

// DetectType reads the metadata.type field from an artifact file to identify
// its type without a full unmarshal. Useful for routing in tools that accept
// multiple artifact types.
func DetectType(path string) (ArtifactType, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return InvalidArtifact, fmt.Errorf("reading artifact at %q: %w", path, err)
	}
	t, err := gemara.DetectType(data)
	if err != nil {
		return InvalidArtifact, fmt.Errorf("detecting type in %q: %w", path, err)
	}
	return t, nil
}
