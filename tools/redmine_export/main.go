package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gigol/irecall/core"
)

type config struct {
	host               string
	port               int
	user               string
	password           string
	database           string
	output             string
	baseURL            string
	sslMode            string
	includeIssues      bool
	includeJournals    bool
	includePrivate     bool
	projectIdentifiers stringList
	issueIDs           intList
}

type stringList []string

func (s *stringList) String() string {
	return strings.Join(*s, ",")
}

func (s *stringList) Set(value string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	*s = append(*s, value)
	return nil
}

type intList []int

func (l *intList) String() string {
	out := make([]string, 0, len(*l))
	for _, v := range *l {
		out = append(out, strconv.Itoa(v))
	}
	return strings.Join(out, ",")
}

func (l *intList) Set(value string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	n, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("invalid integer %q", value)
	}
	*l = append(*l, n)
	return nil
}

type issueRow struct {
	IssueID           int
	ProjectIdentifier string
	ProjectName       string
	TrackerName       string
	StatusName        string
	Subject           string
	Description       string
	AuthorID          int
	AuthorLogin       string
	AuthorFirstname   string
	AuthorLastname    string
	CreatedOn         time.Time
	UpdatedOn         time.Time
}

type journalRow struct {
	JournalID         int
	IssueID           int
	ProjectIdentifier string
	ProjectName       string
	TrackerName       string
	StatusName        string
	Subject           string
	Notes             string
	UserID            int
	UserLogin         string
	UserFirstname     string
	UserLastname      string
	CreatedOn         time.Time
	PrivateNotes      bool
}

func main() {
	cfg, err := parseFlags(os.Args[1:])
	if err != nil {
		exitErr(err)
	}
	if err := run(context.Background(), cfg); err != nil {
		exitErr(err)
	}
}

