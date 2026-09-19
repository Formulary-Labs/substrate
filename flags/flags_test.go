package flags_test

import (
	goflag "flag"
	"testing"

	"github.com/Formulary-Labs/substrate/flags"
	"github.com/Formulary-Labs/substrate/format"
)

func TestDefault_hasJSONFormat(t *testing.T) {
	cfg := flags.Default()
	if cfg.Format != format.JSON {
		t.Errorf("Default format = %q, want %q", cfg.Format, format.JSON)
	}
	if cfg.DryRun {
		t.Error("Default DryRun must be false")
	}
	if cfg.Quiet {
		t.Error("Default Quiet must be false")
	}
	if cfg.Interactive {
		t.Error("Default Interactive must be false")
	}
}

func TestRegisterOn_registersFlags(t *testing.T) {
	fs := goflag.NewFlagSet("test", goflag.ContinueOnError)
	cfg, fmtStr := flags.RegisterOn(fs)

	// Flags should be registered without error.
	if err := fs.Parse([]string{"-program", "myprogram", "-dry-run"}); err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if cfg.Program != "myprogram" {
		t.Errorf("Program = %q, want %q", cfg.Program, "myprogram")
	}
	if !cfg.DryRun {
		t.Error("DryRun should be true after -dry-run flag")
	}
	if *fmtStr != string(format.JSON) {
		t.Errorf("fmtStr = %q, want %q", *fmtStr, string(format.JSON))
	}
}

func TestFinalize_resolvesFormat(t *testing.T) {
	fs := goflag.NewFlagSet("test", goflag.ContinueOnError)
	cfg, fmtStr := flags.RegisterOn(fs)

	if err := fs.Parse([]string{"-format", "md"}); err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if err := flags.Finalize(cfg, *fmtStr); err != nil {
		t.Fatalf("Finalize error: %v", err)
	}

	if cfg.Format != format.MD {
		t.Errorf("Format = %q after Finalize, want %q", cfg.Format, format.MD)
	}
}

func TestFinalize_invalidFormat(t *testing.T) {
	fs := goflag.NewFlagSet("test", goflag.ContinueOnError)
	cfg, _ := flags.RegisterOn(fs)
	if err := fs.Parse(nil); err != nil {
		t.Fatal(err)
	}

	err := flags.Finalize(cfg, "badformat")
	if err == nil {
		t.Error("expected error for invalid format string")
	}
}

func TestDescriptions_hasAllKeys(t *testing.T) {
	descs := flags.Descriptions()
	required := []string{"format", "program", "dry-run", "interactive", "quiet", "output-dir"}
	for _, key := range required {
		if _, ok := descs[key]; !ok {
			t.Errorf("Descriptions missing key %q", key)
		}
	}
}

func TestRegisterOn_outputDirDefault(t *testing.T) {
	fs := goflag.NewFlagSet("test", goflag.ContinueOnError)
	cfg, _ := flags.RegisterOn(fs)

	if err := fs.Parse(nil); err != nil {
		t.Fatal(err)
	}

	if cfg.OutputDir != "." {
		t.Errorf("OutputDir default = %q, want %q", cfg.OutputDir, ".")
	}
}
