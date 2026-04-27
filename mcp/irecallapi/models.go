package irecallapi

import (
	appbackend "github.com/gigol/irecall/app"
	"github.com/gigol/irecall/core"
)

type BootstrapState = appbackend.BootstrapState
type RecallResult = appbackend.RecallResult
type Quote = core.Quote
type RecallHistorySummary = core.RecallHistorySummary
type RecallHistoryEntry = core.RecallHistoryEntry

type AddQuoteRequest struct {
	Content string `json:"content"`
}

type RunRecallRequest struct {
	Question string `json:"question"`
}

type SaveRecallAsQuoteRequest struct {
	Question string   `json:"question"`
	Response string   `json:"response"`
	Keywords []string `json:"keywords"`
}

type UpdateQuoteRequest struct {
	ID      int64  `json:"id"`
	Content string `json:"content"`
}

type DeleteQuotesRequest struct {
	IDs []int64 `json:"ids"`
}

type DeleteRecallHistoryRequest struct {
	IDs []int64 `json:"ids"`
}

type OKResponse struct {
	OK bool `json:"ok"`
}
