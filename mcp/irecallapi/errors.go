package irecallapi

import (
	"encoding/json"
	"fmt"
	"strings"
)

type APIError struct {
	StatusCode int
	Message    string
	Body       string
}

func (e *APIError) Error() string {
	message := strings.TrimSpace(e.Message)
	if message == "" {
		message = "unexpected API error"
	}
	if e.StatusCode <= 0 {
		return message
	}
	return fmt.Sprintf("iRecall API returned %d: %s", e.StatusCode, message)
}

func parseAPIError(statusCode int, body []byte) error {
	var payload struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(body, &payload); err == nil && strings.TrimSpace(payload.Error) != "" {
		return &APIError{
			StatusCode: statusCode,
			Message:    strings.TrimSpace(payload.Error),
			Body:       strings.TrimSpace(string(body)),
		}
	}
	return &APIError{
		StatusCode: statusCode,
		Message:    strings.TrimSpace(string(body)),
		Body:       strings.TrimSpace(string(body)),
	}
}
