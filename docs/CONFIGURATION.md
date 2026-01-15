# ArchonCLI Configuration ⚙️

ArchonCLI uses a combination of configuration files, environment variables, and command-line arguments for maximum flexibility.

## 📄 The `.archon.yaml` File

This file is created when you run `archon init`. It can be stored in your project folder or your home directory (`~/.archon.yaml`).

Example file content:
```yaml
gemini_key: "AIzaSy..."
model_id: "gemini-3-pro-preview"
project_hash: "a1b2c3d4..."
cache_name: "cached_contents/..."
```

### Configuration Parameters:
- `gemini_key`: Your Google Gemini API Key.
- `model_id`: The ID of the model being used (Default: `gemini-3-pro-preview`).
- `project_hash`: The last hash of your project for caching purposes.
- `cache_name`: The ID of the active context cache on Google's servers.

## 🌍 Environment Variables

Environment variables have higher priority than the configuration file.

- `ARCHON_GEMINI_KEY`: Sets the API Key.
- `ARCHON_MODEL_ID`: Sets the AI model to be used.

Example:
```bash
export ARCHON_GEMINI_KEY="AIzaSy_YOUR_KEY"
archon ask "..."
```

## 🧠 Context Caching (Gemini 3)

ArchonCLI utilizes the **Context Caching** feature from Google Gemini 3 to improve response speed and reduce token costs on large codebases.

### How It Works:
1. **Project Hashing**: ArchonCLI calculates a unique hash based on your project's file contents (ignoring folders like `.git`, `node_modules`, `build`, etc).
2. **Upload & Cache**: The first time you ask a question (`ask`), ArchonCLI sends the entire code context to Google Gemini and requests it to be stored as a cache.
3. **ID Utilization**: Google returns a `cache_name`. ArchonCLI saves this ID along with the project hash into the `.archon.yaml` file.
4. **Fast Responses**: On subsequent questions, ArchonCLI only sends that cache ID. Google doesn't need to re-process thousands of lines of code, making responses much faster.

### Cache Activation Requirements:
- **Token Threshold**: Google Gemini will only create a cache if the content sent is **at least 32,768 tokens** (approximately 130,000 characters of code). If your project is too small, this feature will remain "Inactive".
- **Matching Hash**: If you change the code, the project hash changes. Context Cache will become "Inactive (Hash Mismatch)" until you ask a new question, at which point ArchonCLI will automatically update the cache.
- **TTL (Time To Live)**: Caches on Google's servers are temporary (usually lasting several hours/days depending on usage). ArchonCLI will auto-sync if the cache has been deleted from the server.

---

## 🔐 API Key Security

ArchonCLI prioritizes your data security:
1. **Local Storage**: API Keys are stored in the `.archon.yaml` file with restricted access permissions (0600 on Unix systems).
2. **No Training Data Usage**: If using Enterprise mode (Vertex AI), your data will not be used to train Google's public models.
3. **File Filtering**: You can use a `.gitignore` file or specific configuration to prevent ArchonCLI from reading sensitive files like `.env` or `secrets.go`.

## 🔄 Updating Configuration

You can change the configuration at any time via CLI:
```bash
archon config set model "gemini-3-flash"
archon config list
```
