// Package losstable provides a shared organizational loss table for
// consistent FAIR (Factor Analysis of Information Risk) ALE calculations
// across subsidiaries and business units.
//
// A loss table is a YAML configuration file maintained at the org level that
// converts human-readable system descriptors — data classifications and system
// tiers — into USD figures. Using a shared table ensures that ALE values from
// different teams and subsidiaries are calculated on the same baseline
// assumptions, making cross-org risk prioritization meaningful.
//
// Example loss table (org-loss-table.yaml):
//
//	schema_version: "1.0"
//	cost_per_record:
//	  pii:       150   # USD per PII record (GDPR/CCPA baseline)
//	  phi:       250   # USD per PHI record (HIPAA exposure)
//	  financial: 200
//	  public:    0
//	downtime_cost_per_hour:
//	  tier_1: 50000   # customer-facing systems
//	  tier_2: 10000   # internal tools
//	  tier_3: 1000    # dev/test
//	assumed_exposure_hours: 72   # default RTO window for ALE calculation
//	regulatory_fines:
//	  gdpr:  500000
//	  hipaa: 250000
package losstable

import (
	"fmt"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
)

const currentSchemaVersion = "1.0"

// LossTable holds the org-level financial parameters used to derive
// AssetValue for FAIR ALE calculations. Load from a YAML file with Load().
type LossTable struct {
	// SchemaVersion identifies the loss table format version.
	SchemaVersion string `yaml:"schema_version"`

	// CostPerRecord maps data classification labels to USD cost per record
	// lost or exposed. Standard labels: pii, phi, financial, public.
	// Custom labels are allowed and referenced by --data-class.
	CostPerRecord map[string]float64 `yaml:"cost_per_record"`

	// DowntimeCostPerHour maps system tier labels to USD cost per hour of
	// downtime. Standard labels: tier_1, tier_2, tier_3.
	DowntimeCostPerHour map[string]float64 `yaml:"downtime_cost_per_hour"`

	// AssumedExposureHours is the default recovery time objective (RTO) in
	// hours used when computing AssetValue for tier-based risks.
	// Overridable per-risk with --exposure-hours.
	AssumedExposureHours float64 `yaml:"assumed_exposure_hours"`

	// RegulatoryFines maps jurisdiction labels to USD per incident. These
	// are informational and not currently included in ALE computation.
	RegulatoryFines map[string]float64 `yaml:"regulatory_fines,omitempty"`
}

// Load reads a YAML loss table from the given path and validates it.
func Load(path string) (*LossTable, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("losstable: reading %q: %w", path, err)
	}
	var t LossTable
	if err := yaml.Unmarshal(data, &t); err != nil {
		return nil, fmt.Errorf("losstable: parsing %q: %w", path, err)
	}
	if t.SchemaVersion != "" && t.SchemaVersion != currentSchemaVersion {
		return nil, fmt.Errorf("losstable: unsupported schema version %q in %q (expected %s)",
			t.SchemaVersion, path, currentSchemaVersion)
	}
	if t.CostPerRecord == nil && t.DowntimeCostPerHour == nil {
		return nil, fmt.Errorf("losstable: %q must define at least cost_per_record or downtime_cost_per_hour", path)
	}
	return &t, nil
}

// AssetValueForRecords returns the USD asset value for a risk that exposes
// recordCount records of the given data classification. Returns 0 if the
// data class is not defined in the table.
//
// Example: AssetValueForRecords("pii", 500000) → 500000 × $150 = $75,000,000
func (t *LossTable) AssetValueForRecords(dataClass string, recordCount int) float64 {
	if t == nil || t.CostPerRecord == nil {
		return 0
	}
	cost, ok := t.CostPerRecord[strings.ToLower(dataClass)]
	if !ok {
		return 0
	}
	return cost * float64(recordCount)
}

// AssetValueForTier returns the USD asset value for a risk affecting a system
// of the given tier, calculated as downtime cost × exposureHours. If
// exposureHours is 0, AssumedExposureHours from the table is used. Returns 0
// if the tier is not defined in the table.
//
// Example: AssetValueForTier("tier_1", 0) → $50,000/hr × 72 hr = $3,600,000
func (t *LossTable) AssetValueForTier(tier string, exposureHours float64) float64 {
	if t == nil || t.DowntimeCostPerHour == nil {
		return 0
	}
	cost, ok := t.DowntimeCostPerHour[strings.ToLower(tier)]
	if !ok {
		return 0
	}
	hours := exposureHours
	if hours == 0 {
		hours = t.AssumedExposureHours
	}
	if hours == 0 {
		hours = 72 // hard fallback if not configured
	}
	return cost * hours
}
