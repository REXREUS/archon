package lsp

import (
	"archon/internal/adapters/gemini"
	"archon/internal/adapters/vectordb"
	"archon/internal/config"
	"archon/internal/core"
	"context"
	"encoding/json"
	"fmt"
	"os"
)

type ExecuteCommandParams struct {
	Command   string            `json:"command"`
	Arguments []json.RawMessage `json:"arguments"`
}

func (s *Server) handleMessage(req Request) {
	switch req.Method {
	case "initialize":
		s.handleInitialize(req)
	case "initialized":
		// Notification, no response
	case "shutdown":
		s.sendResponse(req.ID, nil)
	case "exit":
		os.Exit(0)
	case "workspace/executeCommand":
		s.handleExecuteCommand(req)
	default:
		if req.ID != nil {
			s.sendError(req.ID, -32601, fmt.Sprintf("Method '%s' not found", req.Method))
		}
	}
}

func (s *Server) handleInitialize(req Request) {
	result := map[string]interface{}{
		"capabilities": map[string]interface{}{
			"executeCommandProvider": map[string]interface{}{
				"commands": []string{
					"archon/ask",
					"archon/index",
					"archon/explain",
				},
			},
		},
		"serverInfo": map[string]interface{}{
			"name":    "archon-lsp",
			"version": "0.1.0",
		},
	}
	s.sendResponse(req.ID, result)
}

func (s *Server) handleExecuteCommand(req Request) {
	var params ExecuteCommandParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		s.sendError(req.ID, -32602, "Invalid params")
		return
	}

	ctx := context.Background()
	cfg, _ := config.LoadConfig()

	switch params.Command {
	case "archon/ask":
		if len(params.Arguments) == 0 {
			s.sendError(req.ID, -32602, "Query argument missing")
			return
		}
		var query string
		json.Unmarshal(params.Arguments[0], &query)
		
		response, err := s.executeAsk(ctx, cfg, query)
		if err != nil {
			s.sendError(req.ID, -32000, err.Error())
		} else {
			s.sendResponse(req.ID, response)
		}

	case "archon/index":
		err := s.executeIndex(ctx, cfg)
		if err != nil {
			s.sendError(req.ID, -32000, err.Error())
		} else {
			s.sendResponse(req.ID, "Indexing complete")
		}

	case "archon/explain":
		if len(params.Arguments) == 0 {
			s.sendError(req.ID, -32602, "Target argument missing")
			return
		}
		var target string
		json.Unmarshal(params.Arguments[0], &target)
		
		response, err := s.executeExplain(ctx, cfg, target)
		if err != nil {
			s.sendError(req.ID, -32000, err.Error())
		} else {
			s.sendResponse(req.ID, response)
		}

	default:
		s.sendError(req.ID, -32601, fmt.Sprintf("Command '%s' not found", params.Command))
	}
}

func (s *Server) executeAsk(ctx context.Context, cfg *config.Config, query string) (string, error) {
	if cfg.GeminiKey == "" {
		return "", fmt.Errorf("Gemini API key not found")
	}

	store, err := vectordb.NewStore(ctx, "./chromem_db", cfg.GeminiKey)
	var contextText string
	if err == nil {
		defer store.Close()
		orchestrator := core.NewOrchestrator(store)
		contextText, _ = orchestrator.SearchContext(ctx, query)
	}

	client, err := gemini.NewClient(ctx, cfg.GeminiKey, cfg.ModelID)
	if err != nil {
		return "", err
	}
	defer client.Close()

	prompt := query
	if contextText != "" {
		prompt = fmt.Sprintf("%s\n\nQuestion: %s", contextText, query)
	}

	resp, err := client.Ask(ctx, prompt)
	if err != nil {
		return "", err
	}

	return resp.Text, nil
}

func (s *Server) executeIndex(ctx context.Context, cfg *config.Config) error {
	if cfg.GeminiKey == "" {
		return fmt.Errorf("Gemini API key not found")
	}

	store, err := vectordb.NewStore(ctx, "./chromem_db", cfg.GeminiKey)
	if err != nil {
		return err
	}
	defer store.Close()

	orchestrator := core.NewOrchestrator(store)
	return orchestrator.IndexDirectory(ctx, ".", nil)
}

func (s *Server) executeExplain(ctx context.Context, cfg *config.Config, target string) (string, error) {
	if cfg.GeminiKey == "" {
		return "", fmt.Errorf("Gemini API key not found")
	}

	store, err := vectordb.NewStore(ctx, "./chromem_db", cfg.GeminiKey)
	var contextText string
	if err == nil {
		defer store.Close()
		orchestrator := core.NewOrchestrator(store)
		contextText, _ = orchestrator.SearchContext(ctx, "Explain "+target)
	}

	client, err := gemini.NewClient(ctx, cfg.GeminiKey, cfg.ModelID)
	if err != nil {
		return "", err
	}
	defer client.Close()

	var prompt string
	if contextText != "" {
		prompt = fmt.Sprintf("%s\n\nTask: Explain in detail about: %s", contextText, target)
	} else {
		content, err := os.ReadFile(target)
		if err == nil {
			prompt = fmt.Sprintf("Explain the following code:\n```\n%s\n```", string(content))
		} else {
			prompt = fmt.Sprintf("Explain about: %s", target)
		}
	}

	resp, err := client.Ask(ctx, prompt)
	if err != nil {
		return "", err
	}

	return resp.Text, nil
}
