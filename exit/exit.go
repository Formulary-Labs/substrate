// Package exit defines the exit code contract for all Formulary tools.
//
// Exit codes are meaningful:
//   - 0  OK          — clean run, output is valid
//   - 1  Validation  — input or output failed schema/logic validation
//   - 2  ToolError   — unexpected failure (missing file, network error, etc.)
//
// Tools must exit with one of these three codes. No other exit codes are
// defined. Callers (CI pipelines, agent orchestrators) may rely on this
// three-value contract.
package exit

import "os"

const (
	// OK indicates a clean run. All checks passed and output is valid.
	OK = 0

	// Validation indicates that input or output failed schema or logic
	// validation. The tool ran correctly; the data did not meet requirements.
	Validation = 1

	// ToolError indicates an unexpected failure: missing file, parse error,
	// network failure, or internal error. Distinct from a validation failure.
	ToolError = 2
)

// With calls os.Exit with the given code.
func With(code int) {
	os.Exit(code)
}

// IfErr calls os.Exit(ToolError) when err is non-nil.
func IfErr(err error) {
	if err != nil {
		os.Exit(ToolError)
	}
}

// IfValidationErr calls os.Exit(Validation) when err is non-nil.
func IfValidationErr(err error) {
	if err != nil {
		os.Exit(Validation)
	}
}
