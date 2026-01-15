# ArchonCLI Usage Guide 📖

This document explains how to install, set up, and use the various commands available in ArchonCLI.

## 📥 Installation

Currently, ArchonCLI is available as a single binary. You can build it from source if you have Go installed:

```bash
# Clone the repository
git clone https://github.com/rexreus/archon.git
cd archon

# Build the binary
go build -o archon ./cmd/archon/main.go

# Move to PATH (Optional)
mv archon /usr/local/bin/
```

## 🛠️ Initial Setup

Before you start asking questions, you need to initialize and authenticate.

### 1. Project Initialization
Run this command in your project's root directory:
```bash
archon init
```
This will create a `.archon.yaml` configuration file.

### 2. Authentication
Enter your Google Gemini API Key:
```bash
archon auth --key "AIzaSy..."
```
Alternatively, set it via environment variable:
```bash
export ARCHON_GEMINI_KEY="AIzaSy..."
```

## ⌨️ CLI Command Reference

> **Note:** Di dokumentasi ini kami menggunakan perintah `archon`. Jika Anda belum melakukan build biner, gunakan `go run cmd/archon/main.go` sebagai gantinya.

### `archon index`
Scans the codebase and updates the local vector database.
- `--force`, `-f`: Force re-indexing of all files.
- `--watch`, `-w`: Monitor file changes in real-time.

### `archon ask [question]`
Ask a question about your code.
```bash
archon ask "Explain how the authentication system works here"
```

### `archon explain [file/symbol]`
Explain a specific file or symbol (function/class).
```bash
archon explain ./internal/core/orchestrator.go
archon explain --symbol "LoadConfig"
```

### `archon refactor [file]`
Analyze code and provide improvement suggestions.
```bash
# Contoh menggunakan path yang benar sesuai struktur proyek
archon refactor ./internal/adapters/parser/parser.go --goal "improve performance"
archon refactor ./internal/adapters/parser/parser.go --apply # Terapkan perubahan langsung
```

### `archon review`
Perform an automated code review on staged changes (`git add`).
```bash
archon review
```

### `archon commit`
Analyze staged changes and generate a smart commit message, with an option to commit immediately.
```bash
archon commit
```

### `archon test [file]`
Generate automated unit tests for the selected file.
```bash
archon test ./internal/core/orchestrator.go
```

### `archon analyze`
Perform a deep scan to detect code smells or design pattern violations.

### `archon diagram`
Generate diagram code (Mermaid/PlantUML).
```bash
archon diagram --type sequence --focus "OrderProcess"
```

### `archon doc [file]`
Automatically generate code documentation (docstrings/comments).

### `archon status`
Show system health status, vector index statistics, and API quota usage.

### `archon uninstall`
Removes the configuration file and local vector database. To remove the binary, use the uninstall script provided in the repository.

### `archon config [list/set]`
View or modify configuration settings.

### `archon lsp`
Start Language Server mode (for IDE integration).

### `archon version`
Display build version information.

## 🖥️ Using TUI Mode

Simply type `archon` without arguments to enter interactive mode.

### Navigation
- **Arrow Up/Down**: Move between menus or scroll chat history.
- **Enter**: Select menu or send message.
- **Tab**: Switch between panels (e.g., between Chat and Context panels).
- **Esc / Ctrl+C**: Return to menu or exit the application.

### TUI Features
- **Chat Mode**: Interactive discussion with AI.
- **AI Code Review**: Analyze staged changes directly from TUI.
- **Smart Commit**: Generate commit message suggestions based on staged changes.
- **System Status**: View vector database statistics and API status.
- **Index Progress**: Visualized file indexing process.
- **Token Monitor**: Monitor session token usage in the status bar at the bottom.
