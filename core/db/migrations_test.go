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

func countRows(t *testing.T, db *sql.DB, table string) int {
	t.Helper()

	var n int
	query := `SELECT COUNT(*) FROM ` + table
	if err := db.QueryRow(query).Scan(&n); err != nil {
		t.Fatalf("count rows in %s: %v", table, err)
	}
	return n
}
