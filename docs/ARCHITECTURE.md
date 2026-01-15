# ArchonCLI System Architecture 🏛️

ArchonCLI is built with a focus on high performance, portability, and deep semantic understanding of source code.

## 🏗️ Technical Stack

- **Language**: Go (Golang) - Selected for its static binary compilation, high-performance concurrency (Goroutines), and cross-platform compatibility.
- **LLM Backend**: Google Gemini 3 (Pro & Flash) - Used for deep reasoning and fast processing.
- **Vector Database**: `chromem-go` - A pure Go, embedded vector database that allows ArchonCLI to remain a single binary without external C++ dependencies.
- **Code Parser**: Tree-Sitter (`go-tree-sitter`) - Provides incremental, syntax-aware parsing for multiple programming languages.
- **CLI Framework**: Cobra - Handles command management, flags, and automated documentation.
- **TUI Framework**: Bubble Tea - Implements "The Elm Architecture" for a visually rich and interactive terminal interface.

## 🧩 Core Components

### 1. Orchestrator (`internal/core`)
The brain of the application. It coordinates the flow between the user interface, the code parser, the vector database, and the Gemini API.

### 2. Semantic Parser (`internal/adapters/parser`)
Instead of simple line-based chunking, ArchonCLI uses Tree-Sitter to parse code into an Abstract Syntax Tree (AST). It extracts:
- Function definitions
- Class/Struct declarations
- Documentation comments (docstrings)
- Interface definitions

### 3. Vector Store (`internal/adapters/vectordb`)
Manages the local embedding storage. It uses Google's `text-embedding-004` model to vectorize code symbols and stores them in `chromem-go` for fast retrieval.

### 4. Gemini Client (`internal/adapters/gemini`)
A wrapper around the official Google Generative AI Go SDK. It implements:
- **Rate Limiting**: A token bucket algorithm to stay within API quotas.
- **Context Caching**: Management of server-side state to reduce token consumption.
- **Deep Think**: Integration with Gemini 3's reasoning capabilities.

## 🔄 Workflow: The RAG Pipeline

1. **Indexing Phase**:
   - Files are scanned and filtered (respecting `.gitignore`).
   - Tree-Sitter parses relevant files into symbols.
   - Each symbol is embedded and stored in the local vector DB.
   
2. **Query Phase**:
   - The user asks a question via CLI or TUI.
   - The question is embedded and a similarity search is performed against the vector DB.
   - Relevant code snippets (context) are retrieved.
   
3. **Reasoning Phase**:
   - The context and user query are sent to Gemini 3.
   - If a context cache is available and matching, it is used to speed up processing.
   - Gemini generates a detailed answer based on the provided code context.

## 📂 Project Structure

```
archon/
├── cmd/archon/          # Application entry point
├── internal/
│   ├── core/            # Business logic & Orchestration
│   ├── adapters/        # External system integrations (Gemini, VectorDB, Parser)
│   ├── ui/              # Presentation layer (CLI, TUI, LSP)
│   ├── config/          # Configuration management
│   └── utils/           # Shared utility functions
├── scripts/             # Build and deployment scripts
└── docs/                # Project documentation
```
