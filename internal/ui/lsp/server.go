package lsp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
)

// Request mewakili pesan JSON-RPC 2.0 Request
type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// Response mewakili pesan JSON-RPC 2.0 Response
type Response struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Result  interface{}     `json:"result,omitempty"`
	Error   *ResponseError  `json:"error,omitempty"`
}

type ResponseError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Notification mewakili pesan JSON-RPC 2.0 Notification
type Notification struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type Server struct {
	reader *bufio.Reader
	writer io.Writer
	mu     sync.Mutex
}

func NewServer() *Server {
	return &Server{
		reader: bufio.NewReader(os.Stdin),
		writer: os.Stdout,
	}
}

func (s *Server) Start() error {
	for {
		contentLength, err := s.readHeader()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		body := make([]byte, contentLength)
		_, err = io.ReadFull(s.reader, body)
		if err != nil {
			return err
		}

		var req Request
		if err := json.Unmarshal(body, &req); err != nil {
			continue
		}

		go s.handleMessage(req)
	}
}

func (s *Server) readHeader() (int, error) {
	var contentLength int
	for {
		line, err := s.reader.ReadString('\n')
		if err != nil {
			return 0, err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		if strings.HasPrefix(line, "Content-Length:") {
			parts := strings.Split(line, ":")
			if len(parts) == 2 {
				val, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
				contentLength = val
			}
		}
	}
	return contentLength, nil
}

func (s *Server) sendResponse(id interface{}, result interface{}) {
	res := Response{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
	s.write(res)
}

func (s *Server) sendError(id interface{}, code int, message string) {
	res := Response{
		JSONRPC: "2.0",
		ID:      id,
		Error: &ResponseError{
			Code:    code,
			Message: message,
		},
	}
	s.write(res)
}

func (s *Server) write(msg interface{}) {
	data, _ := json.Marshal(msg)
	s.mu.Lock()
	defer s.mu.Unlock()
	fmt.Fprintf(s.writer, "Content-Length: %d\r\n\r\n%s", len(data), data)
}
