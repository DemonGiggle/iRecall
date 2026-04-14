package main

import (
	"reflect"
	"testing"
	"time"
)

func TestBuildIssueEntry(t *testing.T) {
	cfg := config{database: "redmine_production", baseURL: "https://redmine.example.com"}
	row := issueRow{
		IssueID:           123,
		ProjectIdentifier: "core-app",
		ProjectName:       "Core App",
		TrackerName:       "Bug",
		StatusName:        "In Progress",
		Subject:           "Crash on startup",
		Description:       "App crashes after launch.",
		AuthorID:          7,
		AuthorLogin:       "alice",
		AuthorFirstname:   "Alice",
		AuthorLastname:    "Chen",
		CreatedOn:         time.Unix(100, 0).UTC(),
		UpdatedOn:         time.Unix(200, 0).UTC(),
	}

	entry := issueRowToEntry(cfg, row)
	if entry.GlobalID != "redmine:redmine_production:issue:123:description" {
		t.Fatalf("GlobalID = %q", entry.GlobalID)
	}
	if entry.SourceBackend != "redmine" || entry.SourceEntityType != "issue_description" || entry.SourceEntityID != "123" {
		t.Fatalf("source provenance = %+v", entry)
	}
	if entry.SourceURL != "https://redmine.example.com/issues/123" {
		t.Fatalf("SourceURL = %q", entry.SourceURL)
	}
	if entry.AuthorUserID != "redmine:user:7" || entry.AuthorName != "Alice Chen" {
		t.Fatalf("author = %q / %q", entry.AuthorUserID, entry.AuthorName)
	}
	if !reflect.DeepEqual(entry.Tags, []string{"redmine", "issue-description", "core-app", "bug", "in-progress"}) {
		t.Fatalf("tags = %#v", entry.Tags)
	}
}

func TestBuildJournalEntry(t *testing.T) {
	cfg := config{database: "redmine_production", baseURL: "https://redmine.example.com/"}
	row := journalRow{
		JournalID:         456,
		IssueID:           123,
		ProjectIdentifier: "core-app",
		ProjectName:       "Core App",
		TrackerName:       "Feature",
		StatusName:        "Resolved",
		Subject:           "Export support",
		Notes:             "Implemented the first pass.",
		UserID:            8,
		UserLogin:         "bob",
		UserFirstname:     "",
		UserLastname:      "",
		CreatedOn:         time.Unix(300, 0).UTC(),
	}

	entry := journalRowToEntry(cfg, row)
	if entry.GlobalID != "redmine:redmine_production:journal:456" {
		t.Fatalf("GlobalID = %q", entry.GlobalID)
	}
	if entry.SourceURL != "https://redmine.example.com/issues/123#note-456" {
		t.Fatalf("SourceURL = %q", entry.SourceURL)
	}
	if entry.AuthorName != "bob" {
		t.Fatalf("AuthorName = %q", entry.AuthorName)
	}
	if entry.UpdatedAtUTC != entry.CreatedAtUTC {
		t.Fatalf("UpdatedAtUTC = %v, want CreatedAtUTC %v", entry.UpdatedAtUTC, entry.CreatedAtUTC)
	}
}

func TestBuildTagsDeduplicatesAndNormalizes(t *testing.T) {
	got := buildTags("issue_description", "core-app", "Core App", "Bug", "In Progress")
	want := []string{"redmine", "issue-description", "core-app", "bug", "in-progress"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("buildTags() = %#v, want %#v", got, want)
	}
}

func TestParseCSVRecords(t *testing.T) {
	data := []byte("issue_id,subject\n123,Hello\n124,World\n")
	records, err := parseCSVRecords(data)
	if err != nil {
		t.Fatalf("parseCSVRecords() error = %v", err)
	}
	if len(records) != 2 || records[0]["issue_id"] != "123" || records[1]["subject"] != "World" {
		t.Fatalf("records = %#v", records)
	}
}

func TestIssueFilters(t *testing.T) {
	cfg := config{
		projectIdentifiers: stringList{"backend", "frontend"},
		issueIDs:           intList{1, 2, 3},
	}
	filters := issueFilters(cfg)
	got := []string{
		"p.identifier IN ('backend', 'frontend')",
		"i.id IN (1, 2, 3)",
	}
	if !reflect.DeepEqual(filters, got) {
		t.Fatalf("issueFilters() = %#v, want %#v", filters, got)
	}
}
