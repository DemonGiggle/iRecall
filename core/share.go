package core

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gigol/irecall/core/db"
)

func (e *Engine) ExportQuotes(ctx context.Context, ids []int64) ([]byte, error) {
	_ = ctx
	if len(ids) == 0 {
		return nil, fmt.Errorf("no quotes selected for export")
	}

	entries := make([]SharedQuoteEntry, 0, len(ids))
	for _, id := range ids {
		q, err := e.loadQuote(id)
		if err != nil {
			return nil, err
		}
		entries = append(entries, SharedQuoteEntry{
			GlobalID:         q.GlobalID,
			AuthorUserID:     q.AuthorUserID,
			AuthorName:       q.AuthorName,
			SourceUserID:     q.SourceUserID,
			SourceName:       q.SourceName,
			SourceBackend:    q.SourceBackend,
			SourceNamespace:  q.SourceNamespace,
			SourceEntityType: q.SourceEntityType,
			SourceEntityID:   q.SourceEntityID,
			SourceLabel:      q.SourceLabel,
			SourceURL:        q.SourceURL,
			Version:          q.Version,
			Content:          q.Content,
			Tags:             append([]string(nil), q.Tags...),
			CreatedAtUTC:     q.CreatedAt.UTC(),
			UpdatedAtUTC:     q.UpdatedAt.UTC(),
		})
	}

	payload, err := json.MarshalIndent(SharedQuoteEnvelope{
		SchemaVersion: ShareSchemaVersion,
		ExportedAt:    time.Now().UTC(),
		Quotes:        entries,
	}, "", "  ")
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func (e *Engine) ImportSharedQuotes(ctx context.Context, payload []byte) (ImportResult, error) {
	_ = ctx

	var env SharedQuoteEnvelope
	if err := json.Unmarshal(payload, &env); err != nil {
		return ImportResult{}, fmt.Errorf("decode share payload: %w", err)
	}
	if env.SchemaVersion != 1 && env.SchemaVersion != ShareSchemaVersion {
		return ImportResult{}, fmt.Errorf("unsupported share schema version: %d", env.SchemaVersion)
	}

	var result ImportResult
	for _, entry := range env.Quotes {
		if err := validateSharedQuoteEntry(entry); err != nil {
			return ImportResult{}, err
		}

		existing, lookupErr := e.store.GetQuoteByGlobalID(entry.GlobalID)
		if lookupErr != nil && lookupErr != sql.ErrNoRows {
			return ImportResult{}, lookupErr
		}

		identity := db.QuoteIdentity{
			GlobalID:         entry.GlobalID,
			AuthorUserID:     entry.AuthorUserID,
			AuthorName:       entry.AuthorName,
			SourceUserID:     entry.SourceUserID,
			SourceName:       entry.SourceName,
			SourceBackend:    entry.SourceBackend,
			SourceNamespace:  entry.SourceNamespace,
			SourceEntityType: entry.SourceEntityType,
			SourceEntityID:   entry.SourceEntityID,
			SourceLabel:      entry.SourceLabel,
			SourceURL:        entry.SourceURL,
			Version:          entry.Version,
		}
		identity = normalizeImportedQuoteIdentity(env.SchemaVersion, identity)

		tagIDs, err := e.store.UpsertTags(entry.Tags)
		if err != nil {
			return ImportResult{}, err
		}

		switch {
		case lookupErr == sql.ErrNoRows:
			id, err := e.store.InsertImportedQuote(entry.Content, identity, entry.CreatedAtUTC.Unix(), entry.UpdatedAtUTC.Unix())
			if err != nil {
				return ImportResult{}, err
			}
			if err := e.store.ReplaceQuoteTags(id, tagIDs); err != nil {
				return ImportResult{}, err
			}
			if err := e.store.UpdateQuoteFTS(id, entry.Tags); err != nil {
				return ImportResult{}, err
			}
			result.Inserted++
		case entry.Version > existing.Version:
			if err := e.store.UpdateImportedQuote(existing.ID, entry.Content, identity, entry.CreatedAtUTC.Unix(), entry.UpdatedAtUTC.Unix()); err != nil {
				return ImportResult{}, err
			}
			if err := e.store.ReplaceQuoteTags(existing.ID, tagIDs); err != nil {
				return ImportResult{}, err
			}
			if err := e.store.UpdateQuoteFTS(existing.ID, entry.Tags); err != nil {
				return ImportResult{}, err
			}
			result.Updated++
		case entry.Version == existing.Version:
			if err := e.store.UpdateImportedQuote(existing.ID, existing.Content, identity, existing.CreatedAt, existing.UpdatedAt); err != nil {
				return ImportResult{}, err
			}
			result.Duplicates++
		default:
			result.Stale++
		}
	}

	return result, nil
}

func validateSharedQuoteEntry(entry SharedQuoteEntry) error {
	switch {
	case strings.TrimSpace(entry.GlobalID) == "":
		return fmt.Errorf("shared quote missing global_id")
	case strings.TrimSpace(entry.AuthorUserID) == "":
		return fmt.Errorf("shared quote %s missing author_user_id", entry.GlobalID)
	case strings.TrimSpace(entry.SourceUserID) == "":
		return fmt.Errorf("shared quote %s missing source_user_id", entry.GlobalID)
	case entry.Version < 1:
		return fmt.Errorf("shared quote %s has invalid version %d", entry.GlobalID, entry.Version)
	case strings.TrimSpace(entry.Content) == "":
		return fmt.Errorf("shared quote %s has empty content", entry.GlobalID)
	}
	return nil
}

func normalizeImportedQuoteIdentity(schemaVersion int, identity db.QuoteIdentity) db.QuoteIdentity {
	if schemaVersion == 1 || identity.SourceBackend == "" {
		identity.SourceBackend = "shared_import"
	}
	if schemaVersion == 1 || identity.SourceNamespace == "" {
		sourceUserID := strings.TrimSpace(identity.SourceUserID)
		if sourceUserID == "" {
			sourceUserID = "unknown"
		}
		identity.SourceNamespace = "share:" + sourceUserID
	}
	if schemaVersion == 1 || identity.SourceEntityType == "" {
		identity.SourceEntityType = "shared_quote"
	}
	if schemaVersion == 1 || identity.SourceEntityID == "" {
		identity.SourceEntityID = identity.GlobalID
	}
	if schemaVersion == 1 || identity.SourceLabel == "" {
		sourceName := strings.TrimSpace(identity.SourceName)
		if sourceName == "" {
			sourceName = "Shared import"
		}
		identity.SourceLabel = sourceName
	}
	return identity
}
