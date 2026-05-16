# CodeDNA 🧬

**Automated Engineering Standard Extraction Engine**

CodeDNA is a forensic tool designed to ingest an existing codebase, extract implicit engineering patterns, and generate a machine-readable "constitution" (like `.cursorrules` or `CLAUDE.md`) for AI agents to follow.

## 🚀 Key Features

- **AST-Powered Analysis:** Uses Tree-sitter for deep static analysis of Go, TypeScript, and more.
- **Privacy-First:** Context scrubbing and local-first parsing ensure your proprietary logic stays on your machine.
- **Semantic Clustering:** Identifies repeating patterns (logging, error handling, DI) using local embeddings.
- **Multi-Target Export:** Generates optimized rules for Cursor, VS Code, GitHub Copilot, and Continue.dev.
- **Premium TUI:** A sleek, high-end terminal interface for real-time extraction monitoring.

- **CLI:** Cobra

## 🌍 Supported Languages & Integrations

| Category | Supported Targets |
| :--- | :--- |
| **Languages** | Go, TypeScript, Python (Coming Soon), Rust (Coming Soon) |
| **IDEs** | Cursor, VS Code, JetBrains, Windsurf, Claude Code |
| **AI Agents** | Antigravity, Copilot, Claude, Cursor, Continue, Windsurf, JetBrains AI |
| **Rule Formats** | Modular Rules (`.agents/`, `.github/`, `.claude/`, `.cursor/`, `.continue/`, `.windsurf/`, `.aiassistant/`) |

## 📦 Installation

```bash
go install github.com/everglowlabs/codedna/cmd/codedna@latest
```

## 📖 Usage

```bash
# Scan a project and generate rules
codedna scan ./path/to/project
```

## 🛡 Security & Privacy

CodeDNA is built with privacy in mind. It performs all AST traversal locally and scrubs sensitive identifiers before any LLM processing (if enabled).

## 📄 License

Open Source under the MIT License.
