package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// IsIgnored checking if path should be ignored by Archon (indexing, hashing, etc)
func IsIgnored(path string) bool {
	name := filepath.Base(path)
	
	// Folders to ignore
	ignoredDirs := []string{".git", "node_modules", "vendor", "chromem_db", "bin", "build", "obj", ".idea", ".vscode"}
	
	// Files to ignore specifically
	ignoredFiles := []string{".archon.yaml", "archon.exe"}
	
	// Extensions to ignore (binaries, logs, etc)
	ignoredExts := []string{".exe", ".dll", ".so", ".dylib", ".bin", ".log", ".test"}

	// Check specific files
	for _, f := range ignoredFiles {
		if name == f {
			return true
		}
	}

	// Check if path contains any ignored directory
	parts := strings.Split(filepath.ToSlash(path), "/")
	for _, part := range parts {
		for _, dir := range ignoredDirs {
			if part == dir {
				return true
			}
		}
	}

	// Check extension
	ext := filepath.Ext(name)
	for _, e := range ignoredExts {
		if ext == e {
			return true
		}
	}

	// Ignore common temporary files
	if strings.HasPrefix(name, ".") && name != ".archon.yaml" && !strings.Contains(path, "github") {
		// allow hidden files that are not common config folders
		// but keep it simple for now
	}

	return false
}

// IsDir checks if path is a directory
func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
