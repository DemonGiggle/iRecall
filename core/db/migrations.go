package db

import (
	"database/sql"
	"fmt"
	"slices"
	"time"
)

type migration struct {
	version int
	name    string
	up      func(*sql.Tx) error
}

var migrations = []migration{
	{
		version: 1,
		name:    "initial_schema",
		up:      upInitialSchema,
	},
}

const initialSchemaSQL = `
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
`

func runMigrations(db *sql.DB) error {
	if err := validateMigrations(migrations); err != nil {
		return err
	}
	if err := ensureMigrationTables(db); err != nil {
		return err
	}
	if err := importLegacySchemaVersions(db); err != nil {
		return err
	}

	applied, err := appliedMigrationVersions(db)
	if err != nil {
		return err
	}

	for _, m := range migrations {
		if _, ok := applied[m.version]; ok {
			continue
		}
		if err := applyMigration(db, m); err != nil {
			return err
		}
	}

	return nil
}

func validateMigrations(ms []migration) error {
	if len(ms) == 0 {
		return nil
	}

	versions := make([]int, 0, len(ms))
	for _, m := range ms {
		if m.version < 1 {
			return fmt.Errorf("invalid migration version %d", m.version)
		}
		if m.name == "" {
			return fmt.Errorf("migration v%d has empty name", m.version)
		}
		if m.up == nil {
			return fmt.Errorf("migration v%d has nil up function", m.version)
		}
		versions = append(versions, m.version)
	}

	sorted := slices.Clone(versions)
	slices.Sort(sorted)
	for i, v := range sorted {
		if i > 0 && v == sorted[i-1] {
			return fmt.Errorf("duplicate migration version %d", v)
		}
		if v != i+1 {
			return fmt.Errorf("migration versions must be contiguous starting at 1: got %d at position %d", v, i+1)
		}
	}

	return nil
}

func ensureMigrationTables(db *sql.DB) error {
	if _, err := db.Exec(`
CREATE TABLE IF NOT EXISTS schema_migrations (
    version    INTEGER PRIMARY KEY,
    name       TEXT NOT NULL,
    applied_at INTEGER NOT NULL
)`); err != nil {
		return fmt.Errorf("create schema_migrations: %w", err)
	}

	// Keep the legacy table in place for backward compatibility with databases
	// created by the original migration runner.
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_version (version INTEGER NOT NULL)`); err != nil {
		return fmt.Errorf("create schema_version: %w", err)
	}

	return nil
}

func importLegacySchemaVersions(db *sql.DB) error {
	var imported int
	if err := db.QueryRow(`SELECT COUNT(*) FROM schema_migrations`).Scan(&imported); err != nil {
		return fmt.Errorf("count schema_migrations: %w", err)
	}
	if imported > 0 {
		return nil
	}

	rows, err := db.Query(`SELECT DISTINCT version FROM schema_version ORDER BY version`)
	if err != nil {
		return fmt.Errorf("read legacy schema_version: %w", err)
	}
	defer rows.Close()

	var legacyVersions []int
	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return fmt.Errorf("scan legacy schema version: %w", err)
		}
		legacyVersions = append(legacyVersions, version)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate legacy schema versions: %w", err)
	}
	if len(legacyVersions) == 0 {
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin legacy import: %w", err)
	}
	defer tx.Rollback()

	appliedAt := time.Now().Unix()
	for _, version := range legacyVersions {
		name := fmt.Sprintf("legacy_v%d", version)
		if _, err := tx.Exec(
			`INSERT OR IGNORE INTO schema_migrations(version, name, applied_at) VALUES (?, ?, ?)`,
			version, name, appliedAt,
		); err != nil {
			return fmt.Errorf("import legacy schema version %d: %w", version, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit legacy import: %w", err)
	}
	return nil
}

func appliedMigrationVersions(db *sql.DB) (map[int]string, error) {
	rows, err := db.Query(`SELECT version, name FROM schema_migrations ORDER BY version`)
	if err != nil {
		return nil, fmt.Errorf("read applied migrations: %w", err)
	}
	defer rows.Close()

	applied := make(map[int]string)
	for rows.Next() {
		var version int
		var name string
		if err := rows.Scan(&version, &name); err != nil {
			return nil, fmt.Errorf("scan applied migration: %w", err)
		}
		applied[version] = name
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate applied migrations: %w", err)
	}
	return applied, nil
}

func applyMigration(db *sql.DB, m migration) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin migration v%d: %w", m.version, err)
	}
	defer tx.Rollback()

	if err := m.up(tx); err != nil {
		return fmt.Errorf("migration v%d (%s): %w", m.version, m.name, err)
	}
	if _, err := tx.Exec(
		`INSERT INTO schema_migrations(version, name, applied_at) VALUES (?, ?, ?)`,
		m.version, m.name, time.Now().Unix(),
	); err != nil {
		return fmt.Errorf("record migration v%d (%s): %w", m.version, m.name, err)
	}
	if _, err := tx.Exec(`INSERT INTO schema_version(version) VALUES (?)`, m.version); err != nil {
		return fmt.Errorf("record legacy schema version v%d (%s): %w", m.version, m.name, err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit migration v%d (%s): %w", m.version, m.name, err)
	}
	return nil
}

func upInitialSchema(tx *sql.Tx) error {
	if _, err := tx.Exec(initialSchemaSQL); err != nil {
		return err
	}
	return nil
}
