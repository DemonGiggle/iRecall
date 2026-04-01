package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

// Store wraps the SQLite database and exposes all persistence operations.
type Store struct {
	db *sql.DB
}

// Open opens (or creates) the SQLite database at path and runs migrations.
func Open(path string) (*Store, error) {
	slog.Info("db: opening database", "path", path)
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	if _, err := db.Exec(`PRAGMA journal_mode = WAL`); err != nil {
		return nil, err
	}
	if _, err := db.Exec(`PRAGMA foreign_keys = ON`); err != nil {
		return nil, err
	}
	if _, err := db.Exec(`PRAGMA busy_timeout = 5000`); err != nil {
		return nil, err
	}
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}
	slog.Info("db: database ready")
	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	slog.Info("db: closing database")
	return s.db.Close()
}

// --- Quotes ---

// InsertQuote stores a new quote and returns its assigned ID.
func (s *Store) InsertQuote(content string) (int64, error) {
	now := time.Now().Unix()
	slog.Info("db: inserting quote", "content_len", len(content))
	res, err := s.db.Exec(
		`INSERT INTO quotes(content, created_at, updated_at) VALUES (?, ?, ?)`,
		content, now, now,
	)
	if err != nil {
		slog.Error("db: insert quote failed", "error", err)
		return 0, fmt.Errorf("insert quote: %w", err)
	}
	id, _ := res.LastInsertId()
	slog.Info("db: quote inserted", "id", id)
	return id, nil
}

// UpdateQuoteFTS refreshes the FTS index for a quote with its current tags.
// Must be called after tag associations are saved.
func (s *Store) UpdateQuoteFTS(id int64, tags []string) error {
	slog.Debug("db: updating FTS for quote", "id", id, "tags", tags)
	row := s.db.QueryRow(`SELECT content FROM quotes WHERE id = ?`, id)
	var content string
	if err := row.Scan(&content); err != nil {
		slog.Error("db: fetch quote for FTS failed", "id", id, "error", err)
		return fmt.Errorf("fetch quote for fts: %w", err)
	}
	tagStr := strings.Join(tags, " ")
	// delete old FTS entry then reinsert with tag text
	if _, err := s.db.Exec(
		`INSERT INTO quotes_fts(quotes_fts, rowid, content, tags) VALUES ('delete', ?, ?, '')`,
		id, content,
	); err != nil {
		slog.Error("db: FTS delete failed", "id", id, "error", err)
		return fmt.Errorf("fts delete: %w", err)
	}
	if _, err := s.db.Exec(
		`INSERT INTO quotes_fts(rowid, content, tags) VALUES (?, ?, ?)`,
		id, content, tagStr,
	); err != nil {
		slog.Error("db: FTS insert failed", "id", id, "error", err)
		return fmt.Errorf("fts insert: %w", err)
	}
	slog.Info("db: FTS updated", "id", id, "tag_str", tagStr)
	return nil
}

// DeleteQuote removes a quote and its tag associations.
func (s *Store) DeleteQuote(id int64) error {
	slog.Info("db: deleting quote", "id", id)
	_, err := s.db.Exec(`DELETE FROM quotes WHERE id = ?`, id)
	if err != nil {
		slog.Error("db: delete quote failed", "id", id, "error", err)
	}
	return err
}

// ListQuotes returns all quotes with their tags, newest first.
func (s *Store) ListQuotes() ([]QuoteRow, error) {
	slog.Debug("db: listing all quotes")
	rows, err := s.db.Query(`
		SELECT q.id, q.content, q.created_at, q.updated_at,
		       COALESCE(GROUP_CONCAT(t.name, ','), '') AS tags
		FROM quotes q
		LEFT JOIN quote_tags qt ON qt.quote_id = q.id
		LEFT JOIN tags t        ON t.id = qt.tag_id
		GROUP BY q.id
		ORDER BY q.created_at DESC
	`)
	if err != nil {
		slog.Error("db: list quotes failed", "error", err)
		return nil, err
	}
	defer rows.Close()
	result, err := scanQuoteRows(rows)
	slog.Debug("db: listed quotes", "count", len(result))
	return result, err
}

