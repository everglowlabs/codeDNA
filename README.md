# CodeDNA 🧬

> **Automated Engineering Standard Extraction Engine**
> Scan a codebase. Extract its implicit patterns. Generate AI agent rules — automatically.

---

CodeDNA forensically analyzes an existing codebase using AST-level static analysis, clusters the detected patterns semantically, and produces a machine-readable "constitution" that AI agents (Cursor, Copilot, Claude, Continue, Windsurf, JetBrains, and Antigravity) natively follow.

Every mature project implicitly enforces standards — how errors are wrapped, how loggers are initialized, how HTTP handlers are structured. CodeDNA makes those standards **explicit** and **enforceable** in under a minute.

---

## Features

- **AST-Powered Analysis** — Tree-sitter parses Go, TypeScript, JavaScript, Python, Rust, and Java at the syntax tree level. No regex heuristics.
- **Semantic Clustering** — Local vector embeddings (`chromem-go`) group similar code blocks so the LLM receives coherent, topically-focused batches.
- **Offline & Local LLM Support** — Run completely offline and air-gapped using local reasoning providers like Ollama and LM Studio.
- **Privacy-First** — All parsing runs 100% locally. Active string literal scrubbing (`--scrub-strings`) sanitizes payloads before reasoning.
- **Interactive Rule Editor** — Review, edit, or reject synthesized rules directly inside the fullscreen BubbleTea terminal dashboard.
- **Multi-Target Export** — One scan generates rules for 7+ AI agents and IDEs simultaneously.
- **Rule Sandbox** — Generated rules are validated against your existing AST before being written to disk.
- **CI/CD Lint Mode** — `codedna lint` exits non-zero on violations, offering configurable severity gates (`--fail-on`).

---

## Installation

```bash
go install github.com/everglowlabs/codedna/cmd/codedna@latest
```

**Requirements:** Go 1.21+

---

## Quick Start

```bash
# Scan the current directory and output codedna.json
codedna scan .

# Scan and generate rules for Cursor + GitHub Copilot
codedna scan . --format cursor,copilot

# Scan and export for every supported AI agent
codedna scan . --format json,cursor,copilot,antigravity,claude,continue,windsurf,jetbrains

# Generate a human-readable engineering standards document
codedna scan . --format markdown
```

---

## CLI Reference

### `codedna scan`

```
Usage:
  codedna scan [path] [flags]

Flags:
  -f, --format strings   Output format(s), comma-separated (default: [json])
                         Options: json, markdown, cursor, copilot,
                                  antigravity, claude, continue, windsurf, jetbrains
  -p, --provider string  LLM reasoning provider: anthropic, openai, gemini,
                         ollama, lmstudio (default: "anthropic")
  -m, --model string     LLM model name (default: "claude-3-5-sonnet")
      --scrub-strings    Scrub string literals from AST blocks before reasoning for privacy
      --no-interactive   Disable the interactive TUI rule review session
```

Running `scan` launches the forensic TUI, walks the directory tree (respecting `.gitignore`), runs the extraction pipeline, prompts for interactive rule reviews, and writes the selected output formats to disk.

### `codedna lint`

```
Usage:
  codedna lint [path] [flags]

Flags:
  -c, --config string    Path to codedna.json (default: "codedna.json")
      --fail-on string   Severity threshold to fail the build: error, warning, info
                         (default: "warning")
```

Validates the current codebase against a previously generated `codedna.json`. Exits `0` on success, or `1` on violations that meet or exceed the `--fail-on` threshold.

---

## Supported Languages & Output Formats

### Languages

| Language | Extensions | Status |
| :--- | :--- | :--- |
| Go | `.go` | ✅ Stable |
| TypeScript | `.ts`, `.tsx` | ✅ Stable |
| JavaScript | `.js` | ✅ Stable |
| Python | `.py` | ✅ Stable |
| Rust | `.rs` | ✅ Stable |
| Java | `.java` | ✅ Stable |

### Output Formats

| Format Flag | Output Path | Used By |
| :--- | :--- | :--- |
| `json` | `codedna.json` | Lint command, custom tooling |
| `markdown` | `CODEDNA_STANDARDS.md` | Human onboarding docs |
| `cursor` | `.cursor/rules/*.mdc` | Cursor IDE |
| `copilot` | `.github/instructions/*.instructions.md` | GitHub Copilot |
| `antigravity` | `.agents/rules/*.md` | Antigravity |
| `claude` | `CLAUDE.md` + `.claude/rules/*.md` | Claude Code |
| `continue` | `.continue/rules/*.md` | Continue.dev |
| `windsurf` | `.windsurf/rules/*.md` | Windsurf Cascade |
| `jetbrains` | `.aiassistant/rules/*.md` | JetBrains AI |

---

## CI/CD Integration

Add CodeDNA as a PR quality gate in GitHub Actions:

```yaml
name: CodeDNA Lint
on: [pull_request]

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      - run: go install github.com/everglowlabs/codedna/cmd/codedna@latest
      - run: codedna lint . --config codedna.json
```

Commit the `codedna.json` generated by your initial `scan` to the repository. The lint step will enforce those standards on every PR.

---

## How It Works

```
codedna scan ./
      │
      ▼
  Scanner ──── respects .gitignore
      │
      ▼
  Parser ───── Tree-sitter AST extraction (Go, TS, JS, Python, Rust, Java)
      │
      ▼
  Optimizer ── selects "Golden Samples" within LLM token limits
      │
      ▼
  Cluster ──── chromem-go local embeddings group similar patterns
      │
      ▼
  LLM ───────── Chain-of-Thought prompts → DNA_Schema (JSON)
      │
      ▼
  Sandbox ──── validates rules against existing AST
      │
      ▼
  Compilers ── writes rules for each target agent
```

The `codedna.json` produced at the end is the canonical source of truth. All other output formats are compiled from it.

---

## Security & Privacy

- **Local-first parsing:** AST traversal and structural extraction happen 100% on your machine.
- **Scrubbing:** String literals and high-entropy identifiers are removed before the LLM payload is constructed.
- **Zero-retention headers:** API calls default to opt-out-of-training headers (`anthropic-no-log: true`, etc.).
- **No source upload:** Only anonymized structural patterns reach the LLM — never raw source files.

---

## Architecture & Technical Details

See [`CodeDNA_Implementation_Spec.md`](./CodeDNA_Implementation_Spec.md) for the full technical specification, including:
- Internal package API reference
- `DNA_Schema` JSON schema with field descriptions
- Per-format output examples (Cursor, Copilot, Claude, Antigravity, etc.)
- Full implementation roadmap and status

---

## License

Open Source under the [MIT License](LICENSE).
