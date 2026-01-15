package gemini

import (
	"archon/internal/utils"
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
)

type CacheManager struct {
	client *genai.Client
}

func NewCacheManager(client *genai.Client) *CacheManager {
	return &CacheManager{client: client}
}

func (cm *CacheManager) CreateContextCache(ctx context.Context, modelID string, files []string) (string, error) {
	var contents []*genai.Content
	
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		
		contents = append(contents, &genai.Content{
			Role: "user",
			Parts: []genai.Part{
				genai.Text(fmt.Sprintf("File: %s\n\n%s", file, string(data))),
			},
		})
	}

	cachedContent := &genai.CachedContent{
		Model:      modelID,
		Contents:   contents,
	}

	res, err := cm.client.CreateCachedContent(ctx, cachedContent)
	if err != nil {
		return "", fmt.Errorf("failed to create cached content: %w (Note: Context caching requires at least 32,768 tokens)", err)
	}

	return res.Name, nil
}

func (cm *CacheManager) ListCaches(ctx context.Context) ([]*genai.CachedContent, error) {
	var caches []*genai.CachedContent
	it := cm.client.ListCachedContents(ctx)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		caches = append(caches, resp)
	}
	return caches, nil
}

func (cm *CacheManager) DeleteCache(ctx context.Context, name string) error {
	return cm.client.DeleteCachedContent(ctx, name)
}

func CalculateProjectHash(dir string) (string, error) {
	hash := sha256.New()
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Ignore irrelevant files (node_modules, .git, etc.)
		if utils.IsIgnored(path) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if info.IsDir() {
			return nil
		}
		
		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		hash.Write(data)
		return nil
	})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
