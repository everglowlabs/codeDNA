# CodeDNA — Implementation Specification

> **Automated Engineering Standard Extraction Engine**
> Transform a codebase's "tribal knowledge" into machine-readable AI agent constitutions.

---

## Table of Contents

1. [Overview](#1-overview)
2. [Architecture](#2-architecture)
3. [Internal Packages](#3-internal-packages)
4. [DNA Schema Reference](#4-dna-schema-reference)
5. [CLI Reference](#5-cli-reference)
6. [Output Formats](#6-output-formats)
7. [Security & Privacy Model](#7-security--privacy-model)
8. [TUI Experience](#8-tui-experience)
9. [Implementation Status](#9-implementation-status)
10. [Advanced Features & Roadmap](#10-advanced-features--roadmap)

---

## 1. Overview

CodeDNA ingests an existing codebase, performs deep AST-level analysis, and produces a structured "constitution" — a set of enforced engineering standards that AI agents (Cursor, Copilot, Claude, Continue, Windsurf, JetBrains, and Antigravity) natively consume.

**The core insight:** Every mature codebase implicitly enforces patterns — how errors are wrapped, how loggers are initialized, how HTTP handlers are structured. CodeDNA makes those patterns *explicit* and *machine-readable* in under a minute.

### Key Capabilities

| Capability | Description |
| :--- | :--- |
| **AST-Powered Analysis** | Tree-sitter parses Go, TypeScript, JavaScript, and Python at the syntax tree level — no regex heuristics |
| **Semantic Clustering** | Local vector embeddings group similar code blocks (e.g., all HTTP handlers, all logger init patterns) |
| **Privacy-First** | All parsing is 100% local; scrubbing removes secrets before any LLM call |
| **Multi-Target Export** | One scan → rules for 7+ AI agents and IDEs simultaneously |
| **CI/CD Lint Mode** | `codedna lint` validates PRs against the generated `codedna.json` baseline |

---

## 2. Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        codedna scan ./                       │
└───────────────────────────┬─────────────────────────────────┘
                            │
              ┌─────────────▼──────────────┐
              │    Scanner (fs walker)      │
              │  .gitignore-aware           │
              └─────────────┬──────────────┘
                            │ file paths
              ┌─────────────▼──────────────┐
              │    Parser (Tree-sitter)     │
              │  Go / TS / JS / Python      │
              └─────────────┬──────────────┘
                            │ []CodeBlock
              ┌─────────────▼──────────────┐
              │    Optimizer               │
              │  Context selection &        │
              │  deduplication              │
              └─────────────┬──────────────┘
                            │ representative blocks
              ┌─────────────▼──────────────┐
              │    Cluster                  │
              │  chromem-go embeddings      │
              │  groups by pattern type     │
              └─────────────┬──────────────┘
                            │ cluster groups
              ┌─────────────▼──────────────┐
              │  LLM Reasoning Pipeline     │
              │  langchaingo (CoT prompts)  │
              │  Anthropic / Gemini / OAI   │
              └─────────────┬──────────────┘
                            │ DNA_Schema (JSON)
              ┌─────────────▼──────────────┐
              │    Rule Sandbox             │
              │  Validates generated rules  │
              │  against existing AST       │
              └─────────────┬──────────────┘
                            │ validated standards
              ┌─────────────▼──────────────┐
              │  Multi-Format Compilers     │
              │  Cursor / Copilot / Claude  │
              │  Antigravity / Continue /   │
              │  Windsurf / JetBrains       │
              └─────────────────────────────┘
```

### Tech Stack

| Layer | Library | Purpose |
| :--- | :--- | :--- |
| **CLI** | `spf13/cobra` | Command routing (`scan`, `lint`) |
| **TUI** | `charmbracelet/bubbletea` + `lipgloss` + `bubbles` | Real-time forensic terminal UI |
| **AST** | `smacker/go-tree-sitter` | Polyglot syntax tree parsing |
| **Embeddings** | `philippgille/chromem-go` | Local vector DB for semantic clustering |
| **LLM** | `tmc/langchaingo` | Multi-provider orchestration (Anthropic, Gemini, OpenAI) |

---

## 3. Internal Packages

### `internal/scanner`
Walks the target directory respecting `.gitignore` rules. Returns a filtered list of file paths ready for parsing.

```
scanner.NewScanner(path) → Scanner
scanner.Scan(ignoreList) → []string (file paths)
```

**Key behavior:** Files matching `.gitignore` patterns, vendor directories, and binary files are excluded automatically.

---

### `internal/parser`
Wraps Tree-sitter for polyglot AST extraction. Each file extension maps to a language-specific grammar and query set.

```go
parser.NewParser(ext) → *Parser   // ext: ".go", ".ts", ".tsx", ".js", ".py"
parser.ExtractBlocks(path) → []CodeBlock
```

**Supported constructs per language:**

| Language | Extensions | Captured Constructs |
| :--- | :--- | :--- |
| **Go** | `.go` | `function_declaration`, `method_declaration` |
| **TypeScript** | `.ts`, `.tsx` | `function_declaration`, `method_definition`, `interface_declaration`, `class_declaration` |
| **JavaScript** | `.js` | `function_declaration`, `method_definition`, `class_declaration` |
| **Python** | `.py` | `function_definition`, `class_definition` |

Each `CodeBlock` carries:
- `Content` — the full source text of the construct
- `Kind` — `function`, `method`, `class`, or `interface`
- `FilePath` — the originating file

---

### `internal/cluster`
Groups `CodeBlock` slices into semantic clusters using `chromem-go` local vector embeddings. Similar patterns (e.g., all structured logging calls) are grouped together so the LLM receives coherent, topically-focused batches.

---

### `internal/optimizer`
Performs semantic deduplication and structural ranking. Selects "Golden Samples" — the most representative `CodeBlock` instances per cluster — to stay within LLM token budgets while maximizing pattern coverage.

---

### `internal/sandbox`
The Rule Sandbox validates generated standards against the project's actual AST. After the LLM produces a rule, the sandbox checks that a synthetic code example following that rule would be structurally consistent with existing project patterns. Rules that fail validation are flagged for refinement.

```go
sandbox.NewSandbox()
sandbox.LintFile(filePath, []Standard) → LintResult{Passed bool, Issues []string}
```

---

### `internal/compiler`
Implements the `Compiler` interface for each target agent. Each compiler consumes a `DNA_Schema` and emits a `map[string]string` of `filePath → content`.

```go
type Compiler interface {
    Compile(dna schema.DNA_Schema) (map[string]string, error)
}
```

**Available compilers:**

| Compiler | Output Target |
| :--- | :--- |
| `JSONCompiler` | `codedna.json` (raw schema, used by `lint`) |
| `MarkdownCompiler` | `CODEDNA_STANDARDS.md` (human-readable) |
| `CursorCompiler` | `.cursor/rules/<title>.mdc` |
| `CopilotCompiler` | `.github/instructions/<title>.instructions.md` |
| `AntigravityCompiler` | `.agents/rules/<title>.md` (≤ 12,000 chars) |
| `ClaudeCompiler` | `CLAUDE.md` + `.claude/rules/<title>.md` |
| `ContinueCompiler` | `.continue/rules/<title>.md` |
| `WindsurfCompiler` | `.windsurf/rules/<title>.md` |
| `JetBrainsCompiler` | `.aiassistant/rules/<title>.md` |

---

### `internal/schema`

The canonical `DNA_Schema` data model (see §4 for full reference).

---

### `internal/tui`

BubbleTea model implementing the forensic terminal UI. Provides real-time progress tracking, file-by-file status updates, and a structured analytics panel.

---

## 4. DNA Schema Reference

The `codedna.json` file is the source of truth for all compilers and the lint command.

```jsonc
{
  "project_name": "my-api",
  "version": "1.0.0",
  "generated_at": "2026-05-17T00:00:00Z",
  "standards": [
    {
      "id": "std-001",
      "category": "Error Handling",
      "title": "Wrapped Errors with Context",
      "description": "All errors must be wrapped with fmt.Errorf using the %w verb to preserve the error chain.",
      "rationale": "Error wrapping enables callers to use errors.Is / errors.As and preserves the full stack of context through the call chain.",
      "patterns": [
        "fmt.Errorf(\"...: %w\", err)",
        "errors.Is(err, target)"
      ],
      "samples": [
        {
          "file_path": "internal/db/postgres.go",
          "language": "go",
          "content": "func (r *Repo) GetUser(id int) (*User, error) {\n  u, err := r.db.QueryRow(id)\n  if err != nil {\n    return nil, fmt.Errorf(\"GetUser %d: %w\", id, err)\n  }\n  return u, nil\n}"
        }
      ],
      "rules": [
        "Always wrap errors with fmt.Errorf(\"context: %w\", err) — never return bare errors.",
        "Use errors.Is for comparison, never == on error values.",
        "Log errors at the boundary where they are handled, not where they are created."
      ]
    }
  ]
}
```

### Schema Fields

| Field | Type | Description |
| :--- | :--- | :--- |
| `project_name` | `string` | Project identifier |
| `version` | `string` | Schema version |
| `generated_at` | `time.Time` | ISO 8601 timestamp |
| `standards[].id` | `string` | Unique standard identifier |
| `standards[].category` | `string` | e.g., `"Error Handling"`, `"Logging"`, `"Testing"` |
| `standards[].title` | `string` | Short, human-readable name (used as filename) |
| `standards[].description` | `string` | One-sentence summary |
| `standards[].rationale` | `string` | *Why* this pattern exists |
| `standards[].patterns` | `[]string` | Abstracted code patterns detected |
| `standards[].samples` | `[]Sample` | Golden code samples with language and file path |
| `standards[].rules` | `[]string` | Concrete directives for AI agents |

---

## 5. CLI Reference

### Installation

```bash
go install github.com/everglowlabs/codedna/cmd/codedna@latest
```

---

### `codedna scan`

Scans a directory, runs the extraction pipeline, and exports agent rules.

```
Usage:
  codedna scan [path] [flags]

Flags:
  -f, --format strings   Output formats (default [json])
                         Choices: json, markdown, cursor, copilot,
                         antigravity, claude, continue, windsurf, jetbrains

Examples:
  # Scan current directory, output JSON only (default)
  codedna scan .

  # Scan and generate rules for Cursor + GitHub Copilot
  codedna scan . --format cursor,copilot

  # Scan a specific repo and export for all supported agents
  codedna scan ~/projects/my-api --format json,cursor,copilot,antigravity,claude,continue,windsurf,jetbrains

  # Scan and generate a human-readable standards document
  codedna scan . --format markdown
```

**What happens during a scan:**

1. The TUI launches in fullscreen forensic mode.
2. The scanner walks the directory tree, respecting `.gitignore`.
3. Tree-sitter extracts AST blocks from each supported file.
4. The optimizer selects golden samples within token limits.
5. The cluster engine groups semantically similar blocks.
6. The LLM reasoning pipeline (Chain-of-Thought) generates `Standard` entries.
7. The sandbox validates each rule against the existing AST.
8. Compilers write the selected output formats to disk.

---

### `codedna lint`

Checks that the current codebase adheres to a previously generated `codedna.json`. Designed for CI/CD pipelines.

```
Usage:
  codedna lint [path] [flags]

Flags:
  -c, --config string   Path to codedna.json (default "codedna.json")

Examples:
  # Lint current directory against codedna.json in the repo root
  codedna lint .

  # Lint with a custom config path
  codedna lint . --config ./standards/codedna.json

  # Lint a subdirectory
  codedna lint ./services/auth --config codedna.json
```

**Exit codes:**

| Code | Meaning |
| :--- | :--- |
| `0` | All files adhere to DNA standards |
| `1` | One or more violations found |

**Example CI step (GitHub Actions):**

```yaml
- name: CodeDNA Lint
  run: codedna lint . --config codedna.json
```

---

## 6. Output Formats

### Cursor (`.cursor/rules/<title>.mdc`)

Uses YAML frontmatter with glob scoping so rules activate only for relevant files:

```markdown
---
description: Wrapped Errors with Context
globs: internal/**/*.go, cmd/**/*.go
alwaysApply: false
---

# Wrapped Errors with Context

All errors must be wrapped with fmt.Errorf using the %w verb.

## Guidelines
- Always wrap errors with fmt.Errorf("context: %w", err).
- Use errors.Is for comparison, never == on error values.
```

---

### GitHub Copilot (`.github/instructions/<title>.instructions.md`)

Path-specific instructions with YAML frontmatter:

```markdown
---
globs:
  - "internal/**/*.go"
  - "cmd/**/*.go"
---

# Wrapped Errors with Context

All errors must be wrapped with fmt.Errorf using the %w verb.

## Rationale
Error wrapping enables callers to use errors.Is / errors.As...

## Guidelines
- Always wrap errors with fmt.Errorf("context: %w", err).
```

---

### Claude (`.claude/rules/<title>.md` + `CLAUDE.md`)

Generates a lean `CLAUDE.md` index and modular rule files using `paths:` frontmatter:

```markdown
---
paths:
  - "internal/**/*.go"
---

# Wrapped Errors with Context
...
```

`CLAUDE.md` contains an overview linking to each modular rule file.

---

### Antigravity (`.agents/rules/<title>.md`)

Modular rule files with a hard 12,000-character cap (Antigravity's enforced limit). Includes `Category`, `Rationale`, and code examples:

```markdown
# Wrapped Errors with Context

Category: Error Handling
Rationale: Error wrapping enables callers to use errors.Is / errors.As...

## Guidelines
- Always wrap errors with fmt.Errorf("context: %w", err).

## Examples
```go
// File: internal/db/postgres.go
func (r *Repo) GetUser(id int) (*User, error) { ... }
```
```

---

### Other Formats

| Format | Path Pattern | Frontmatter Key |
| :--- | :--- | :--- |
| Continue.dev | `.continue/rules/<title>.md` | `globs:`, `alwaysApply:` |
| Windsurf | `.windsurf/rules/<title>.md` | `trigger: glob`, `globs:` |
| JetBrains AI | `.aiassistant/rules/<title>.md` | `type: file patterns`, `pattern:` |
| Human Docs | `CODEDNA_STANDARDS.md` | — (plain Markdown) |
| Raw Schema | `codedna.json` | — (used by `lint`) |

---

## 7. Security & Privacy Model

### What stays local (always)

- Full source code never leaves your machine.
- AST parsing and block extraction are 100% local via Tree-sitter.
- Vector embeddings and clustering run in-process via `chromem-go`.

### What gets sent to the LLM (if enabled)

- Anonymized, deduplicated code *structures* — not raw proprietary strings.
- The scrubber replaces string literals, comments (optional), and high-entropy identifiers before payload construction.
- Zero-retention API headers are set by default (`anthropic-no-log: true`, etc.).

### Privacy settings

```bash
# Future flag — scrub all string literals before LLM submission
codedna scan . --scrub-strings

# Future flag — use a local LLM (Ollama) for fully offline operation
codedna scan . --provider ollama --model codellama
```

---

## 8. TUI Experience

The terminal UI is built with BubbleTea and Lipgloss, designed to feel like a high-end forensic analysis tool.

```
 ┌──────────────────────────────────────────────────────────────────┐
 │  CodeDNA 🧬  //  Standard Extraction Engine                      │
 ├──────────────────────────────────────────────────────────────────┤
 │  Scanning: [██████████████████░░░░░░░░░░] 64%  (42/66 files)    │
 │  Current:  ./internal/db/postgres.go                             │
 ├───────────────────────┬──────────────────────────────────────────┤
 │  PATTERNS DETECTED    │  LOGS                                    │
 │  ─────────────────    │  ─────────────────────────────────────   │
 │  [✔] Error Handling   │  ✔ Extracted 12 function blocks          │
 │  [✔] Logging (JSON)   │  ✔ Cluster: 3 HTTP handler patterns      │
 │  [➔] DI Strategy...   │  ➔ Generating rule: "Wrapped Errors"...  │
 │  [ ] Testing Patterns │                                          │
 └───────────────────────┴──────────────────────────────────────────┘
```

**TUI controls:**

| Key | Action |
| :--- | :--- |
| `q` / `Ctrl+C` | Abort scan |
| Auto-exits | On scan completion |

---

## 9. Implementation Status

### Phase 1 — Foundation ✅ COMPLETE

- [x] **CLI/TUI Skeleton** — Cobra command routing (`scan`, `lint`) with BubbleTea real-time file walker
- [x] **Tree-sitter Integration** — Go, TypeScript, JavaScript, and Python parsers with AST query extraction
- [x] **Filter Logic** — `.gitignore`-aware file scanner via `internal/scanner`

### Phase 2 — Intelligence & Extraction ✅ COMPLETE

- [x] **Pattern Clustering** — `chromem-go` local embeddings group code blocks by semantic similarity
- [x] **Context Payload Generation** — Optimizer selects representative "Golden Samples" within token limits
- [x] **JSON Schema Definition** — `DNA_Schema`, `Standard`, and `Sample` types in `internal/schema`

### Phase 3 — Validation & Output ✅ COMPLETE

- [x] **Rule Sandbox** — `internal/sandbox` validates generated rules against existing project AST patterns
- [x] **Multi-Format Compilers** — 9 compilers: JSON, Markdown, Cursor, Copilot, Antigravity, Claude, Continue, Windsurf, JetBrains
- [x] **CI/CD Integration** — `codedna lint` command with non-zero exit on violations

### Phase 4 — Local-First Intelligence 🔄 NEXT

- [ ] **Local LLM Support** — `--provider` flag supporting Ollama and LM Studio for fully offline operation
- [ ] **Interactive Rule Editor** — In-TUI accept / reject / edit flow before rules are written to disk
- [ ] **Expanded Language Support** — Rust (`.rs`) and Java (`.java`) Tree-sitter grammars
- [ ] **`--scrub-strings` flag** — Explicit opt-in to strip all string literals from the LLM payload
- [ ] **Severity Levels** — Classify each rule as `error`, `warning`, or `info` in the schema and lint output

### Phase 5 — Developer Experience 📋 PLANNED

- [ ] **Semantic Rule Diffing** — `codedna diff` compares two `codedna.json` snapshots and reports standard drift
- [ ] **Watch Mode** — `codedna watch` continuously lints files on save during active development
- [ ] **`codedna init`** — Guided setup wizard that detects the project language stack and pre-configures the scan
- [ ] **Rule Severity in CI** — `--fail-on warning` flag to make the lint step configurable by severity threshold
- [ ] **Structured Lint Output** — `--output json` flag on `lint` for machine-readable violation reports (parseable by CI dashboards)

### Phase 6 — Enterprise & Ecosystem 🔭 FUTURE

- [ ] **Multi-Repo Consensus** — `codedna consensus ./services/...` scans multiple repos to derive a "Company DNA"
- [ ] **First-Class GitHub Action** — `everglowlabs/codedna-action` for zero-config CI linting via the GitHub Marketplace
- [ ] **VS Code Extension** — Inline rule highlighting and quick-fix suggestions powered by the local `codedna.json`
- [ ] **Rule Registry** — A shareable, versioned registry of community-contributed `DNA_Schema` templates (e.g., `codedna use go-microservice-standard`)
- [ ] **Incremental Scanning** — Git-diff-aware mode that only re-analyzes changed files, reducing scan time on large repos

### Language Support Matrix

| Language | Extension(s) | AST Parsing | Clustering | Status |
| :--- | :--- | :--- | :--- | :--- |
| Go | `.go` | ✅ | ✅ | Stable |
| TypeScript | `.ts`, `.tsx` | ✅ | ✅ | Stable |
| JavaScript | `.js` | ✅ | ✅ | Stable |
| Python | `.py` | ✅ | ✅ | Stable |
| Rust | `.rs` | ⬜ | ⬜ | Phase 4 |
| Java | `.java` | ⬜ | ⬜ | Phase 4 |
| Ruby | `.rb` | ⬜ | ⬜ | Phase 6 |
| C/C++ | `.c`, `.cpp` | ⬜ | ⬜ | Phase 6 |

---

## 10. Phase Details

### Phase 4 — Local-First Intelligence

#### Local LLM Support

Introduces a `--provider` flag to route LLM calls to a locally running Ollama or LM Studio instance, enabling fully offline and air-gapped operation.

```bash
# Use a local Ollama model
codedna scan . --provider ollama --model codellama:13b

# Use LM Studio's OpenAI-compatible endpoint
codedna scan . --provider lmstudio --model mistral-7b

# Default cloud provider (unchanged)
codedna scan . --provider anthropic --model claude-3-5-sonnet
```

**Implementation notes:**
- Extend `langchaingo` provider selection in the scan pipeline
- Add `--provider` and `--model` flags to `scanCmd` in `cmd/codedna/commands/scan.go`
- Document recommended local models by language (e.g., CodeLlama for Go/Python, Mistral for TypeScript)

---

#### Interactive Rule Editor

After the LLM generates each `Standard`, the TUI pauses and presents it for user review before writing anything to disk.

```
 ┌─────────────────────────────────────────────────────────────────┐
 │  NEW RULE DETECTED                                              │
 │  ─────────────────────────────────────────────────────────────  │
 │  Title:    Wrapped Errors with Context                          │
 │  Category: Error Handling                                       │
 │  Rules:                                                         │
 │    • Always wrap with fmt.Errorf("context: %w", err)           │
 │    • Use errors.Is — never == on error values                   │
 │                                                                 │
 │  [A] Accept   [E] Edit   [R] Reject   [S] Skip All             │
 └─────────────────────────────────────────────────────────────────┘
```

**Implementation notes:**
- New `ReviewMode` state in the BubbleTea model in `internal/tui/model.go`
- `[E] Edit` opens an inline text editor (using `bubbles/textarea`) pre-filled with the rule list
- Accepted rules are queued; rejected rules are discarded before the compiler runs
- `--no-interactive` flag bypasses the editor for CI use

---

#### Severity Levels

Extends the `Standard` schema with a `Severity` field so linting can be configured by threshold.

```jsonc
// codedna.json
{
  "standards": [
    {
      "id": "std-001",
      "severity": "error",   // "error" | "warning" | "info"
      "title": "Wrapped Errors with Context",
      ...
    }
  ]
}
```

```bash
# Fail CI only on error-severity violations (warnings are logged but non-blocking)
codedna lint . --fail-on error

# Fail on warnings and above (default after this phase)
codedna lint . --fail-on warning
```

---

### Phase 5 — Developer Experience

#### `codedna diff`

Compares two `codedna.json` snapshots to surface how engineering standards have evolved between commits or sprints.

```bash
# Diff between a baseline and the current state
codedna diff --base codedna.v1.json --head codedna.json

# Diff between two git refs (auto-fetches the file at each ref)
codedna diff --base main --head feature/new-auth
```

**Sample output:**
```
📊 CodeDNA Standard Drift Report

[+] NEW    std-007  API Rate Limiting Pattern
[~] CHANGED std-002  Error Handling
           rule 2 changed: "log at boundary" → "log and emit metric at boundary"
[-] REMOVED std-005  Global Mutex Pattern  (no longer detected in codebase)
```

**Implementation notes:**
- New `cmd/codedna/commands/diff.go` with `--base` and `--head` flags
- Compute symmetric diff on `Standards` slices by `ID` and content hash
- Output modes: terminal (colored), JSON (`--output json`), and Markdown (`--output markdown`) for PR comments

---

#### `codedna watch`

Continuously monitors the working directory and lints changed files against the local `codedna.json` on every save.

```bash
codedna watch .
codedna watch . --config codedna.json --fail-on warning
```

**Sample output (live):**
```
👁  Watching ./  (codedna.json v1.2.0)

[12:04:01] CHANGED  internal/auth/handler.go
           [WARN] std-002: Error returned without context wrapping (line 47)
[12:04:33] CHANGED  internal/auth/handler.go
           [OK]  All standards pass.
```

**Implementation notes:**
- Use `fsnotify` to watch for `WRITE` events on supported extensions
- Run `sandbox.LintFile` on each changed path, print inline results
- Debounce rapid saves with a 200 ms window

---

#### Structured Lint Output

Makes `codedna lint` usable as a data source for CI dashboards, GitHub annotations, and custom tooling.

```bash
# Emit machine-readable JSON violation report
codedna lint . --output json > violations.json
```

```jsonc
// violations.json
{
  "passed": false,
  "violations": [
    {
      "file": "internal/auth/handler.go",
      "line": 47,
      "standard_id": "std-002",
      "severity": "error",
      "message": "Error returned without context wrapping."
    }
  ],
  "summary": { "error": 1, "warning": 0, "info": 0 }
}
```

---

### Phase 6 — Enterprise & Ecosystem

#### Multi-Repo Consensus (`codedna consensus`)

Scans multiple repositories and synthesizes a "Company DNA" — a superset schema representing standards that appear consistently across the org.

```bash
# Scan all services in a monorepo
codedna consensus ./services/*

# Scan separate repos (pass multiple paths)
codedna consensus ~/org/api ~/org/worker ~/org/gateway --format json,cursor

# Set a quorum threshold: only include standards present in ≥ 70% of repos
codedna consensus ./services/* --quorum 0.7
```

**Implementation notes:**
- Run the existing scan pipeline on each path in parallel (bounded goroutine pool)
- Merge `DNA_Schema` slices by `category` + pattern similarity (cosine similarity > 0.85)
- Output a single unified `company-dna.json` with a `sources[]` field listing contributing repos
- `--quorum` flag controls the minimum fraction of repos a pattern must appear in to be included

---

#### Rule Registry (`codedna use`)

A shareable, versioned registry of community-contributed `DNA_Schema` templates. Lets teams bootstrap standards without an initial scan.

```bash
# Apply a community-contributed standard template
codedna use go-microservice-standard
codedna use typescript-react-standard --format cursor,copilot

# List available templates
codedna registry list

# Publish your own DNA to the registry
codedna registry publish --name my-org/api-standard --from codedna.json
```

**Implementation notes:**
- Registry hosted at `registry.codedna.dev` (or self-hostable via `--registry-url`)
- Templates are versioned `DNA_Schema` JSON files with a `registry.yaml` manifest
- `codedna use` fetches the template, merges with any existing local `codedna.json`, and re-runs the compilers

---

#### Incremental Scanning

Git-diff-aware mode that only re-analyzes files changed since the last scan, reducing runtime on large repos from minutes to seconds.

```bash
# Only re-analyze files changed since last codedna.json was generated
codedna scan . --incremental

# Compare against a specific git ref
codedna scan . --incremental --since main
```

**Implementation notes:**
- Run `git diff --name-only <ref>` to get the changed file list
- Load existing `codedna.json`, retain unchanged standards, re-run the pipeline only on the diff
- Merge new/updated standards back into the existing schema, preserving `ID` stability
- Store a `git_ref` field in `DNA_Schema` metadata for traceability

---

> **CodeDNA** transforms tribal knowledge into machine-readable standards — one scan at a time.
