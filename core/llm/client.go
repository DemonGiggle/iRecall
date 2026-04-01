package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sort"
	"strings"
	"time"
)

// Client is an HTTP client for any OpenAI-compatible chat completions API.
type Client struct {
	cfg        ProviderConfig
	httpClient *http.Client
}

func NewClient(cfg ProviderConfig) *Client {
	slog.Debug("llm: creating client", "base_url", cfg.BaseURL(), "model", cfg.Model, "has_api_key", cfg.APIKey != "")
	return &Client{
		cfg:        cfg,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// Message is a single chat turn.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatOptions controls optional per-request parameters.
type ChatOptions struct {
	Temperature *float64 // nil = provider default
	MaxTokens   *int     // nil = provider default
}

// Chat sends a chat completion request.
// If tokenCh is non-nil, streaming is enabled; tokens are sent to the channel
// which is closed when the stream ends. Returns empty string when streaming.
func (c *Client) Chat(ctx context.Context, msgs []Message, tokenCh chan<- string, opts ...ChatOptions) (string, error) {
	stream := tokenCh != nil
	url := c.cfg.BaseURL() + "/chat/completions"

	slog.Info("llm: chat request", "url", url, "model", c.cfg.Model, "stream", stream, "msg_count", len(msgs))
	for i, m := range msgs {
		slog.Debug("llm: chat message", "index", i, "role", m.Role, "content_len", len(m.Content))
	}

	body := map[string]any{
		"model":    c.cfg.Model,
		"messages": msgs,
		"stream":   stream,
	}
	if len(opts) > 0 {
		o := opts[0]
		if o.Temperature != nil {
			body["temperature"] = *o.Temperature
		}
		if o.MaxTokens != nil {
			body["max_tokens"] = *o.MaxTokens
		}
	}
	data, err := json.Marshal(body)
	if err != nil {
		slog.Error("llm: failed to marshal request body", "error", err)
		return "", err
	}
	slog.Debug("llm: request body", "size_bytes", len(data))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		slog.Error("llm: failed to create request", "error", err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.cfg.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.cfg.APIKey)
	}

	// For streaming, remove the client-level timeout so the stream can be long.
	client := c.httpClient
	if stream {
		client = &http.Client{}
	}

	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("llm: request failed", "error", err, "elapsed", time.Since(start))
		return "", fmt.Errorf("chat request: %w", err)
	}
	defer resp.Body.Close()

	slog.Info("llm: response received", "status", resp.StatusCode, "elapsed", time.Since(start))

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		slog.Error("llm: non-OK response", "status", resp.StatusCode, "body", string(b))
		return "", fmt.Errorf("provider returned %d: %s", resp.StatusCode, string(b))
	}

	if stream {
		slog.Debug("llm: starting SSE stream parse")
		go func() {
			defer close(tokenCh)
			parseSSE(resp.Body, tokenCh)
			slog.Debug("llm: SSE stream finished")
		}()
		return "", nil
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		slog.Error("llm: failed to decode response", "error", err)
		return "", fmt.Errorf("decode response: %w", err)
	}
	if len(result.Choices) == 0 {
		slog.Error("llm: no choices in response")
		return "", fmt.Errorf("no choices in response")
	}
	content := result.Choices[0].Message.Content
	slog.Info("llm: chat completed", "response_len", len(content))
	slog.Debug("llm: chat response content", "content", content)
	return content, nil
}

// FetchModels calls GET /v1/models and returns sorted model IDs.
func (c *Client) FetchModels(ctx context.Context) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	url := c.cfg.BaseURL() + "/models"
	slog.Info("llm: fetching models", "url", url)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	if c.cfg.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.cfg.APIKey)
	}

	start := time.Now()
	resp, err := c.httpClient.Do(req)
	if err != nil {
		slog.Error("llm: fetch models failed", "error", err, "elapsed", time.Since(start))
		return nil, fmt.Errorf("fetch models: %w", err)
	}
	defer resp.Body.Close()

	slog.Debug("llm: models response", "status", resp.StatusCode, "elapsed", time.Since(start))

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		slog.Error("llm: fetch models non-OK", "status", resp.StatusCode, "body", string(b))
		return nil, fmt.Errorf("provider returned %d: %s", resp.StatusCode, string(b))
	}

	var result struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		slog.Error("llm: decode models failed", "error", err)
		return nil, fmt.Errorf("decode models: %w", err)
	}

	ids := make([]string, 0, len(result.Data))
	for _, m := range result.Data {
		ids = append(ids, m.ID)
	}
	sort.Strings(ids)
	slog.Info("llm: fetched models", "count", len(ids), "models", ids)
	return ids, nil
}

// parseSSE reads a streaming SSE response and sends each content token to ch.
func parseSSE(r io.Reader, ch chan<- string) {
	scanner := bufio.NewScanner(r)
	tokenCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		payload := strings.TrimPrefix(line, "data: ")
		if payload == "[DONE]" {
			slog.Debug("llm: SSE received [DONE]", "total_tokens", tokenCount)
			return
		}
		var chunk struct {
			Choices []struct {
				Delta struct {
					Content string `json:"content"`
				} `json:"delta"`
			} `json:"choices"`
		}
		if err := json.Unmarshal([]byte(payload), &chunk); err != nil {
			slog.Debug("llm: SSE unmarshal error", "error", err, "payload", payload)
			continue
		}
		if len(chunk.Choices) > 0 {
			if tok := chunk.Choices[0].Delta.Content; tok != "" {
				tokenCount++
				ch <- tok
			}
		}
	}
	if err := scanner.Err(); err != nil {
		slog.Error("llm: SSE scanner error", "error", err)
	}
	slog.Debug("llm: SSE stream ended", "total_tokens", tokenCount)
}
