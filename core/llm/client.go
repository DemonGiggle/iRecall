package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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

// Chat sends a chat completion request.
// If tokenCh is non-nil, streaming is enabled; tokens are sent to the channel
// which is closed when the stream ends. Returns empty string when streaming.
func (c *Client) Chat(ctx context.Context, msgs []Message, tokenCh chan<- string) (string, error) {
	stream := tokenCh != nil

	body := map[string]any{
		"model":    c.cfg.Model,
		"messages": msgs,
		"stream":   stream,
	}
	data, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.cfg.BaseURL()+"/chat/completions", bytes.NewReader(data))
	if err != nil {
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

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("chat request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("provider returned %d: %s", resp.StatusCode, string(b))
	}

	if stream {
		go func() {
			defer close(tokenCh)
			parseSSE(resp.Body, tokenCh)
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
		return "", fmt.Errorf("decode response: %w", err)
	}
	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}
	return result.Choices[0].Message.Content, nil
}

// FetchModels calls GET /v1/models and returns sorted model IDs.
func (c *Client) FetchModels(ctx context.Context) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.cfg.BaseURL()+"/models", nil)
	if err != nil {
		return nil, err
	}
	if c.cfg.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.cfg.APIKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch models: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("provider returned %d: %s", resp.StatusCode, string(b))
	}

	var result struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode models: %w", err)
	}

	ids := make([]string, 0, len(result.Data))
	for _, m := range result.Data {
		ids = append(ids, m.ID)
	}
	sort.Strings(ids)
	return ids, nil
}

// parseSSE reads a streaming SSE response and sends each content token to ch.
func parseSSE(r io.Reader, ch chan<- string) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		payload := strings.TrimPrefix(line, "data: ")
		if payload == "[DONE]" {
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
			continue
		}
		if len(chunk.Choices) > 0 {
			if tok := chunk.Choices[0].Delta.Content; tok != "" {
				ch <- tok
			}
		}
	}
}
