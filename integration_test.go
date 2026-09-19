// Package integration runs smoke tests for the Formulary CLI tools.
// Each test builds the binary for its tool and verifies that --version
// exits 0 and emits a version string. These tests catch broken builds and
// missing main packages early in CI.
package integration_test

import (
	"os/exec"
	"strings"
	"testing"
)

// tool describes a single CLI binary to smoke-test.
type tool struct {
	name   string
	dir    string
	mainPkg string
}

var tools = []tool{
	{"assay",    "../assay",    "./cmd/assay"},
	{"bind",     "../bind",     "./cmd/bind"},
	{"challenge","../challenge","./cmd/challenge"},
	{"compound", "../compound", "./cmd/compound"},
	{"decay",    "../decay",    "./cmd/decay"},
	{"dose",     "../dose",     "./cmd/dose"},
	{"exhibit",  "../exhibit",  "./cmd/exhibit"},
	{"formula",  "../formula",  "./cmd/formula"},
	{"probe",    "../probe",    "./cmd/probe"},
	{"scan",     "../scan",     "./cmd/scan"},
	{"specimen", "../specimen", "./cmd/specimen"},
	{"titer",    "../titer",    "./cmd/titer"},
	{"vital",    "../vital",    "./cmd/vital"},
}

func TestVersionFlag(t *testing.T) {
	for _, tool := range tools {
		tool := tool
		t.Run(tool.name, func(t *testing.T) {
			t.Parallel()

			// Build the binary into a temp dir.
			bin := t.TempDir() + "/" + tool.name
			buildCmd := exec.Command("go", "build", "-o", bin, tool.mainPkg)
			buildCmd.Dir = tool.dir
			if out, err := buildCmd.CombinedOutput(); err != nil {
				t.Fatalf("build %s: %v\n%s", tool.name, err, out)
			}

			// Run --version and check it exits 0.
			runCmd := exec.Command(bin, "--version")
			out, err := runCmd.CombinedOutput()
			if err != nil {
				t.Fatalf("%s --version failed: %v\noutput: %s", tool.name, err, out)
			}
			if !strings.Contains(string(out), "version") {
				t.Errorf("%s --version output %q does not contain 'version'", tool.name, string(out))
			}
		})
	}
}
