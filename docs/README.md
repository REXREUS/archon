# ArchonCLI 🏛️

**ArchonCLI** is an autonomous AI-powered code architect assistant designed to help developers efficiently understand, analyze, and manage complex codebases.

Powered by **Google Gemini 3**, ArchonCLI goes beyond traditional code assistants by using a **Semantic Syntax-Aware RAG** approach that understands code structures deeply using **Tree-Sitter**.

## 🚀 Key Features

- **Semantic Analysis**: Understands code structure (functions, classes, dependencies) rather than just plain text.
- **Dual Interface**: Use the visually rich **TUI (Terminal User Interface)** interactive mode or the fast and automatable **CLI** mode.
- **RAG (Retrieval-Augmented Generation)**: Searches for relevant context across your repository to provide accurate answers.
- **Context Caching**: Optimizes token usage and speeds up responses on large repositories.
- **Security**: Secure API Key management and environment variable support.

## 📂 Documentation

To learn more about ArchonCLI, please refer to the following documents:

1.  [**Features**](FEATURES.md) - Detailed primary functionality.
2.  [**Usage Guide**](USAGE.md) - Installation and command reference.
3.  [**Configuration**](CONFIGURATION.md) - Setting up API Keys and application preferences.
4.  [**System Architecture**](ARCHITECTURE.md) - Behind-the-scenes technical details.

## 🛠️ System Requirements

- **Operating System**: Windows, macOS, or Linux.
- **Go**: Version 1.21 or higher (if building from source).
- **Google Gemini API Key**: Required to access the AI backend.

## 🏁 Quick Start

1.  **Installation**: Add the `archon` binary to your PATH.
2.  **Initialization**: Run `archon init` in your project folder.
3.  **Authentication**: Set your API Key with `archon auth --key "YOUR_API_KEY"`.
4.  **Index**: Run `archon index` to scan your code.
5.  **Ask**: Start asking with `archon ask "How does this application work?"` or simply run `archon` for TUI mode.

---
*Created with ❤️ by the Archon Development Team.*
