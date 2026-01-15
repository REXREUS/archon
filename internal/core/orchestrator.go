package core

import (
	"archon/internal/adapters/parser"
	"archon/internal/adapters/vectordb"
	"archon/internal/utils"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

type Orchestrator struct {
	store  *vectordb.Store
	parser *parser.GenericParser
}

func NewOrchestrator(store *vectordb.Store) *Orchestrator {
	return &Orchestrator{
		store:  store,
		parser: parser.NewGenericParser(),
	}
}

func (o *Orchestrator) IndexDirectory(ctx context.Context, dir string, progress func(current, total int, file string)) error {
	var filesToIndex []string
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if utils.IsIgnored(path) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if info.IsDir() {
			return nil
		}
		lang := parser.DetectLanguage(path)
		if lang != parser.Unknown {
			filesToIndex = append(filesToIndex, path)
		}
		return nil
	})

	total := len(filesToIndex)
	for i, path := range filesToIndex {
		if progress != nil {
			progress(i+1, total, path)
		}

		err := o.IndexFile(ctx, path)
		if err != nil {
			continue
		}
	}
	return nil
}

func (o *Orchestrator) IndexFile(ctx context.Context, path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	symbols, err := o.parser.Parse(ctx, path, content)
	if err != nil {
		err = o.store.AddDocument(ctx, path, string(content), map[string]string{
			"file": path,
			"type": "file",
		})
		return err
	}

	for _, sym := range symbols {
		id := path + ":" + sym.Name
		metadata := map[string]string{
			"file": path,
			"name": sym.Name,
			"type": sym.Type,
		}
		err = o.store.AddDocument(ctx, id, sym.Code, metadata)
		if err != nil {
			continue
		}
	}
	return nil
}

func (o *Orchestrator) SearchContext(ctx context.Context, query string) (string, error) {
	results, err := o.store.Search(ctx, query, 5)
	if err != nil {
		return "", err
	}

	contextText := "Here are some relevant code snippets from the codebase:\n\n"
	for _, res := range results {
		contextText += "File: " + res.Metadata["file"] + "\n"
		if name, ok := res.Metadata["name"]; ok {
			contextText += "Symbol: " + name + " (" + res.Metadata["type"] + ")\n"
		}
		contextText += "```\n" + res.Content + "\n```\n\n"
	}

	return contextText, nil
}

func (o *Orchestrator) GetFilesForIndexing(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if utils.IsIgnored(path) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if info.IsDir() {
			return nil
		}
		lang := parser.DetectLanguage(path)
		if lang != parser.Unknown {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func (o *Orchestrator) WatchDirectory(ctx context.Context, dir string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	// Add directories to watch
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if utils.IsIgnored(path) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if info.IsDir() {
			return watcher.Add(path)
		}
		return nil
	})

	fmt.Println("Watching for changes...")
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			// We care about Write and Create events
			if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
				info, err := os.Stat(event.Name)
				if err != nil {
					continue
				}
				if info.IsDir() {
					// If it's a new directory, we should watch it too
					watcher.Add(event.Name)
					continue
				}

				lang := parser.DetectLanguage(event.Name)
				if lang != parser.Unknown {
					fmt.Printf("File changed: %s, re-indexing...\n", event.Name)
					o.IndexFile(ctx, event.Name)
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			fmt.Printf("Watcher error: %v\n", err)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
