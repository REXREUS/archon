# ArchonCLI Features 🌟

ArchonCLI provides a comprehensive set of features designed to make codebase navigation and understanding effortless.

## 🧠 Semantic Understanding
Unlike traditional search tools that treat code as plain text, ArchonCLI understands the structure of your code:
- **Syntax-Aware Parsing**: Uses Tree-Sitter to identify functions, classes, and variables.
- **Symbol Indexing**: Maps out how different parts of your code relate to each other.
- **Multi-language Support**: Supports Go, TypeScript, JavaScript, Python, and more.

## 🔍 Semantic Search (RAG)
ArchonCLI uses Retrieval-Augmented Generation to find the most relevant parts of your codebase:
- **Vector Embeddings**: Converts code snippets into mathematical vectors to find semantic similarities.
- **Local Storage**: Your code index is stored locally using `chromem-go`, ensuring privacy and speed.
- **Contextual Answers**: When you ask a question, ArchonCLI finds relevant code and provides it as context to the AI.

## 💻 Dual Interface
ArchonCLI offers two ways to interact with your code:

### 1. Interactive TUI (Terminal User Interface)
- **Chat Interface**: Familiar chat experience inside your terminal.
- **Thinking Indicator**: Real-time feedback when the AI is processing.
- **Context Viewer**: View exactly what code the AI is looking at.
- **Resource Monitor**: Track token usage and estimated costs in real-time.

### 2. Command Line Interface (CLI)
- **Scriptable**: Perfect for CI/CD pipelines or local automation.
- **JSON Output Support**: Easy integration with other tools.
- **Watch Mode**: Automatically re-indexes files as you save them.

## 🛠️ Developer Productivity Tools
- **AI Code Review**: Analyze staged changes for bugs and best practices.
- **Smart Commit Messages**: Automatically generate commit messages based on your changes.
- **Unit Test Generation**: Create comprehensive tests for your functions with a single command.
- **Architectural Analysis**: Detect code smells and design pattern violations.
- **Diagram Generation**: Generate Mermaid or PlantUML diagrams of your code structure.

## ⚡ Performance Optimization
- **Context Caching**: Uses Google Gemini 3's caching feature to reduce latency and costs by up to 90%.
- **Rate Limiting**: Intelligent token bucket implementation to stay within API quotas.
- **Incremental Indexing**: Only re-indexes files that have actually changed.
