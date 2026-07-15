package db

import (
	"database/sql"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

func TestValidateMigrationsRejectsGaps(t *testing.T) {
	t.Parallel()

	err := validateMigrations([]migration{
		{version: 1, name: "one", up: upInitialSchema},
		{version: 3, name: "three", up: upInitialSchema},
	})
	if err == nil {
		t.Fatal("expected gap in migration versions to fail validation")
	}
}

func TestRunMigrationsFreshDatabase(t *testing.T) {
	t.Parallel()

	db := openTestSQLiteDB(t)
	if err := configureConnection(db); err != nil {
		t.Fatalf("configure connection: %v", err)
	}

	if err := runMigrations(db); err != nil {
		t.Fatalf("run migrations: %v", err)
	}

	assertTableExists(t, db, "quotes")
	assertTableExists(t, db, "tags")
	assertTableExists(t, db, "quote_tags")
	assertTableExists(t, db, "quotes_fts")
	assertTableExists(t, db, "settings")
	assertTableExists(t, db, "schema_migrations")
	assertTableExists(t, db, "user_profile")

	if got := countRows(t, db, "schema_migrations"); got != len(migrations) {
		t.Fatalf("schema_migrations row count = %d, want %d", got, len(migrations))
	}
	if got := countRows(t, db, "schema_version"); got != len(migrations) {
		t.Fatalf("schema_version row count = %d, want %d", got, len(migrations))
	}
}

func TestRunMigrationsImportsLegacySchemaVersion(t *testing.T) {
	t.Parallel()

	db := openTestSQLiteDB(t)
	if err := configureConnection(db); err != nil {
		t.Fatalf("configure connection: %v", err)
	}

	if _, err := db.Exec(`CREATE TABLE schema_version (version INTEGER NOT NULL)`); err != nil {
		t.Fatalf("create legacy schema_version: %v", err)
	}
	if _, err := db.Exec(initialSchemaSQL); err != nil {
		t.Fatalf("create legacy schema: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO schema_version(version) VALUES (1)`); err != nil {
		t.Fatalf("insert legacy schema version: %v", err)
	}

	if err := runMigrations(db); err != nil {
		t.Fatalf("run migrations: %v", err)
	}

	if got := countRows(t, db, "schema_migrations"); got != len(migrations) {
		t.Fatalf("schema_migrations row count = %d, want %d", got, len(migrations))
	}
	if got := countRows(t, db, "schema_version"); got != len(migrations) {
		t.Fatalf("schema_version row count = %d, want %d", got, len(migrations))
	}

	var name string
	if err := db.QueryRow(`SELECT name FROM schema_migrations WHERE version = 1`).Scan(&name); err != nil {
		t.Fatalf("query imported migration name: %v", err)
	}
	if name != "legacy_v1" {
		t.Fatalf("imported migration name = %q, want %q", name, "legacy_v1")
	}
	if err := db.QueryRow(`SELECT name FROM schema_migrations WHERE version = 2`).Scan(&name); err != nil {
		t.Fatalf("query migration v2 name: %v", err)
	}
	if name != "quote_identity_and_user_profile" {
		t.Fatalf("migration v2 name = %q", name)
	}
}

func TestRunMigrationsV4DatabasePreservesExistingData(t *testing.T) {
	t.Parallel()

	db := openTestSQLiteDB(t)
	if err := configureConnection(db); err != nil {
		t.Fatalf("configure connection: %v", err)
	}
	if err := ensureMigrationTables(db); err != nil {
		t.Fatalf("ensure migration tables: %v", err)
	}
	for _, m := range migrations[:4] {
		if err := applyMigration(db, m); err != nil {
			t.Fatalf("apply v4 fixture migration %d: %v", m.version, err)
		}
	}

	const settingsJSON = `{"Provider":{"Host":"legacy-provider","Port":11434,"HTTPS":false,"APIKey":"","Model":"legacy-model","KeywordModel":"legacy-keyword-model"},"Search":{"MaxResults":7,"MinRelevance":0.25},"Debug":{"MockLLM":true},"Theme":"forest","Web":{"Port":9731},"RootDir":"/legacy/root"}`
	if _, err := db.Exec(`INSERT INTO settings(key, value) VALUES ('settings', ?)`, settingsJSON); err != nil {
		t.Fatalf("insert legacy settings: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO user_profile(user_id, display_name, created_at, updated_at) VALUES ('legacy-user', 'Legacy User', 100, 200)`); err != nil {
		t.Fatalf("insert legacy user profile: %v", err)
	}
	result, err := db.Exec(`INSERT INTO quotes(
		content, global_id, author_user_id, author_name, source_user_id, source_name,
		source_backend, source_namespace, source_entity_type, source_entity_id, source_label, source_url,
		version, created_at, updated_at
	) VALUES ('legacy searchable quote', 'legacy-global-id', 'legacy-user', 'Legacy User',
		'legacy-user', 'Legacy User', 'local', 'local:legacy-user', 'quote', 'legacy-global-id',
		'Local quote', '', 3, 300, 400)`)
	if err != nil {
		t.Fatalf("insert legacy quote: %v", err)
	}
	quoteID, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("legacy quote id: %v", err)
	}
	result, err = db.Exec(`INSERT INTO tags(name) VALUES ('compatibility')`)
	if err != nil {
		t.Fatalf("insert legacy tag: %v", err)
	}
	tagID, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("legacy tag id: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO quote_tags(quote_id, tag_id) VALUES (?, ?)`, quoteID, tagID); err != nil {
		t.Fatalf("link legacy tag: %v", err)
	}
	result, err = db.Exec(`INSERT INTO recall_history(question, response, created_at) VALUES ('legacy question', 'legacy response', 500)`)
	if err != nil {
		t.Fatalf("insert legacy recall history: %v", err)
	}
	historyID, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("legacy history id: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO recall_history_quotes(history_id, quote_id, position) VALUES (?, ?, 0)`, historyID, quoteID); err != nil {
		t.Fatalf("link legacy recall history: %v", err)
	}

	if err := runMigrations(db); err != nil {
		t.Fatalf("upgrade v4 database to current schema: %v", err)
	}

	if got := countRows(t, db, "schema_migrations"); got != len(migrations) {
		t.Fatalf("schema migration count = %d, want %d", got, len(migrations))
	}
	assertTableExists(t, db, "quote_attachments")
	assertIndexExists(t, db, "idx_quote_attachments_quote_id")

	var gotSettings string
	if err := db.QueryRow(`SELECT value FROM settings WHERE key = 'settings'`).Scan(&gotSettings); err != nil {
		t.Fatalf("read settings after upgrade: %v", err)
	}
	if gotSettings != settingsJSON {
		t.Fatalf("settings changed during upgrade:\n got %s\nwant %s", gotSettings, settingsJSON)
	}

	var content, globalID string
	var version int64
	if err := db.QueryRow(`SELECT content, global_id, version FROM quotes WHERE id = ?`, quoteID).Scan(&content, &globalID, &version); err != nil {
		t.Fatalf("read quote after upgrade: %v", err)
	}
	if content != "legacy searchable quote" || globalID != "legacy-global-id" || version != 3 {
		t.Fatalf("quote changed during upgrade: content=%q global_id=%q version=%d", content, globalID, version)
	}

	var tag string
	if err := db.QueryRow(`SELECT t.name FROM tags t JOIN quote_tags qt ON qt.tag_id = t.id WHERE qt.quote_id = ?`, quoteID).Scan(&tag); err != nil {
		t.Fatalf("read quote tag after upgrade: %v", err)
	}
	if tag != "compatibility" {
		t.Fatalf("tag after upgrade = %q, want compatibility", tag)
	}

	var response string
	if err := db.QueryRow(`SELECT h.response FROM recall_history h JOIN recall_history_quotes rhq ON rhq.history_id = h.id WHERE rhq.quote_id = ?`, quoteID).Scan(&response); err != nil {
		t.Fatalf("read recall history after upgrade: %v", err)
	}
	if response != "legacy response" {
		t.Fatalf("recall response after upgrade = %q, want legacy response", response)
	}

	var searchMatches int
	if err := db.QueryRow(`SELECT COUNT(*) FROM quotes_fts WHERE quotes_fts MATCH 'searchable'`).Scan(&searchMatches); err != nil {
		t.Fatalf("search FTS after upgrade: %v", err)
	}
	if searchMatches != 1 {
		t.Fatalf("FTS matches after upgrade = %d, want 1", searchMatches)
	}

	var displayName string
	if err := db.QueryRow(`SELECT display_name FROM user_profile WHERE user_id = 'legacy-user'`).Scan(&displayName); err != nil {
		t.Fatalf("read user profile after upgrade: %v", err)
	}
	if displayName != "Legacy User" {
		t.Fatalf("display name after upgrade = %q, want Legacy User", displayName)
	}
}

func openTestSQLiteDB(t *testing.T) *sql.DB {
	t.Helper()

	path := filepath.Join(t.TempDir(), "test.db")
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}
	t.Cleanup(func() {
		db.Close()
	})
	return db
}

func assertTableExists(t *testing.T, db *sql.DB, name string) {
	t.Helper()

	var got string
	if err := db.QueryRow(`SELECT name FROM sqlite_master WHERE name = ?`, name).Scan(&got); err != nil {
		t.Fatalf("table %q not found: %v", name, err)
	}
}

func assertIndexExists(t *testing.T, db *sql.DB, name string) {
	t.Helper()

	var got string
	if err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type = 'index' AND name = ?`, name).Scan(&got); err != nil {
		t.Fatalf("index %q not found: %v", name, err)
	}
}

func countRows(t *testing.T, db *sql.DB, table string) int {
	t.Helper()

	var n int
	query := `SELECT COUNT(*) FROM ` + table
	if err := db.QueryRow(query).Scan(&n); err != nil {
		t.Fatalf("count rows in %s: %v", table, err)
	}
	return n
}
