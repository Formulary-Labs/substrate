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
	"sort"

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
	// CapabilityCatalog re-exports gemara.CapabilityCatalog for Formulary tools.
	CapabilityCatalog = gemara.CapabilityCatalog
	// ThreatCatalog re-exports gemara.ThreatCatalog for Formulary tools.
	ThreatCatalog = gemara.ThreatCatalog
	// VectorCatalog re-exports gemara.VectorCatalog for Formulary tools.
	VectorCatalog = gemara.VectorCatalog
	// PrincipleCatalog re-exports gemara.PrincipleCatalog for Formulary tools.
	PrincipleCatalog = gemara.PrincipleCatalog
	// EnforcementLog re-exports gemara.EnforcementLog for Formulary tools.
	EnforcementLog = gemara.EnforcementLog
	// Lexicon re-exports gemara.Lexicon for Formulary tools.
	Lexicon = gemara.Lexicon
)

// Re-export published enums so tools use the schema vocabularies without
// importing go-gemara directly. Emitted values must be these values: a tool
// that invents a parallel vocabulary breaks interoperability with every other
// conformant tool.
type (
	// Result re-exports gemara.Result (Not Run, Passed, Failed, Needs Review,
	// Not Applicable, Unknown).
	Result = gemara.Result
	// ConfidenceLevel re-exports gemara.ConfidenceLevel (Undetermined, Low,
	// Medium, High).
	ConfidenceLevel = gemara.ConfidenceLevel
	// RelationshipType re-exports gemara.RelationshipType (implements,
	// implemented-by, supports, supported-by, equivalent, subsumes, no-match,
	// relates-to).
	RelationshipType = gemara.RelationshipType
	// RiskAppetite re-exports gemara.RiskAppetite (Minimal, Low, Moderate, High).
	RiskAppetite = gemara.RiskAppetite
	// MethodType re-exports gemara.MethodType (Behavioral, Intent, Remediation,
	// Gate).
	MethodType = gemara.MethodType
	// ModeType re-exports gemara.ModeType (Manual, Automated).
	ModeType = gemara.ModeType
	// EntityType re-exports gemara.EntityType (Human, Software, Software
	// Assisted). This is the declared-executor vocabulary: substituting an
	// agent for a tool is expressed by changing the executor, not the method.
	EntityType = gemara.EntityType
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

// LoadCapabilityCatalog loads a gemara CapabilityCatalog from the given file path.
func LoadCapabilityCatalog(path string) (*CapabilityCatalog, error) {
	f := &fetcher.File{}
	catalog, err := gemara.Load[CapabilityCatalog](context.Background(), f, path)
	if err != nil {
		return nil, fmt.Errorf("loading capability catalog from %q: %w", path, err)
	}
	return catalog, nil
}

// LoadThreatCatalog loads a gemara ThreatCatalog from the given file path.
func LoadThreatCatalog(path string) (*ThreatCatalog, error) {
	f := &fetcher.File{}
	catalog, err := gemara.Load[ThreatCatalog](context.Background(), f, path)
	if err != nil {
		return nil, fmt.Errorf("loading threat catalog from %q: %w", path, err)
	}
	return catalog, nil
}

// LoadVectorCatalog loads a gemara VectorCatalog from the given file path.
func LoadVectorCatalog(path string) (*VectorCatalog, error) {
	f := &fetcher.File{}
	catalog, err := gemara.Load[VectorCatalog](context.Background(), f, path)
	if err != nil {
		return nil, fmt.Errorf("loading vector catalog from %q: %w", path, err)
	}
	return catalog, nil
}

// LoadPrincipleCatalog loads a gemara PrincipleCatalog from the given file path.
func LoadPrincipleCatalog(path string) (*PrincipleCatalog, error) {
	f := &fetcher.File{}
	catalog, err := gemara.Load[PrincipleCatalog](context.Background(), f, path)
	if err != nil {
		return nil, fmt.Errorf("loading principle catalog from %q: %w", path, err)
	}
	return catalog, nil
}

// LoadEnforcementLog loads a gemara EnforcementLog from the given file path.
func LoadEnforcementLog(path string) (*EnforcementLog, error) {
	f := &fetcher.File{}
	log, err := gemara.Load[EnforcementLog](context.Background(), f, path)
	if err != nil {
		return nil, fmt.Errorf("loading enforcement log from %q: %w", path, err)
	}
	return log, nil
}

// LoadLexicon loads a gemara Lexicon from the given file path.
func LoadLexicon(path string) (*Lexicon, error) {
	f := &fetcher.File{}
	lex, err := gemara.Load[Lexicon](context.Background(), f, path)
	if err != nil {
		return nil, fmt.Errorf("loading lexicon from %q: %w", path, err)
	}
	return lex, nil
}

// loaders maps every supported artifact type to its typed loader. Keeping it as
// a single table is what stops the boundary drifting: a new published artifact
// type either appears here or SupportedTypes does not claim it.
var loaders = map[ArtifactType]func(path string) (any, error){ //nolint:gochecknoglobals // fixed dispatch table
	AuditLogArtifact:          func(p string) (any, error) { return unwrap(LoadAuditLog(p)) },
	CapabilityCatalogArtifact: func(p string) (any, error) { return unwrap(LoadCapabilityCatalog(p)) },
	ControlCatalogArtifact:    func(p string) (any, error) { return unwrap(LoadControlCatalog(p)) },
	EnforcementLogArtifact:    func(p string) (any, error) { return unwrap(LoadEnforcementLog(p)) },
	EvaluationLogArtifact:     func(p string) (any, error) { return unwrap(LoadEvaluationLog(p)) },
	GuidanceCatalogArtifact:   func(p string) (any, error) { return unwrap(LoadGuidanceCatalog(p)) },
	LexiconArtifact:           func(p string) (any, error) { return unwrap(LoadLexicon(p)) },
	MappingDocumentArtifact:   func(p string) (any, error) { return unwrap(LoadMappingDocument(p)) },
	PolicyArtifact:            func(p string) (any, error) { return unwrap(LoadPolicy(p)) },
	PrincipleCatalogArtifact:  func(p string) (any, error) { return unwrap(LoadPrincipleCatalog(p)) },
	RiskCatalogArtifact:       func(p string) (any, error) { return unwrap(LoadRiskCatalog(p)) },
	ThreatCatalogArtifact:     func(p string) (any, error) { return unwrap(LoadThreatCatalog(p)) },
	VectorCatalogArtifact:     func(p string) (any, error) { return unwrap(LoadVectorCatalog(p)) },
}

// unwrap normalizes a typed load result into an any. Without this, a failed
// load returns an interface holding a nil pointer, which is not == nil: callers
// that check for nil would think a failed load succeeded.
func unwrap[T any](v *T, err error) (any, error) {
	if err != nil {
		return nil, err
	}
	return v, nil
}

// SupportedTypes returns the artifact types this package can load, sorted by
// their published name. Callers use it to report what a tool accepts instead of
// hard-coding a list that drifts.
func SupportedTypes() []ArtifactType {
	types := make([]ArtifactType, 0, len(loaders))
	for t := range loaders {
		types = append(types, t)
	}
	sort.Slice(types, func(i, j int) bool { return types[i].String() < types[j].String() })
	return types
}

// LoadByType loads an artifact of a known type and returns it as an any, so a
// caller that routes on DetectType does not need its own type switch. The
// concrete type is the one named by the artifact type.
func LoadByType(path string, t ArtifactType) (any, error) {
	load, ok := loaders[t]
	if !ok {
		return nil, fmt.Errorf("unsupported artifact type %q (supported: %v)", t.String(), typeNames())
	}
	return load(path)
}

// typeNames returns the supported type names for error messages.
func typeNames() []string {
	types := SupportedTypes()
	names := make([]string, len(types))
	for i, t := range types {
		names[i] = t.String()
	}
	return names
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
