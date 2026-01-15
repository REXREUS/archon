package parser

import (
	"path/filepath"
	"regexp"
	"strings"
)

type Language string

const (
	Go         Language = "go"
	Python     Language = "python"
	TypeScript Language = "typescript"
	JavaScript Language = "javascript"
	Java       Language = "java"
	Rust       Language = "rust"
	Cpp        Language = "cpp"
	Php        Language = "php"
	Ruby       Language = "ruby"
	Csharp     Language = "csharp"
	Unknown    Language = "unknown"
)

func DetectLanguage(filename string) Language {
	ext := filepath.Ext(filename)
	switch ext {
	case ".go":
		return Go
	case ".py":
		return Python
	case ".ts", ".tsx":
		return TypeScript
	case ".js", ".jsx":
		return JavaScript
	case ".java":
		return Java
	case ".rs":
		return Rust
	case ".cpp", ".hpp", ".cc", ".cxx":
		return Cpp
	case ".php":
		return Php
	case ".rb":
		return Ruby
	case ".cs":
		return Csharp
	default:
		return Unknown
	}
}

func ExtractSymbols(lang Language, content []byte) ([]Symbol, error) {
	strContent := string(content)
	symbols := []Symbol{}

	type pattern struct {
		re   *regexp.Regexp
		name string
		symType string
	}

	patterns := []pattern{}

	switch lang {
	case Go:
		patterns = []pattern{
			{regexp.MustCompile(`(?m)^func\s+([A-Za-z0-9_]+)\s*\(`), "$1", "function"},
			{regexp.MustCompile(`(?m)^func\s+\([^\)]+\)\s+([A-Za-z0-9_]+)\s*\(`), "$1", "method"},
			{regexp.MustCompile(`(?m)^type\s+([A-Za-z0-9_]+)\s+(struct|interface)`), "$1", "type"},
		}
	case Python:
		patterns = []pattern{
			{regexp.MustCompile(`(?m)^def\s+([A-Za-z0-9_]+)\s*\(`), "$1", "function"},
			{regexp.MustCompile(`(?m)^class\s+([A-Za-z0-9_]+)(\s*\(|\s*:)`), "$1", "class"},
		}
	case TypeScript, JavaScript:
		patterns = []pattern{
			{regexp.MustCompile(`(?m)^(?:export\s+)?function\s+([A-Za-z0-9_]+)\s*\(`), "$1", "function"},
			{regexp.MustCompile(`(?m)^(?:export\s+)?class\s+([A-Za-z0-9_]+)`), "$1", "class"},
			{regexp.MustCompile(`(?m)^(?:export\s+)?interface\s+([A-Za-z0-9_]+)`), "$1", "interface"},
			{regexp.MustCompile(`(?m)^\s*([A-Za-z0-9_]+)\s*\([^\)]*\)\s*\{`), "$1", "method"},
		}
	case Java, Csharp:
		patterns = []pattern{
			{regexp.MustCompile(`(?m)(?:public|private|protected|static|\s)+class\s+([A-Za-z0-9_]+)`), "$1", "class"},
			{regexp.MustCompile(`(?m)(?:public|private|protected|static|\s)+interface\s+([A-Za-z0-9_]+)`), "$1", "interface"},
			{regexp.MustCompile(`(?m)(?:public|private|protected|static|\s)+[A-Za-z0-9_<>\[\]]+\s+([A-Za-z0-9_]+)\s*\([^\)]*\)\s*(?:\{|throws)`), "$1", "method"},
		}
	}

	// If no patterns or language not supported, return the entire file
	if len(patterns) == 0 {
		return []Symbol{{
			Name: "entire_file",
			Type: "file",
			Code: strContent,
		}}, nil
	}

	// Pre-split lines for all patterns
	lines := regexp.MustCompile(`\r?\n`).Split(strContent, -1)

	// Find all matches and extract blocks
	for _, p := range patterns {
		matches := p.re.FindAllStringSubmatchIndex(strContent, -1)
		for _, m := range matches {
			if len(m) < 4 {
				continue
			}
			name := strContent[m[2]:m[3]]
			
			// Find start line using string counting (robust for CRLF/LF)
			startLine := 0
			if m[0] > 0 {
				startLine = strings.Count(strContent[:m[0]], "\n")
			}

			// Capture a few lines after start or until next empty line/pattern
			endLine := startLine + 20 // Default chunk size
			if endLine > len(lines) {
				endLine = len(lines)
			}
			
			code := ""
			for i := startLine; i < endLine; i++ {
				code += lines[i] + "\n"
			}

			symbols = append(symbols, Symbol{
				Name: name,
				Type: p.symType,
				Code: code,
			})
		}
	}

	if len(symbols) == 0 {
		symbols = append(symbols, Symbol{
			Name: "entire_file",
			Type: "file",
			Code: strContent,
		})
	}

	return symbols, nil
}
