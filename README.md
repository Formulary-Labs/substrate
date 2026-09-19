# substrate

Core library for Formulary compliance tools.

Every Formulary tool imports `substrate`. No tool imports another tool.

```bash
go get github.com/Formulary-Labs/substrate
```

## Packages

### `substrate/exit`

The three-value exit code contract that all Formulary tools must honor. Callers — CI pipelines, agent orchestrators, shell scripts — rely on this contract being consistent across every tool.

| Code | Constant | Meaning |
|---|---|---|
| `0` | `exit.OK` | Clean run — all checks passed, output is valid |
| `1` | `exit.Validation` | Validation failure — data did not meet requirements |
| `2` | `exit.ToolError` | Unexpected failure — missing file, network error, internal error |

```go
import "github.com/Formulary-Labs/substrate/exit"

// Fail with a tool error and a structured JSON message
if err != nil {
    fmt.Fprintf(os.Stderr, `{"error": %q, "code": 2}`, err)
    exit.With(exit.ToolError)
}

// Convenience helpers
exit.IfErr(err)            // exits ToolError if err != nil
exit.IfValidationErr(err)  // exits Validation if err != nil
```

Do not call `os.Exit` directly in Formulary tools. Always use `exit.With`, `exit.IfErr`, or `exit.IfValidationErr`.

### `substrate/format`

Output format types, parsing, and MIME/extension helpers. Default format is `json`.

```go
import "github.com/Formulary-Labs/substrate/format"

f, err := format.Parse(flagValue)  // validates user input; errors on unknown format
fmt.Println(f.MIMEType())          // "application/json"
fmt.Println(f.Extension())         // "json"
```

Supported format values: `json`, `md`, `csv`, `html`. Not every tool supports every format — `format.Parse` validates the value; the tool is responsible for handling the parsed format.

### `substrate/flags`

Standard CLI flag definitions. Every Formulary tool must expose these flags. Pulling them from `substrate/flags` keeps definitions consistent and prevents drift.

| Flag | Type | Default | Description |
|---|---|---|---|
| `--format` | string | `json` | Output format |
| `--program` | string | `""` | Program slug for provenance and file resolution |
| `--dry-run` | bool | `false` | Print what would be written without writing it |
| `--interactive` | bool | `false` | Enable interactive prompts |
| `--quiet` | bool | `false` | Suppress progress output |
| `--output-dir` | string | `.` | Directory for output files |

```go
import "github.com/Formulary-Labs/substrate/flags"

cfg := flags.Parse()
// cfg.Format, cfg.Program, cfg.DryRun, cfg.Interactive, cfg.Quiet, cfg.OutputDir
```

No tool may prompt for user input unless `cfg.Interactive` is true. This is required for CI and agent-orchestrated use.

### `substrate/provenance`

Provenance writer. Every tool write must call `provenance.Write`. Entries are appended to a JSONL file compatible with the `regimen` provenance log schema, so agent-layer tooling and Go tools share a single audit trail.

```go
import "github.com/Formulary-Labs/substrate/provenance"

err := provenance.Write("logs/provenance.jsonl", provenance.Entry{
    Spec:        "functions/control-assessment-spec.md",
    Output:      "output/product-iso27001-assessment.json",
    OutputType:  "assessment",
    Program:     cfg.Program,
    Purpose:     "Generated control assessment for ISO 27001 Q3 run",
    Reusability: provenance.Instance,
    QualityGate: provenance.Pass,
    Tool:        "assay",
    ToolVersion: version.Version,
})
```

`provenance.Write` creates the log directory and file if they do not exist. It never truncates. Do not open or write to `provenance.jsonl` directly.

Reusability constants: `Template`, `Reference`, `Instance`, `Artifact`.
Quality gate constants: `Pass`, `FailedOnceCorrected`, `Escalated`, `Skipped`, `NotApplicable`.

### `substrate/artifact`

Gemara artifact loaders. Import `artifact` instead of `go-gemara` directly so tools never duplicate the detection and loading logic.

```go
import "github.com/Formulary-Labs/substrate/artifact"

catalog, err := artifact.LoadControlCatalog("path/to/catalog.yaml")
guidance, err := artifact.LoadGuidanceCatalog("path/to/guidance.yaml")
auditLog, err := artifact.LoadAuditLog("path/to/audit.yaml")

// Auto-detect type when you don't know what you're loading
t, err := artifact.DetectType("path/to/unknown.yaml")
```

## Conventions

All Formulary tools must follow these conventions. `substrate` enforces most of them at the code level:

- `--format json` is the default — machine-readable output first
- Exit codes: `0/1/2` only — use `substrate/exit`
- No interactive prompts without `--interactive`
- `--dry-run` on all write operations
- All writes log provenance via `substrate/provenance`
- Data output to stdout, progress and errors to stderr
- Error output is structured JSON when `--format json` is set: `{"error": "...", "code": 1}`

## Building a new tool

If you're adding a tool to Formulary, `substrate` is the starting point. Import it, wire up `flags.Parse()`, use `exit.With` for exit handling, call `provenance.Write` on every write. The conventions checklist is in [CONTRIBUTING.md](https://github.com/Formulary-Labs/.github/blob/main/CONTRIBUTING.md).

## License

Apache License 2.0
