# substrate

Core library for [Formulary](https://github.com/Formulary-Labs) compliance micro-tools.

Every Formulary tool imports `substrate`. No tool imports another tool.

```
go get github.com/Formulary-Labs/substrate
```

---

## Packages

### `substrate/exit`

Exit code constants for the three-value contract all Formulary tools must honor.

| Code | Constant | Meaning |
|---|---|---|
| `0` | `exit.OK` | Clean run — all checks passed, output is valid |
| `1` | `exit.Validation` | Validation failure — data did not meet requirements |
| `2` | `exit.ToolError` | Unexpected failure — missing file, network error, etc. |

```go
import "github.com/Formulary-Labs/substrate/exit"

if err != nil {
    fmt.Fprintf(os.Stderr, `{"error": %q, "code": 2}`, err)
    exit.With(exit.ToolError)
}
```

### `substrate/format`

Output format types and helpers. Default format is `json`.

```go
import "github.com/Formulary-Labs/substrate/format"

f, err := format.Parse(flagValue)  // validates user input
fmt.Println(f.MIMEType())          // "application/json"
fmt.Println(f.Extension())         // "json"
```

### `substrate/flags`

Standard CLI flag definitions. Every tool must expose these flags.

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
// cfg.Format, cfg.Program, cfg.DryRun, cfg.Interactive, cfg.Quiet
```

### `substrate/provenance`

Provenance writer. All tool writes must call `provenance.Write`. Entries are
appended to a JSONL file compatible with the `prompt` repo's `scripts/provenance_log.py`.

```go
import "github.com/Formulary-Labs/substrate/provenance"

err := provenance.Write("logs/provenance.jsonl", provenance.Entry{
    Spec:        "functions/probe-spec.md",
    Output:      "artifact.yaml",
    OutputType:  "artifact",
    Program:     cfg.Program,
    Purpose:     "Validated gemara artifact",
    Reusability: provenance.Instance,
    QualityGate: provenance.Pass,
    Tool:        "probe",
    ToolVersion: version.Version,
})
```

### `substrate/artifact`

Gemara artifact loaders. Type aliases for all key gemara types — import
`artifact` instead of `go-gemara` directly.

```go
import "github.com/Formulary-Labs/substrate/artifact"

catalog, err := artifact.LoadControlCatalog("path/to/catalog.yaml")
guidance, err := artifact.LoadGuidanceCatalog("path/to/guidance.yaml")
auditLog, err := artifact.LoadAuditLog("path/to/audit.yaml")
t, err := artifact.DetectType("path/to/unknown.yaml")
```

---

## Conventions

All Formulary tools must follow the shared conventions in
[CONTRIBUTING.md](https://github.com/Formulary-Labs/.github/blob/main/CONTRIBUTING.md).

Key points:
- `--format json` is the default — machine-readable output first
- Exit codes: 0/1/2 only — use `substrate/exit`
- No interactive prompts without `--interactive`
- `--dry-run` on all write operations
- All writes log provenance via `substrate/provenance`

---

## License

Apache License 2.0
