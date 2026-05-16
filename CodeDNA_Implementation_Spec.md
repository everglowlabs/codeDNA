# CodeDNA Implementation Spec: Automated Agent Standard Extraction

## 1. Objective
Build an automated engine that ingests an existing codebase, extracts implicit engineering patterns (architecture, error handling, logging, testing), and outputs a highly structured "constitution" (e.g., `.cursorrules`, `CLAUDE.md`, or `.github/copilot-instructions.md`) for AI agents and IDEs to natively follow.

---

## 2. Architecture & Tech Stack (Go-Centric)
To ensure the tool is extremely fast, easily distributable as a single binary, and capable of deep static analysis, the core extraction engine is built in **Go**.

### Core Modules
1.  **Provider-Based AST Engine:** A modular system for parsing codebases. Each language (Go, TypeScript, Python, etc.) has a dedicated provider that uses Tree-sitter to identify high-signal constructs (interfaces, decorators, wrapper functions).
2.  **The Context "Optimizer":** Instead of simple packing, this module performs semantic deduplication and structural ranking. It selects the "Golden Samples" of code that best represent the project's standards.
3.  **Privacy-Centric Payload Scrubbing:** Before transmission, CodeDNA scrubs the selected context of string literals, comments (optional), and sensitive-looking identifiers using entropy-based detection.
4.  **Reasoning Pipeline:** Utilizes a Chain-of-Thought (CoT) prompting strategy. The LLM is asked to:
    *   *Analyze* the provided samples.
    *   *Infer* the rationale behind the patterns.
    *   *Synthesize* the rule in a unified internal JSON schema.
5.  **Multi-Target Rule Compiler:** Translates the internal JSON schema into target-specific formats (Markdown, JSON, or YAML) optimized for different AI agents.

### Tech Stack & Libraries
*   **CLI Framework:** `spf13/cobra`
*   **TUI Framework:** `charmbracelet/bubbletea` + `lipgloss` + `bubbles` (for progress, spinners, and lists).
*   **AST Parsing:** `smacker/go-tree-sitter` for polyglot support.
*   **Local Vector DB:** `chromem-go` for semantic search of pattern "clusters."
*   **LLM Orchestration:** `tmc/langchaingo` for multi-provider support (Anthropic, Gemini, OpenAI).

---

## 3. Security & Privacy Posture (Privacy-First)
*   **Local-First Parsing:** AST traversal and structural extraction happen 100% locally.
*   **Zero-Retention Endpoints:** Default configurations use API headers that opt-out of training/retention.
*   **Anonymization:** A local "scrubber" ensures that proprietary business logic strings are replaced with generic tokens before leaving the machine.

---

## 4. Target Output Formats (Native IDE Integration)
CodeDNA acts as a "Universal Rule Translator":
*   **Cursor / Antigravity:** `.cursorrules` (Root-level rule enforcement).
*   **VS Code / GitHub Copilot:** `.github/copilot-instructions.md`.
*   **Custom Agents / Continue.dev:** `CLAUDE.md` or `DEVELOPER.md`.
*   **Internal Knowledge Base:** Markdown-formatted engineering standards for human onboarding.

---

## 5. User Interface (Premium TUI Experience)
The TUI is designed to feel like a high-end forensic tool.

```text
 ┌──────────────────────────────────────────────────────────────┐
 │ CodeDNA // Standard Extraction Engine                        │
 ├──────────────────────────────────────────────────────────────┤
 │ Scanning: [██████████████████░░░░░░░░░░░] 64%                │
 │ Currently processing: ./internal/db/postgres.go               │
 ├──────────────────────────────────────────────────────────────┤
 │ ANALYTICS                                                    │
 │ ├─ [✔] Logging Pattern: Structured JSON (Logrus detected)    │
 │ ├─ [✔] Error Handling: Wrapped (fmt.Errorf with %w)          │
 │ ├─ [➔] DI Strategy: Wire-based injection detected...          │
 ├──────────────────────────────────────────────────────────────┤
 │ LOGS: Rule "Database-Safe-Queries" generated via CoT...      │
 └──────────────────────────────────────────────────────────────┘
```

---

## 6. Implementation Roadmap

### Phase 1: The Foundation (Week 1) [COMPLETED]
*   **[x] CLI/TUI Skeleton:** Implement the Cobra/BubbleTea loop with a real-time file walker.
*   **[x] Tree-sitter Integration:** Setup Go and TypeScript parsers to extract basic function signatures and interface definitions.
*   **[x] Filter Logic:** Robust `.gitignore` awareness to prevent scanning noise.

### Phase 2: Intelligence & Extraction (Week 2)
*   **Pattern Clustering:** Use local embeddings to group similar code blocks (e.g., all HTTP handlers).
*   **Context Payload Generation:** Smart selection of representative files to stay within LLM token limits while maximizing context density.
*   **JSON Schema Definition:** Establish the `DNA_Schema` for internal rule representation.

### Phase 3: Validation & Output (Week 3)
*   **The "Rule Sandbox":** Have the LLM generate a dummy file following the new rules. If the local parser detects a deviation from existing project AST patterns, the rule is refined.
*   **Multi-Format Compilers:** Implementation of Cursor, Copilot, and Markdown exporters.
*   **CI/CD Integration:** A "Lint" mode that checks if new PRs violate the generated `CodeDNA`.

---

## 7. Advanced Features (Future Scope)
*   **Semantic Rule Diffing:** When running CodeDNA again, it shows you how your project standards have evolved over time.
*   **Multi-Repo Consensus:** Scan multiple microservices to generate a "Company DNA" that enforces consistency across the entire org.
*   **Interactive Refinement:** During the TUI session, the user can "reject" or "edit" a rule before it is written to disk.

---
**CodeDNA** transforms "tribal knowledge" into "machine-readable standards."
