package artifact

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/defenseunicorns/go-oscal/src/pkg/validation"
	oscal "github.com/defenseunicorns/go-oscal/src/types/oscal-1-2-2"
	"gopkg.in/yaml.v3"
)

// OSCAL model names. They match the single root key of a native document.
const (
	OSCALModelCatalog             = "catalog"
	OSCALModelProfile             = "profile"
	OSCALModelComponentDefinition = "component-definition"
	OSCALModelSSP                 = "system-security-plan"
	OSCALModelAssessmentPlan      = "assessment-plan"
	OSCALModelAssessmentResults   = "assessment-results"
	OSCALModelPOAM                = "plan-of-action-and-milestones"
	OSCALModelMappingCollection   = "mapping-collection"
)

// Re-export OSCAL 1.2.2 model types so tools do not import go-oscal directly.
type (
	// OSCALModels is a complete OSCAL document. Exactly one model pointer is set.
	OSCALModels = oscal.OscalModels
	// OSCALCatalog is an OSCAL catalog.
	OSCALCatalog = oscal.Catalog
	// OSCALProfile is an OSCAL profile.
	OSCALProfile = oscal.Profile
	// OSCALAssessmentResults is an OSCAL assessment-results document.
	OSCALAssessmentResults = oscal.AssessmentResults
	// OSCALPOAM is an OSCAL plan-of-action-and-milestones document.
	OSCALPOAM = oscal.PlanOfActionAndMilestones
	// OSCALSSP is an OSCAL system-security-plan.
	OSCALSSP = oscal.SystemSecurityPlan
	// OSCALComponentDefinition is an OSCAL component-definition.
	OSCALComponentDefinition = oscal.ComponentDefinition
	// OSCALMappingCollection is an OSCAL mapping-collection.
	OSCALMappingCollection = oscal.MappingCollection
	// OSCALAssessmentPlan is an OSCAL assessment-plan.
	OSCALAssessmentPlan = oscal.AssessmentPlan
	// OSCALControl is one control inside a catalog group or catalog root.
	OSCALControl = oscal.Control
	// OSCALGroup is a catalog group (control family).
	OSCALGroup = oscal.Group
)

// OSCALIssue is one schema finding from an OSCAL document.
type OSCALIssue struct {
	// Field is the instance path (JSON pointer without the leading #), when known.
	Field string
	// Message is the schema error text.
	Message string
}

// ValidateOSCAL checks data against the OSCAL schema for the oscal-version
// declared in the document. A nil issue slice and a nil error means the
// document is schema-valid. Schema findings are returned as issues with a
// nil error. A non-nil error means the document could not be interpreted
// as OSCAL at all (missing model root, missing oscal-version, unreadable YAML).
func ValidateOSCAL(data []byte) (string, []OSCALIssue, error) {
	validator, err := validation.NewValidator(data)
	if err != nil {
		return "", nil, err
	}
	modelType := validator.GetModelType()
	if err := validator.Validate(); err != nil {
		result, resErr := validator.GetValidationResult()
		if resErr != nil || len(result.Errors) == 0 {
			return modelType, []OSCALIssue{{Message: err.Error()}}, nil
		}
		issues := make([]OSCALIssue, 0, len(result.Errors))
		for _, item := range result.Errors {
			issues = append(issues, OSCALIssue{
				Field:   item.InstanceLocation,
				Message: item.Error,
			})
		}
		return modelType, issues, nil
	}
	return modelType, nil, nil
}

// LoadOSCAL reads and schema-validates an OSCAL document.
func LoadOSCAL(path string) (*OSCALModels, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading OSCAL document %q: %w", path, err)
	}
	modelType, issues, err := ValidateOSCAL(data)
	if err != nil {
		return nil, fmt.Errorf("validating OSCAL document %q: %w", path, err)
	}
	if len(issues) > 0 {
		return nil, fmt.Errorf("validating OSCAL document %q: %s", path, formatOSCALIssues(issues))
	}
	var doc OSCALModels
	if err := unmarshalOSCAL(data, &doc); err != nil {
		return nil, fmt.Errorf("decoding OSCAL document %q: %w", path, err)
	}
	if oscalModelName(&doc) == "" {
		return nil, fmt.Errorf("OSCAL document %q decoded without a model (detected %s)", path, modelType)
	}
	return &doc, nil
}

// LoadOSCALCatalog loads a schema-valid OSCAL catalog.
func LoadOSCALCatalog(path string) (*OSCALCatalog, error) {
	doc, err := LoadOSCAL(path)
	if err != nil {
		return nil, err
	}
	if doc.Catalog == nil {
		return nil, fmt.Errorf("OSCAL document %q is %s, not a catalog", path, oscalModelName(doc))
	}
	return doc.Catalog, nil
}

// LoadOSCALProfile loads a schema-valid OSCAL profile.
func LoadOSCALProfile(path string) (*OSCALProfile, error) {
	doc, err := LoadOSCAL(path)
	if err != nil {
		return nil, err
	}
	if doc.Profile == nil {
		return nil, fmt.Errorf("OSCAL document %q is %s, not a profile", path, oscalModelName(doc))
	}
	return doc.Profile, nil
}

// LoadOSCALAssessmentResults loads a schema-valid OSCAL assessment-results document.
func LoadOSCALAssessmentResults(path string) (*OSCALAssessmentResults, error) {
	doc, err := LoadOSCAL(path)
	if err != nil {
		return nil, err
	}
	if doc.AssessmentResults == nil {
		return nil, fmt.Errorf("OSCAL document %q is %s, not assessment-results", path, oscalModelName(doc))
	}
	return doc.AssessmentResults, nil
}

// LoadOSCALPOAM loads a schema-valid OSCAL plan-of-action-and-milestones document.
func LoadOSCALPOAM(path string) (*OSCALPOAM, error) {
	doc, err := LoadOSCAL(path)
	if err != nil {
		return nil, err
	}
	if doc.PlanOfActionAndMilestones == nil {
		return nil, fmt.Errorf("OSCAL document %q is %s, not a plan-of-action-and-milestones", path, oscalModelName(doc))
	}
	return doc.PlanOfActionAndMilestones, nil
}

func unmarshalOSCAL(data []byte, dest any) error {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) > 0 && trimmed[0] == '{' {
		return json.Unmarshal(trimmed, dest)
	}
	return yaml.Unmarshal(data, dest)
}

func formatOSCALIssues(issues []OSCALIssue) string {
	parts := make([]string, 0, len(issues))
	for _, issue := range issues {
		if issue.Field != "" {
			parts = append(parts, issue.Field+": "+issue.Message)
			continue
		}
		parts = append(parts, issue.Message)
	}
	return strings.Join(parts, "; ")
}

func oscalModelName(doc *OSCALModels) string {
	if doc == nil {
		return ""
	}
	switch {
	case doc.Catalog != nil:
		return OSCALModelCatalog
	case doc.Profile != nil:
		return OSCALModelProfile
	case doc.ComponentDefinition != nil:
		return OSCALModelComponentDefinition
	case doc.SystemSecurityPlan != nil:
		return OSCALModelSSP
	case doc.AssessmentPlan != nil:
		return OSCALModelAssessmentPlan
	case doc.AssessmentResults != nil:
		return OSCALModelAssessmentResults
	case doc.PlanOfActionAndMilestones != nil:
		return OSCALModelPOAM
	case doc.MappingCollection != nil:
		return OSCALModelMappingCollection
	default:
		return ""
	}
}