// SearchQuotes performs a ranked FTS5 query and returns up to limit results.
// keywords is joined as "kw1 OR kw2 OR ..." before querying.
func (s *Store) SearchQuotes(keywords []string, limit int) ([]QuoteRow, error) {
	slog.Info("db: searching quotes", "keywords", keywords, "limit", limit)
	if len(keywords) == 0 {
		slog.Debug("db: search skipped, no keywords")
		return nil, nil
	}
	// Quote each keyword with double quotes to escape FTS5 special characters.
	// Any embedded double quotes are doubled to escape them per FTS5 syntax.
	quoted := make([]string, 0, len(keywords))
	for _, kw := range keywords {
		kw = strings.TrimSpace(kw)
		if kw == "" {
			continue
		}
		kw = strings.ReplaceAll(kw, `"`, `""`)
		quoted = append(quoted, `"`+kw+`"`)
	}
	if len(quoted) == 0 {
		slog.Debug("db: search skipped, all keywords empty after trim")
		return nil, nil
	}
	query := strings.Join(quoted, " OR ")
	slog.Debug("db: FTS query", "match_expr", query)
	rows, err := s.db.Query(`
		SELECT q.id, q.content, q.created_at, q.updated_at,
		       COALESCE(GROUP_CONCAT(t.name, ','), '') AS tags
		FROM quotes_fts AS fts
		JOIN quotes AS q     ON q.id = fts.rowid
		LEFT JOIN quote_tags qt ON qt.quote_id = q.id
		LEFT JOIN tags t        ON t.id = qt.tag_id
		WHERE quotes_fts MATCH ?
		GROUP BY q.id
		ORDER BY fts.rank
		LIMIT ?
	`, query, limit)
	if err != nil {
		slog.Error("db: FTS search failed", "query", query, "error", err)
		return nil, fmt.Errorf("fts search: %w", err)
	}
	defer rows.Close()
	result, err := scanQuoteRows(rows)
	if err != nil {
		slog.Error("db: scan search results failed", "error", err)
		return nil, err
	}
	slog.Info("db: search results", "match_count", len(result))
	for i, r := range result {
		slog.Debug("db: search result", "index", i, "id", r.ID, "content_preview", truncate(r.Content, 80), "tags", r.Tags)
	}
	return result, nil
}

// QuoteRow is the raw DB representation returned by list/search queries.
type QuoteRow struct {
	ID        int64
	Content   string
	CreatedAt int64
	UpdatedAt int64
	Tags      string // comma-separated
}

func scanQuoteRows(rows *sql.Rows) ([]QuoteRow, error) {
	var out []QuoteRow
	for rows.Next() {
		var r QuoteRow
		if err := rows.Scan(&r.ID, &r.Content, &r.CreatedAt, &r.UpdatedAt, &r.Tags); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// --- Tags ---

// UpsertTags inserts tags that don't exist yet and returns all their IDs.
func (s *Store) UpsertTags(names []string) ([]int64, error) {
	slog.Debug("db: upserting tags", "tags", names)
	ids := make([]int64, 0, len(names))
	for _, name := range names {
		name = strings.ToLower(strings.TrimSpace(name))
		if name == "" {
			continue
		}
		_, err := s.db.Exec(`INSERT OR IGNORE INTO tags(name) VALUES (?)`, name)
		if err != nil {
			slog.Error("db: upsert tag failed", "tag", name, "error", err)
			return nil, fmt.Errorf("upsert tag %q: %w", name, err)
		}
		var id int64
		if err := s.db.QueryRow(`SELECT id FROM tags WHERE name = ?`, name).Scan(&id); err != nil {
			slog.Error("db: fetch tag id failed", "tag", name, "error", err)
			return nil, fmt.Errorf("fetch tag id %q: %w", name, err)
		}
		ids = append(ids, id)
	}
	slog.Debug("db: upserted tags", "tag_ids", ids)
	return ids, nil
}

// InsertQuoteTags creates the many-to-many associations.
func (s *Store) InsertQuoteTags(quoteID int64, tagIDs []int64) error {
	slog.Debug("db: inserting quote-tag associations", "quote_id", quoteID, "tag_ids", tagIDs)
	for _, tid := range tagIDs {
		if _, err := s.db.Exec(
			`INSERT OR IGNORE INTO quote_tags(quote_id, tag_id) VALUES (?, ?)`,
			quoteID, tid,
		); err != nil {
			slog.Error("db: insert quote_tag failed", "quote_id", quoteID, "tag_id", tid, "error", err)
			return fmt.Errorf("insert quote_tag: %w", err)
		}
	}
	return nil
}

// --- Settings ---

func (s *Store) GetSetting(key string) (string, error) {
	slog.Debug("db: getting setting", "key", key)
	var val string
	err := s.db.QueryRow(`SELECT value FROM settings WHERE key = ?`, key).Scan(&val)
	if err == sql.ErrNoRows {
		slog.Debug("db: setting not found", "key", key)
		return "", nil
	}
	if err != nil {
		slog.Error("db: get setting failed", "key", key, "error", err)
	}
	return val, err
}

func (s *Store) SetSetting(key, value string) error {
	slog.Debug("db: setting value", "key", key, "value_len", len(value))
	_, err := s.db.Exec(
		`INSERT INTO settings(key, value) VALUES (?, ?)
		 ON CONFLICT(key) DO UPDATE SET value = excluded.value`,
		key, value,
	)
	if err != nil {
		slog.Error("db: set setting failed", "key", key, "error", err)
	}
	return err
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
