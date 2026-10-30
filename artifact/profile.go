package artifact

// ProductProfile is a Formulary-native product context artifact.
//
// It provides structured product metadata consumed by formula's system-card
// and context generators when a full Gemara ProductProfile is not yet available
// from onboarding. Fields map directly to ISO/IEC 42001 Annex A system-card
// sections and to compound's 4.3 scope / 8.1 operational context.
//
// ProductProfile may be embedded in a formula ProgramState JSON (as "profile")
// or loaded from a standalone YAML file by tools that need it separately.
//
// All fields are optional. Tools emit [DATA NEEDED: ...] placeholders for
// missing fields — the profile is additive, never blocking.
type ProductProfile struct {
	// Name is the product or system name (overrides program slug in display).
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// Version is the product release / version string.
	Version string `json:"version,omitempty" yaml:"version,omitempty"`
	// Description is a free-text system description for the AI system card
	// and management system context (ISO 42001 Annex A.2.2).
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// Classification categorizes the AI approach: e.g. "generative-ai",
	// "ml-classification", "rule-based", "rag", "llm-inference".
	Classification string `json:"classification,omitempty" yaml:"classification,omitempty"`
	// IntendedUse describes the intended use cases and target users.
	IntendedUse string `json:"intended_use,omitempty" yaml:"intended_use,omitempty"`
	// Limitations describes known limitations, failure modes, and
	// out-of-scope use cases.
	Limitations string `json:"limitations,omitempty" yaml:"limitations,omitempty"`
	// TrainingData describes training data sources, data governance approach,
	// and bias mitigation steps.
	TrainingData string `json:"training_data,omitempty" yaml:"training_data,omitempty"`
	// Evaluation describes evaluation methodology, metrics, and results.
	Evaluation string `json:"evaluation,omitempty" yaml:"evaluation,omitempty"`
	// Governance describes oversight bodies, review cadence, and
	// accountability structures.
	Governance string `json:"governance,omitempty" yaml:"governance,omitempty"`
	// Owner is the accountable product / program owner name.
	Owner string `json:"owner,omitempty" yaml:"owner,omitempty"`
	// Tags are free-form labels used for routing and reporting.
	Tags []string `json:"tags,omitempty" yaml:"tags,omitempty"`
}
