# ArchonCLI 🏛️

[![Go Version](https://img.shields.io/github/go-mod/go-version/rexreus/archon)](https://go.dev/)
[![License](https://img.shields.io/github/license/rexreus/archon)](LICENSE)
[![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20macOS%20%7C%20Linux-blue)](#)
[![Built with Gemini](https://img.shields.io/badge/Built%20with-Google%20Gemini%203-orange)](https://ai.google.dev/)

**ArchonCLI** is an autonomous AI-powered code architect assistant designed to help developers master complex codebases through deep semantic analysis and a modern terminal interface.

---

## ✨ Why ArchonCLI?

Traditional code assistants often struggle with large-scale projects due to fragmented context. ArchonCLI solves this by combining **Semantic Syntax-Aware Indexing** with the reasoning power of **Google Gemini 3**.

### 🚀 Key Features

*   **🧠 Syntax-Aware RAG**: Uses **Tree-Sitter** to understand code structure (functions, classes, types) instead of just raw text.
*   **💻 Dual Interface**: Choose between a visually rich **Interactive TUI** (Bubble Tea) or a scriptable **CLI** (Cobra).
*   **⚡ Smart Optimization**: Leverage **Context Caching** to reduce latency and API costs by up to 90%.
*   **🛠️ Developer Power Tools**:
    *   `review`: Automated AI code review for staged changes.
    *   `commit`: Generate Conventional Commit messages from your diffs.
    *   `test`: Instant unit test generation for any file.
    *   `diagram`: Generate Mermaid/PlantUML diagrams of your architecture.
    *   `refactor`: Get AI-driven refactoring suggestions or apply them directly.

---

## 🖥️ Modern Terminal Experience

ArchonCLI provides a beautiful **TUI (Terminal User Interface)** built with the [Charm](https://charm.sh/) ecosystem.

```text
  ArchonCLI - AI Architect Assistant
  ──────────────────────────────────
  Welcome, Architect. What would you like to do?
  
  > Chat Mode
    Index Codebase
    AI Code Review
    Smart Commit
    System Status
    ...
```

---

## 🚀 Quick Start

### 1. Install

**Automated (Recommended):**
You can use the provided installation scripts:

- **Linux/macOS**: `curl -sSL https://raw.githubusercontent.com/rexreus/archon/main/scripts/install.sh | bash`
- **Windows (PowerShell)**: `irm https://raw.githubusercontent.com/rexreus/archon/main/scripts/install.ps1 | iex`

**Manual Build:**
```bash
go build -o archon ./cmd/archon/main.go
```

### 2. Initialize & Authenticate
Run these commands in your project's root directory:

```bash
# Create config file
archon init

# Set your Gemini API Key
archon auth --key "YOUR_API_KEY"
```

### 3. Index & Ask
Index your codebase to enable semantic search, then start asking:

```bash
# Build the vector index
archon index

# Ask away!
archon ask "How is the data flow handled in this project?"
```

*Or just run `archon` to enter the **Interactive TUI Mode**.*

---

## 📖 Command Overview

| Command | Description |
| :--- | :--- |
| `index` | Scan and index the codebase (with `--watch` support). |
| `ask` | Ask general questions about your project. |
| `review` | Review staged changes for bugs and best practices. |
| `commit` | Generate and apply smart commit messages. |
| `test` | Generate unit tests for specific files. |
| `refactor` | Analyze and suggest improvements for code. |
| `explain` | Deep explanation of files or symbols. |
| `diagram` | Generate architecture diagram code. |
| `status` | View system health and token usage stats. |
| `uninstall` | Uninstall ArchonCLI and remove its data. |

---

## 🛠️ Tech Stack

- **Core**: [Go (Golang)](https://go.dev/)
- **AI Backend**: [Google Gemini 3](https://ai.google.dev/)
- **Vector DB**: [chromem-go](https://github.com/philippgille/chromem-go) (Pure Go Vector DB)
- **Parser**: [Tree-Sitter](https://tree-sitter.github.io/tree-sitter/) (Syntax-aware indexing)
- **TUI Framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea) & [Lip Gloss](https://github.com/charmbracelet/lipgloss)
- **CLI Framework**: [Cobra](https://github.com/spf13/cobra) & [Viper](https://github.com/spf13/viper)

---

## 📂 Project Structure & Docs

Detailed information is available in the [**docs/**](./docs/) directory:

- 📘 [**Introduction**](./docs/README.md) - Project summary.
- 🌟 [**Features**](./docs/FEATURES.md) - Deep dive into functionality.
- 📖 [**Usage Guide**](./docs/USAGE.md) - Command reference and setup.
- ⚙️ [**Configuration**](./docs/CONFIGURATION.md) - API Key and settings.
- 🏗️ [**Architecture**](./docs/ARCHITECTURE.md) - Technical implementation details.

---

## 🛡️ Security & Privacy

ArchonCLI respects your privacy. The **Vector Database** is stored locally on your machine. Data sent to Google Gemini is governed by Google's API terms, and we recommend using **Vertex AI Endpoints** for enterprise-grade security.

---

*ArchonCLI - Empowering Architects with AI Intelligence.*
