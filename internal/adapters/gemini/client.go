package gemini

import (
	"context"
	"fmt"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type Client struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

func NewClient(ctx context.Context, apiKey string, modelID string) (*Client, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}

	model := client.GenerativeModel(modelID)

	return &Client{
		client: client,
		model:  model,
	}, nil
}

func (c *Client) SetCachedContent(name string) {
	c.model.CachedContentName = name
}

func (c *Client) Client() *genai.Client {
	return c.client
}

type Response struct {
	Text         string
	PromptTokens int
	AnswerTokens int
	TotalTokens  int
}

func (c *Client) Ask(ctx context.Context, prompt string) (*Response, error) {
	resp, err := c.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	if len(resp.Candidates) == 0 {
		return nil, fmt.Errorf("no candidates in response")
	}

	// Simple extraction of text from the first candidate
	var result string
	for _, part := range resp.Candidates[0].Content.Parts {
		if text, ok := part.(genai.Text); ok {
			result += string(text)
		}
	}

	return &Response{
		Text:         result,
		PromptTokens: int(resp.UsageMetadata.PromptTokenCount),
		AnswerTokens: int(resp.UsageMetadata.CandidatesTokenCount),
		TotalTokens:  int(resp.UsageMetadata.TotalTokenCount),
	}, nil
}

func (c *Client) Close() {
	c.client.Close()
}
