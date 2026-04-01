package db

import (
	"database/sql"
	"fmt"
)

type migration struct {
	version int
	up      string
}

var migrations = []migration{
	{
		version: 1,
		up: `
CREATE TABLE IF NOT EXISTS quotes (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    content    TEXT    NOT NULL,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS tags (
    id   INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT    NOT NULL UNIQUE COLLATE NOCASE
);

CREATE TABLE IF NOT EXISTS quote_tags (
    quote_id INTEGER NOT NULL REFERENCES quotes(id) ON DELETE CASCADE,
    tag_id   INTEGER NOT NULL REFERENCES tags(id)   ON DELETE CASCADE,
    PRIMARY KEY (quote_id, tag_id)
);

CREATE VIRTUAL TABLE IF NOT EXISTS quotes_fts USING fts5(
    content,
    tags,
    content='quotes',
    content_rowid='id',
    tokenize='porter unicode61'
);

CREATE TRIGGER IF NOT EXISTS quotes_ai AFTER INSERT ON quotes BEGIN
    INSERT INTO quotes_fts(rowid, content, tags)
    VALUES (new.id, new.content, '');
END;

CREATE TRIGGER IF NOT EXISTS quotes_ad AFTER DELETE ON quotes BEGIN
    INSERT INTO quotes_fts(quotes_fts, rowid, content, tags)
    VALUES ('delete', old.id, old.content, '');
END;

CREATE TRIGGER IF NOT EXISTS quotes_au AFTER UPDATE ON quotes BEGIN
    INSERT INTO quotes_fts(quotes_fts, rowid, content, tags)
    VALUES ('delete', old.id, old.content, '');
    INSERT INTO quotes_fts(rowid, content, tags)
    VALUES (new.id, new.content, '');
END;

CREATE TABLE IF NOT EXISTS settings (
    key   TEXT PRIMARY KEY,
    value TEXT NOT NULL
);
`,
	},
}

func runMigrations(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_version (version INTEGER NOT NULL)`)
	if err != nil {
		return fmt.Errorf("create schema_version: %w", err)
	}

	var current int
	row := db.QueryRow(`SELECT COALESCE(MAX(version), 0) FROM schema_version`)
	if err := row.Scan(&current); err != nil {
		return fmt.Errorf("read schema version: %w", err)
	}

	for _, m := range migrations {
		if m.version <= current {
			continue
		}
		if _, err := db.Exec(m.up); err != nil {
			return fmt.Errorf("migration v%d: %w", m.version, err)
		}
		if _, err := db.Exec(`INSERT INTO schema_version(version) VALUES (?)`, m.version); err != nil {
			return fmt.Errorf("record migration v%d: %w", m.version, err)
		}
	}
	return nil
}
