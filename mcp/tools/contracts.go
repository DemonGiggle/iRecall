package tools

import (
	"time"

	"github.com/gigol/irecall/mcp/irecallapi"
)

type quoteResponse struct {
	ID               int64     `json:"id"`
	GlobalID         string    `json:"globalId"`
	AuthorUserID     string    `json:"authorUserId"`
	AuthorName       string    `json:"authorName"`
	SourceUserID     string    `json:"sourceUserId"`
	SourceName       string    `json:"sourceName"`
	SourceBackend    string    `json:"sourceBackend"`
	SourceNamespace  string    `json:"sourceNamespace"`
	SourceEntityType string    `json:"sourceEntityType"`
	SourceEntityID   string    `json:"sourceEntityId"`
	SourceLabel      string    `json:"sourceLabel"`
	SourceURL        string    `json:"sourceUrl"`
	Content          string    `json:"content"`
	Tags             []string  `json:"tags"`
	Version          int64     `json:"version"`
	IsOwnedByMe      bool      `json:"isOwnedByMe"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type listQuotesResponse struct {
	Limit  int             `json:"limit"`
	Offset int             `json:"offset"`
	Quotes []quoteResponse `json:"quotes"`
}

type recallResultResponse struct {
	Question string          `json:"question"`
	Keywords []string        `json:"keywords"`
	Quotes   []quoteResponse `json:"quotes"`
	Response string          `json:"response"`
}

type recallHistorySummaryResponse struct {
	ID        int64     `json:"id"`
	Question  string    `json:"question"`
	Response  string    `json:"response"`
	CreatedAt time.Time `json:"createdAt"`
}

type recallHistoryEntryResponse struct {
	ID        int64           `json:"id"`
	Question  string          `json:"question"`
	Response  string          `json:"response"`
	CreatedAt time.Time       `json:"createdAt"`
	Quotes    []quoteResponse `json:"quotes"`
}

func newListQuotesResponse(limit, offset int, quotes []irecallapi.Quote) listQuotesResponse {
	return listQuotesResponse{
		Limit:  limit,
		Offset: offset,
		Quotes: newQuoteResponses(quotes),
	}
}

func newQuoteResponse(quote *irecallapi.Quote) *quoteResponse {
	if quote == nil {
		return nil
	}
	return &quoteResponse{
		ID:               quote.ID,
		GlobalID:         quote.GlobalID,
		AuthorUserID:     quote.AuthorUserID,
		AuthorName:       quote.AuthorName,
		SourceUserID:     quote.SourceUserID,
		SourceName:       quote.SourceName,
		SourceBackend:    quote.SourceBackend,
		SourceNamespace:  quote.SourceNamespace,
		SourceEntityType: quote.SourceEntityType,
		SourceEntityID:   quote.SourceEntityID,
		SourceLabel:      quote.SourceLabel,
		SourceURL:        quote.SourceURL,
		Content:          quote.Content,
		Tags:             quote.Tags,
		Version:          quote.Version,
		IsOwnedByMe:      quote.IsOwnedByMe,
		CreatedAt:        quote.CreatedAt,
		UpdatedAt:        quote.UpdatedAt,
	}
}

func newQuoteResponses(quotes []irecallapi.Quote) []quoteResponse {
	if quotes == nil {
		return nil
	}
	out := make([]quoteResponse, 0, len(quotes))
	for _, quote := range quotes {
		out = append(out, *newQuoteResponse(&quote))
	}
	return out
}

func newRecallResultResponse(result *irecallapi.RecallResult) *recallResultResponse {
	if result == nil {
		return nil
	}
	return &recallResultResponse{
		Question: result.Question,
		Keywords: result.Keywords,
		Quotes:   newQuoteResponses(result.Quotes),
		Response: result.Response,
	}
}

func newRecallHistorySummaryResponses(history []irecallapi.RecallHistorySummary) []recallHistorySummaryResponse {
	if history == nil {
		return nil
	}
	out := make([]recallHistorySummaryResponse, 0, len(history))
	for _, entry := range history {
		out = append(out, recallHistorySummaryResponse{
			ID:        entry.ID,
			Question:  entry.Question,
			Response:  entry.Response,
			CreatedAt: entry.CreatedAt,
		})
	}
	return out
}

func newRecallHistoryEntryResponse(entry *irecallapi.RecallHistoryEntry) *recallHistoryEntryResponse {
	if entry == nil {
		return nil
	}
	return &recallHistoryEntryResponse{
		ID:        entry.ID,
		Question:  entry.Question,
		Response:  entry.Response,
		CreatedAt: entry.CreatedAt,
		Quotes:    newQuoteResponses(entry.Quotes),
	}
}
