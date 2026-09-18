// Package flags defines the shared CLI flag configuration for all Formulary
// tools. Every tool must accept these flags. The Config struct captures their
// values. Tools may add tool-specific flags on top of this baseline.
//
// Usage in a tool:
//
//	import "github.com/Formulary-Labs/substrate/flags"
//
//	func main() {
//	    cfg := flags.Parse()
//	    // cfg.Format, cfg.Program, cfg.DryRun, cfg.Interactive, cfg.Quiet
//	}
//
// With cobra:
//
//	flags.Bind(cmd.PersistentFlags())
//	cobra.OnInitialize(func() { cfg = flags.FromFlagSet(cmd.Flags()) })
package flags

import (
	goflag "flag"
	"fmt"
	"os"

	"github.com/Formulary-Labs/substrate/format"
)

// Config holds the values of the standard Formulary CLI flags.
type Config struct {
	// Format is the output format. Default: "json".
	Format format.Format

	// Program is the program slug used for provenance logging and file
	// resolution. Empty string means no program context.
	Program string

	// DryRun, when true, causes all write operations to print what would be
	// written without touching the filesystem.
	DryRun bool

	// Interactive, when true, enables interactive prompts. Default false.
	// Tools must not emit interactive prompts when this is false.
	Interactive bool

	// Quiet suppresses progress output. Only the final result is emitted.
	Quiet bool

	// OutputDir is the directory for output files. Default: current directory.
	OutputDir string
}

// Default returns a Config with default values.
func Default() Config {
	return Config{
		Format: format.JSON,
	}
}

// Parse parses standard Formulary flags from os.Args using the standard
// library flag package. Returns a Config with resolved values. Exits 2 on
// invalid flag input.
func Parse() Config {
	cfg := Default()
	var fmtStr string

	goflag.StringVar(&fmtStr, "format", string(format.JSON), "Output format: json (default), yaml, md, csv, html")
	goflag.StringVar(&cfg.Program, "program", "", "Program slug for provenance logging and file resolution")
	goflag.BoolVar(&cfg.DryRun, "dry-run", false, "Print what would be written without writing it")
	goflag.BoolVar(&cfg.Interactive, "interactive", false, "Enable interactive prompts")
	goflag.BoolVar(&cfg.Quiet, "quiet", false, "Suppress progress output")
	goflag.StringVar(&cfg.OutputDir, "output-dir", ".", "Directory for output files")

	goflag.CommandLine.SetOutput(os.Stderr)
	goflag.Parse()

	f, err := format.Parse(fmtStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "flag error: %v\n", err)
		os.Exit(2)
	}
	cfg.Format = f

	return cfg
}

// Descriptions returns the flag descriptions for documentation and help output.
func Descriptions() map[string]string {
	return map[string]string{
		"format":      "Output format: json (default), yaml, md, csv, html",
		"program":     "Program slug for provenance logging and file resolution",
		"dry-run":     "Print what would be written without writing it",
		"interactive": "Enable interactive prompts (default: off)",
		"quiet":       "Suppress progress output; only emit the final result",
		"output-dir":  "Directory for output files",
	}
}