func parseFlags(args []string) (config, error) {
	var cfg config
	fs := flag.NewFlagSet("redmine_export", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	fs.StringVar(&cfg.host, "host", "", "Optional PostgreSQL host or IP; omit to use the local Unix socket")
	fs.IntVar(&cfg.port, "port", 5432, "PostgreSQL port")
	fs.StringVar(&cfg.user, "user", "", "PostgreSQL user")
	fs.StringVar(&cfg.password, "password", "", "PostgreSQL password")
	fs.StringVar(&cfg.database, "database", "", "PostgreSQL database name")
	fs.StringVar(&cfg.output, "output", "", "Output path for iRecall import JSON")
	fs.StringVar(&cfg.baseURL, "base-url", "", "Optional Redmine base URL used to build source URLs")
	fs.StringVar(&cfg.sslMode, "sslmode", "disable", "PostgreSQL sslmode passed to psql")
	fs.BoolVar(&cfg.includeIssues, "include-issues", true, "Include issue descriptions")
	fs.BoolVar(&cfg.includeJournals, "include-journals", true, "Include journal notes")
	fs.BoolVar(&cfg.includePrivate, "include-private-notes", false, "Include Redmine private journal notes")
	fs.Var(&cfg.projectIdentifiers, "project", "Project identifier filter, repeatable")
	fs.Var(&cfg.issueIDs, "issue-id", "Issue ID filter, repeatable")

	if err := fs.Parse(args); err != nil {
		return config{}, err
	}
	if strings.TrimSpace(cfg.user) == "" {
		return config{}, errors.New("missing --user")
	}
	if strings.TrimSpace(cfg.database) == "" {
		return config{}, errors.New("missing --database")
	}
	if strings.TrimSpace(cfg.output) == "" {
		return config{}, errors.New("missing --output")
	}
	if !cfg.includeIssues && !cfg.includeJournals {
		return config{}, errors.New("at least one of --include-issues or --include-journals must be enabled")
	}
	return cfg, nil
}

func run(ctx context.Context, cfg config) error {
	var entries []core.SharedQuoteEntry

	if cfg.includeIssues {
		rows, err := fetchIssueRows(ctx, cfg)
		if err != nil {
			return err
		}
		for _, row := range rows {
			entries = append(entries, issueRowToEntry(cfg, row))
		}
	}

	if cfg.includeJournals {
		rows, err := fetchJournalRows(ctx, cfg)
		if err != nil {
			return err
		}
		for _, row := range rows {
			entries = append(entries, journalRowToEntry(cfg, row))
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		if !entries[i].CreatedAtUTC.Equal(entries[j].CreatedAtUTC) {
			return entries[i].CreatedAtUTC.Before(entries[j].CreatedAtUTC)
		}
		return entries[i].GlobalID < entries[j].GlobalID
	})

	payload, err := json.MarshalIndent(core.SharedQuoteEnvelope{
		SchemaVersion: core.ShareSchemaVersion,
		ExportedAt:    time.Now().UTC(),
		Quotes:        entries,
	}, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal output: %w", err)
	}
	if err := os.WriteFile(cfg.output, payload, 0o600); err != nil {
		return fmt.Errorf("write output: %w", err)
	}
	return nil
}

func fetchIssueRows(ctx context.Context, cfg config) ([]issueRow, error) {
	sql := buildIssueQuery(cfg)
	records, err := runPSQL(ctx, cfg, sql)
	if err != nil {
		return nil, err
	}
	rows := make([]issueRow, 0, len(records))
	for _, record := range records {
		row, err := parseIssueRow(record)
		if err != nil {
			return nil, err
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func fetchJournalRows(ctx context.Context, cfg config) ([]journalRow, error) {
	sql := buildJournalQuery(cfg)
	records, err := runPSQL(ctx, cfg, sql)
	if err != nil {
		return nil, err
	}
	rows := make([]journalRow, 0, len(records))
	for _, record := range records {
		row, err := parseJournalRow(record)
		if err != nil {
			return nil, err
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func buildIssueQuery(cfg config) string {
	var where []string
	where = append(where, "COALESCE(BTRIM(i.description), '') <> ''")
	where = append(where, issueFilters(cfg)...)
	return `
SELECT
  i.id AS issue_id,
  p.identifier AS project_identifier,
  p.name AS project_name,
  t.name AS tracker_name,
  s.name AS status_name,
  i.subject,
  i.description,
  i.author_id,
  COALESCE(u.login, '') AS author_login,
  COALESCE(u.firstname, '') AS author_firstname,
  COALESCE(u.lastname, '') AS author_lastname,
  COALESCE(i.created_on, NOW()) AS created_on,
  COALESCE(i.updated_on, COALESCE(i.created_on, NOW())) AS updated_on
FROM issues i
JOIN projects p ON p.id = i.project_id
JOIN trackers t ON t.id = i.tracker_id
JOIN issue_statuses s ON s.id = i.status_id
LEFT JOIN users u ON u.id = i.author_id
WHERE ` + strings.Join(where, " AND ") + `
ORDER BY i.id`
}

func buildJournalQuery(cfg config) string {
	var where []string
	where = append(where, "j.journalized_type = 'Issue'")
	where = append(where, "COALESCE(BTRIM(j.notes), '') <> ''")
	if !cfg.includePrivate {
		where = append(where, "j.private_notes = FALSE")
	}
	where = append(where, issueFilters(cfg)...)
	return `
SELECT
  j.id AS journal_id,
  i.id AS issue_id,
  p.identifier AS project_identifier,
  p.name AS project_name,
  t.name AS tracker_name,
  s.name AS status_name,
  i.subject,
  j.notes,
  j.user_id,
  COALESCE(u.login, '') AS user_login,
  COALESCE(u.firstname, '') AS user_firstname,
  COALESCE(u.lastname, '') AS user_lastname,
  COALESCE(j.created_on, NOW()) AS created_on,
  j.private_notes
FROM journals j
JOIN issues i ON i.id = j.journalized_id
JOIN projects p ON p.id = i.project_id
JOIN trackers t ON t.id = i.tracker_id
JOIN issue_statuses s ON s.id = i.status_id
LEFT JOIN users u ON u.id = j.user_id
WHERE ` + strings.Join(where, " AND ") + `
ORDER BY j.id`
}

func issueFilters(cfg config) []string {
	var where []string
	if len(cfg.projectIdentifiers) > 0 {
		quoted := make([]string, 0, len(cfg.projectIdentifiers))
		for _, id := range cfg.projectIdentifiers {
			quoted = append(quoted, sqlString(id))
		}
		where = append(where, "p.identifier IN ("+strings.Join(quoted, ", ")+")")
	}
	if len(cfg.issueIDs) > 0 {
		parts := make([]string, 0, len(cfg.issueIDs))
		for _, id := range cfg.issueIDs {
			parts = append(parts, strconv.Itoa(id))
		}
		where = append(where, "i.id IN ("+strings.Join(parts, ", ")+")")
	}
	return where
}

func runPSQL(ctx context.Context, cfg config, sql string) ([]map[string]string, error) {
	cmd := exec.CommandContext(
		ctx,
		"psql",
		"--csv",
		"--no-psqlrc",
		"--port", strconv.Itoa(cfg.port),
		"--username", cfg.user,
		"--dbname", cfg.database,
		"--set", "ON_ERROR_STOP=1",
		"--command", sql,
	)
	if strings.TrimSpace(cfg.host) != "" {
		cmd.Args = append(cmd.Args[:3], append([]string{"--host", cfg.host}, cmd.Args[3:]...)...)
	}
	cmd.Env = append(os.Environ(), "PGPASSWORD="+cfg.password, "PGSSLMODE="+cfg.sslMode)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return nil, errors.New("psql is required but was not found in PATH")
		}
		return nil, fmt.Errorf("psql query failed: %v: %s", err, strings.TrimSpace(stderr.String()))
	}

	records, err := parseCSVRecords(stdout.Bytes())
	if err != nil {
		return nil, err
	}
	return records, nil
}

func parseCSVRecords(data []byte) ([]map[string]string, error) {
	r := csv.NewReader(bytes.NewReader(data))
	rows, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("parse psql csv: %w", err)
	}
	if len(rows) == 0 {
		return nil, nil
	}
	headers := rows[0]
	out := make([]map[string]string, 0, max(0, len(rows)-1))
	for _, row := range rows[1:] {
		record := make(map[string]string, len(headers))
		for i, header := range headers {
			if i < len(row) {
				record[header] = row[i]
			} else {
				record[header] = ""
			}
		}
		out = append(out, record)
	}
	return out, nil
}

func parseIssueRow(record map[string]string) (issueRow, error) {
	issueID, err := parseIntField(record, "issue_id")
	if err != nil {
		return issueRow{}, err
	}
	authorID, err := parseIntField(record, "author_id")
	if err != nil {
		return issueRow{}, err
	}
	createdOn, err := parseTimeField(record, "created_on")
	if err != nil {
		return issueRow{}, err
	}
	updatedOn, err := parseTimeField(record, "updated_on")
	if err != nil {
		return issueRow{}, err
	}
	return issueRow{
		IssueID:           issueID,
		ProjectIdentifier: strings.TrimSpace(record["project_identifier"]),
		ProjectName:       strings.TrimSpace(record["project_name"]),
		TrackerName:       strings.TrimSpace(record["tracker_name"]),
		StatusName:        strings.TrimSpace(record["status_name"]),
		Subject:           strings.TrimSpace(record["subject"]),
		Description:       strings.TrimSpace(record["description"]),
		AuthorID:          authorID,
		AuthorLogin:       strings.TrimSpace(record["author_login"]),
		AuthorFirstname:   strings.TrimSpace(record["author_firstname"]),
		AuthorLastname:    strings.TrimSpace(record["author_lastname"]),
		CreatedOn:         createdOn.UTC(),
		UpdatedOn:         updatedOn.UTC(),
	}, nil
}

func parseJournalRow(record map[string]string) (journalRow, error) {
	journalID, err := parseIntField(record, "journal_id")
	if err != nil {
		return journalRow{}, err
	}
	issueID, err := parseIntField(record, "issue_id")
	if err != nil {
		return journalRow{}, err
	}
	userID, err := parseIntField(record, "user_id")
	if err != nil {
		return journalRow{}, err
	}
	createdOn, err := parseTimeField(record, "created_on")
	if err != nil {
		return journalRow{}, err
	}
	privateNotes, err := parseBoolField(record, "private_notes")
	if err != nil {
		return journalRow{}, err
	}
	return journalRow{
		JournalID:         journalID,
		IssueID:           issueID,
		ProjectIdentifier: strings.TrimSpace(record["project_identifier"]),
		ProjectName:       strings.TrimSpace(record["project_name"]),
		TrackerName:       strings.TrimSpace(record["tracker_name"]),
		StatusName:        strings.TrimSpace(record["status_name"]),
		Subject:           strings.TrimSpace(record["subject"]),
		Notes:             strings.TrimSpace(record["notes"]),
		UserID:            userID,
		UserLogin:         strings.TrimSpace(record["user_login"]),
		UserFirstname:     strings.TrimSpace(record["user_firstname"]),
		UserLastname:      strings.TrimSpace(record["user_lastname"]),
		CreatedOn:         createdOn.UTC(),
		PrivateNotes:      privateNotes,
	}, nil
}

func parseIntField(record map[string]string, key string) (int, error) {
	v := strings.TrimSpace(record[key])
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("parse %s as int: %w", key, err)
	}
	return n, nil
}

func parseBoolField(record map[string]string, key string) (bool, error) {
	v := strings.TrimSpace(record[key])
	b, err := strconv.ParseBool(v)
	if err != nil {
		return false, fmt.Errorf("parse %s as bool: %w", key, err)
	}
	return b, nil
}

func parseTimeField(record map[string]string, key string) (time.Time, error) {
	raw := strings.TrimSpace(record[key])
	layouts := []string{
		time.RFC3339Nano,
		"2006-01-02 15:04:05.999999-07",
		"2006-01-02 15:04:05.999999-07:00",
		"2006-01-02 15:04:05-07",
		"2006-01-02 15:04:05-07:00",
		"2006-01-02 15:04:05.999999",
		"2006-01-02 15:04:05",
	}
	for _, layout := range layouts {
		if ts, err := time.Parse(layout, raw); err == nil {
			return ts, nil
		}
	}
	return time.Time{}, fmt.Errorf("parse %s as time: unsupported format %q", key, raw)
}

func issueRowToEntry(cfg config, row issueRow) core.SharedQuoteEntry {
	authorID := redmineUserID(row.AuthorID)
	authorName := redmineUserName(row.AuthorFirstname, row.AuthorLastname, row.AuthorLogin, row.AuthorID)
	issueURL := joinURL(cfg.baseURL, fmt.Sprintf("/issues/%d", row.IssueID))

	return core.SharedQuoteEntry{
		GlobalID:         fmt.Sprintf("redmine:%s:issue:%d:description", cfg.database, row.IssueID),
		AuthorUserID:     authorID,
		AuthorName:       authorName,
		SourceUserID:     authorID,
		SourceName:       authorName,
		SourceBackend:    "redmine",
		SourceNamespace:  "redmine:" + cfg.database,
		SourceEntityType: "issue_description",
		SourceEntityID:   strconv.Itoa(row.IssueID),
		SourceLabel:      fmt.Sprintf("Redmine issue #%d", row.IssueID),
		SourceURL:        issueURL,
		Version:          1,
		Content:          buildIssueContent(row),
		Tags:             buildTags("issue_description", row.ProjectIdentifier, row.ProjectName, row.TrackerName, row.StatusName),
		CreatedAtUTC:     row.CreatedOn.UTC(),
		UpdatedAtUTC:     row.UpdatedOn.UTC(),
	}
}

func journalRowToEntry(cfg config, row journalRow) core.SharedQuoteEntry {
	authorID := redmineUserID(row.UserID)
	authorName := redmineUserName(row.UserFirstname, row.UserLastname, row.UserLogin, row.UserID)
	journalURL := ""
	if base := joinURL(cfg.baseURL, fmt.Sprintf("/issues/%d", row.IssueID)); base != "" {
		journalURL = base + "#note-" + strconv.Itoa(row.JournalID)
	}

	return core.SharedQuoteEntry{
		GlobalID:         fmt.Sprintf("redmine:%s:journal:%d", cfg.database, row.JournalID),
		AuthorUserID:     authorID,
		AuthorName:       authorName,
		SourceUserID:     authorID,
		SourceName:       authorName,
		SourceBackend:    "redmine",
		SourceNamespace:  "redmine:" + cfg.database,
		SourceEntityType: "issue_journal",
		SourceEntityID:   strconv.Itoa(row.JournalID),
		SourceLabel:      fmt.Sprintf("Redmine journal #%d", row.JournalID),
		SourceURL:        journalURL,
		Version:          1,
		Content:          buildJournalContent(row, authorName),
		Tags:             buildTags("issue_journal", row.ProjectIdentifier, row.ProjectName, row.TrackerName, row.StatusName),
		CreatedAtUTC:     row.CreatedOn.UTC(),
		UpdatedAtUTC:     row.CreatedOn.UTC(),
	}
}

func buildIssueContent(row issueRow) string {
	var parts []string
	title := fmt.Sprintf("Redmine issue #%d: %s", row.IssueID, strings.TrimSpace(row.Subject))
	parts = append(parts, title)
	meta := []string{}
	if row.ProjectName != "" {
		meta = append(meta, "Project: "+row.ProjectName)
	}
	if row.TrackerName != "" {
		meta = append(meta, "Tracker: "+row.TrackerName)
	}
	if row.StatusName != "" {
		meta = append(meta, "Status: "+row.StatusName)
	}
	if len(meta) > 0 {
		parts = append(parts, strings.Join(meta, "\n"))
	}
	parts = append(parts, strings.TrimSpace(row.Description))
	return strings.Join(parts, "\n\n")
}

func buildJournalContent(row journalRow, authorName string) string {
	var parts []string
	parts = append(parts, fmt.Sprintf("Redmine issue #%d: %s", row.IssueID, strings.TrimSpace(row.Subject)))
	meta := []string{}
	if row.ProjectName != "" {
		meta = append(meta, "Project: "+row.ProjectName)
	}
	if row.TrackerName != "" {
		meta = append(meta, "Tracker: "+row.TrackerName)
	}
	if row.StatusName != "" {
		meta = append(meta, "Status: "+row.StatusName)
	}
	if authorName != "" {
		meta = append(meta, "Note by: "+authorName)
	}
	if len(meta) > 0 {
		parts = append(parts, strings.Join(meta, "\n"))
	}
	parts = append(parts, strings.TrimSpace(row.Notes))
	return strings.Join(parts, "\n\n")
}

func buildTags(entityType, projectIdentifier, projectName, trackerName, statusName string) []string {
	candidates := []string{
		"redmine",
		entityType,
		projectIdentifier,
		projectName,
		trackerName,
		statusName,
	}
	seen := make(map[string]struct{}, len(candidates))
	var tags []string
	for _, candidate := range candidates {
		tag := normalizeTag(candidate)
		if tag == "" {
			continue
		}
		if _, ok := seen[tag]; ok {
			continue
		}
		seen[tag] = struct{}{}
		tags = append(tags, tag)
	}
	return tags
}

func normalizeTag(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	replacer := strings.NewReplacer(
		" ", "-",
		"/", "-",
		"_", "-",
		":", "-",
	)
	s = replacer.Replace(s)
	for strings.Contains(s, "--") {
		s = strings.ReplaceAll(s, "--", "-")
	}
	return strings.Trim(s, "-")
}

func redmineUserID(id int) string {
	return fmt.Sprintf("redmine:user:%d", id)
}

func redmineUserName(firstname, lastname, login string, id int) string {
	full := strings.TrimSpace(strings.TrimSpace(firstname) + " " + strings.TrimSpace(lastname))
	switch {
	case full != "":
		return full
	case strings.TrimSpace(login) != "":
		return strings.TrimSpace(login)
	default:
		return fmt.Sprintf("Redmine User %d", id)
	}
}

func joinURL(base, path string) string {
	base = strings.TrimRight(strings.TrimSpace(base), "/")
	if base == "" {
		return ""
	}
	return base + path
}

func sqlString(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func exitErr(err error) {
	fmt.Fprintln(os.Stderr, "error:", err)
	os.Exit(1)
}
