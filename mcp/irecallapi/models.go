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
