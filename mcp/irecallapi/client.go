package irecallapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Config struct {
	BaseURL     string
	APIToken    string
	HTTPTimeout time.Duration
}

type Client struct {
	baseURL    string
	apiToken   string
	httpClient *http.Client
}

const maxResponseBodySize = 2 << 20

var errResponseTooLarge = errors.New("response exceeds 2MB limit")

func NewClient(cfg Config) (*Client, error) {
	if strings.TrimSpace(cfg.BaseURL) == "" {
		return nil, errors.New("base URL is required")
	}
	if strings.TrimSpace(cfg.APIToken) == "" {
		return nil, errors.New("API token is required")
	}
	if cfg.HTTPTimeout <= 0 {
		cfg.HTTPTimeout = 15 * time.Second
	}
	if _, err := url.Parse(cfg.BaseURL); err != nil {
		return nil, fmt.Errorf("parse base URL: %w", err)
	}
	return &Client{
		baseURL:  strings.TrimRight(cfg.BaseURL, "/"),
		apiToken: cfg.APIToken,
		httpClient: &http.Client{
			Timeout: cfg.HTTPTimeout,
		},
	}, nil
}

func (c *Client) BootstrapState(ctx context.Context) (*BootstrapState, error) {
	var value BootstrapState
	if err := c.doJSON(ctx, http.MethodGet, "/api/app/bootstrap-state", nil, &value); err != nil {
		return nil, err
	}
	return &value, nil
}

func (c *Client) ListQuotes(ctx context.Context) ([]Quote, error) {
	var value []Quote
	if err := c.doJSON(ctx, http.MethodGet, "/api/app/list-quotes", nil, &value); err != nil {
		return nil, err
	}
	return value, nil
}

func (c *Client) AddQuote(ctx context.Context, content string) (*Quote, error) {
	var value Quote
	if err := c.doJSON(ctx, http.MethodPost, "/api/app/add-quote", AddQuoteRequest{Content: content}, &value); err != nil {
		return nil, err
	}
	return &value, nil
}

func (c *Client) RunRecall(ctx context.Context, question string) (*RecallResult, error) {
	var value RecallResult
	if err := c.doJSON(ctx, http.MethodPost, "/api/app/run-recall", RunRecallRequest{Question: question}, &value); err != nil {
		return nil, err
	}
	return &value, nil
}

func (c *Client) SaveRecallAsQuote(ctx context.Context, question, response string, keywords []string) (*Quote, error) {
	var value Quote
	if err := c.doJSON(ctx, http.MethodPost, "/api/app/save-recall-as-quote", SaveRecallAsQuoteRequest{
		Question: question,
		Response: response,
		Keywords: keywords,
	}, &value); err != nil {
		return nil, err
	}
	return &value, nil
}

func (c *Client) doJSON(ctx context.Context, method, path string, payload any, dst any) error {
	var body io.Reader
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("marshal request body: %w", err)
		}
		body = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	applyAuth(req, c.apiToken)
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("call iRecall API: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBodySize+1))
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}
	if len(respBody) > maxResponseBodySize {
		return errResponseTooLarge
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return parseAPIError(resp.StatusCode, respBody)
	}
	if dst == nil || len(bytes.TrimSpace(respBody)) == 0 {
		return nil
	}
	if err := json.Unmarshal(respBody, dst); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}
