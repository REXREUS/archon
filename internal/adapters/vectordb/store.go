
package vectordb

import (
	"context"
	"fmt"
	"runtime"

	"github.com/google/generative-ai-go/genai"
	"github.com/philippgille/chromem-go"
	"golang.org/x/time/rate"
	"google.golang.org/api/option"
)

type Store struct {
	db          *chromem.DB
	col         *chromem.Collection
	genaiClient *genai.Client
	limiter     *rate.Limiter
	embFunc     chromem.EmbeddingFunc
}

func NewStore(ctx context.Context, path string, apiKey string) (*Store, error) {
	db, err := chromem.NewPersistentDB(path, false)
	if err != nil {
		return nil, fmt.Errorf("failed to create persistent db: %w", err)
	}

	// Rate limiter: 1500 RPM (Requests Per Minute) -> 25 RPS
	limiter := rate.NewLimiter(rate.Limit(25), 1)

	var embFunc chromem.EmbeddingFunc
	var genaiClient *genai.Client
	if apiKey != "" {
		genaiClient, err = genai.NewClient(ctx, option.WithAPIKey(apiKey))
		if err != nil {
			return nil, fmt.Errorf("failed to create genai client: %w", err)
		}
		
		embFunc = func(ctx context.Context, text string) ([]float32, error) {
			// Wait for rate limiter
			if err := limiter.Wait(ctx); err != nil {
				return nil, err
			}

			model := genaiClient.EmbeddingModel("text-embedding-004")
			res, err := model.EmbedContent(ctx, genai.Text(text))
			if err != nil {
				return nil, err
			}
			return res.Embedding.Values, nil
		}
	}

	col, err := db.GetOrCreateCollection("codebase", nil, embFunc)
	if err != nil {
		return nil, fmt.Errorf("failed to get or create collection: %w", err)
	}

	return &Store{
		db:          db,
		col:         col,
		genaiClient: genaiClient,
		limiter:     limiter,
		embFunc:     embFunc,
	}, nil
}

func (s *Store) Clear(ctx context.Context) error {
	err := s.db.DeleteCollection("codebase")
	if err != nil {
		return fmt.Errorf("failed to delete collection: %w", err)
	}

	col, err := s.db.GetOrCreateCollection("codebase", nil, s.embFunc)
	if err != nil {
		return fmt.Errorf("failed to recreate collection: %w", err)
	}

	s.col = col
	return nil
}

func (s *Store) AddDocument(ctx context.Context, id string, content string, metadata map[string]string) error {
	doc := chromem.Document{
		ID:       id,
		Content:  content,
		Metadata: metadata,
	}

	return s.col.AddDocuments(ctx, []chromem.Document{doc}, runtime.NumCPU())
}

func (s *Store) Search(ctx context.Context, query string, n int) ([]chromem.Result, error) {
	return s.col.Query(ctx, query, n, nil, nil)
}

func (s *Store) Close() error {
	if s.genaiClient != nil {
		s.genaiClient.Close()
	}
	return nil
}
