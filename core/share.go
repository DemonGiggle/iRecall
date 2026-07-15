package core

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/gigol/irecall/core/db"
)

const maxBundleUncompressedBytes = 100 << 20
const maxBundleEntries = 512

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

// PreviewQuoteBundle returns the v3 manifest shown before a binary export.
func (e *Engine) PreviewQuoteBundle(ctx context.Context, ids []int64) ([]byte, error) {
	env, _, err := e.bundleEnvelope(ctx, ids)
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(env, "", "  ")
}

// ExportQuoteBundle creates a portable .irecall ZIP archive for all new exports.
func (e *Engine) ExportQuoteBundle(ctx context.Context, ids []int64) ([]byte, error) {
	env, files, err := e.bundleEnvelope(ctx, ids)
	if err != nil {
		return nil, err
	}
	manifest, err := json.MarshalIndent(env, "", "  ")
	if err != nil {
		return nil, err
	}
	var out bytes.Buffer
	zw := zip.NewWriter(&out)
	w, err := zw.Create("manifest.json")
	if err != nil {
		return nil, err
	}
	if _, err := w.Write(manifest); err != nil {
		return nil, err
	}
	for archivePath, data := range files {
		w, err := zw.Create(archivePath)
		if err != nil {
			return nil, err
		}
		if _, err := w.Write(data); err != nil {
			return nil, err
		}
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func (e *Engine) bundleEnvelope(ctx context.Context, ids []int64) (SharedQuoteEnvelope, map[string][]byte, error) {
	_ = ctx
	if len(ids) == 0 {
		return SharedQuoteEnvelope{}, nil, fmt.Errorf("no quotes selected for export")
	}
	env := SharedQuoteEnvelope{SchemaVersion: BundleSchemaVersion, ExportedAt: time.Now().UTC()}
	files := map[string][]byte{}
	for _, id := range ids {
		q, err := e.loadQuote(id)
		if err != nil {
			return SharedQuoteEnvelope{}, nil, err
		}
		entry := SharedQuoteEntry{GlobalID: q.GlobalID, AuthorUserID: q.AuthorUserID, AuthorName: q.AuthorName,
			SourceUserID: q.SourceUserID, SourceName: q.SourceName, SourceBackend: q.SourceBackend,
			SourceNamespace: q.SourceNamespace, SourceEntityType: q.SourceEntityType, SourceEntityID: q.SourceEntityID,
			SourceLabel: q.SourceLabel, SourceURL: q.SourceURL, Version: q.Version, Content: q.Content,
			Tags: append([]string(nil), q.Tags...), CreatedAtUTC: q.CreatedAt.UTC(), UpdatedAtUTC: q.UpdatedAt.UTC()}
		rows, err := e.store.ListQuoteAttachments(q.ID)
		if err != nil {
			return SharedQuoteEnvelope{}, nil, err
		}
		for _, row := range rows {
			data, err := e.attachmentData(row)
			if err != nil {
				return SharedQuoteEnvelope{}, nil, err
			}
			archivePath := path.Join("attachments", row.ID+filepath.Ext(row.StoragePath))
			entry.Attachments = append(entry.Attachments, SharedAttachmentEntry{ID: row.ID, Filename: row.Filename,
				MediaType: row.MediaType, Size: row.Size, Width: row.Width, Height: row.Height,
				SHA256: row.SHA256, ArchivePath: archivePath})
			files[archivePath] = data
		}
		env.Quotes = append(env.Quotes, entry)
	}
	return env, files, nil
}

func (e *Engine) ImportSharedQuotes(ctx context.Context, payload []byte) (ImportResult, error) {
	if len(payload) >= 4 && bytes.Equal(payload[:4], []byte{'P', 'K', 3, 4}) {
		return e.importQuoteBundle(ctx, payload)
	}
	return e.importSharedQuoteJSON(ctx, payload)
}

func (e *Engine) importSharedQuoteJSON(ctx context.Context, payload []byte) (ImportResult, error) {
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

func (e *Engine) importQuoteBundle(ctx context.Context, payload []byte) (ImportResult, error) {
	zr, err := zip.NewReader(bytes.NewReader(payload), int64(len(payload)))
	if err != nil {
		return ImportResult{}, fmt.Errorf("open quote bundle: %w", err)
	}
	files := make(map[string]*zip.File, len(zr.File))
	if len(zr.File) > maxBundleEntries {
		return ImportResult{}, fmt.Errorf("bundle has too many entries")
	}
	var total uint64
	for _, f := range zr.File {
		clean := path.Clean(f.Name)
		if clean != f.Name || strings.HasPrefix(clean, "../") || path.IsAbs(clean) {
			return ImportResult{}, fmt.Errorf("unsafe bundle path %q", f.Name)
		}
		total += f.UncompressedSize64
		if total > maxBundleUncompressedBytes {
			return ImportResult{}, fmt.Errorf("bundle exceeds 100 MiB uncompressed limit")
		}
		if files[f.Name] != nil {
			return ImportResult{}, fmt.Errorf("bundle contains duplicate path %q", f.Name)
		}
		files[f.Name] = f
	}
	manifestFile := files["manifest.json"]
	if manifestFile == nil {
		return ImportResult{}, fmt.Errorf("bundle is missing manifest.json")
	}
	manifest, err := readZipFileLimited(manifestFile, 2<<20)
	if err != nil {
		return ImportResult{}, err
	}
	var env SharedQuoteEnvelope
	if err := json.Unmarshal(manifest, &env); err != nil {
		return ImportResult{}, fmt.Errorf("decode bundle manifest: %w", err)
	}
	if env.SchemaVersion != BundleSchemaVersion {
		return ImportResult{}, fmt.Errorf("unsupported bundle schema version: %d", env.SchemaVersion)
	}
	if len(env.Quotes) > 100 {
		return ImportResult{}, fmt.Errorf("bundle has too many quotes")
	}

	type importedImages struct {
		entry   SharedQuoteEntry
		images  []ImageInput
		replace bool
	}
	validated := make([]importedImages, 0, len(env.Quotes))
	actualBytes := int64(len(manifest))
	for _, entry := range env.Quotes {
		if err := validateSharedQuoteEntry(entry); err != nil {
			return ImportResult{}, err
		}
		if len(entry.Attachments) > MaxQuoteImages {
			return ImportResult{}, fmt.Errorf("shared quote %s has too many images", entry.GlobalID)
		}
		existing, lookupErr := e.store.GetQuoteByGlobalID(entry.GlobalID)
		replace := lookupErr == sql.ErrNoRows || (lookupErr == nil && entry.Version > existing.Version)
		if lookupErr != nil && lookupErr != sql.ErrNoRows {
			return ImportResult{}, lookupErr
		}
		item := importedImages{entry: entry, replace: replace}
		for _, attachment := range entry.Attachments {
			if attachment.ArchivePath == "" || path.Clean(attachment.ArchivePath) != attachment.ArchivePath || !strings.HasPrefix(attachment.ArchivePath, "attachments/") {
				return ImportResult{}, fmt.Errorf("invalid attachment path %q", attachment.ArchivePath)
			}
			zf := files[attachment.ArchivePath]
			if zf == nil {
				return ImportResult{}, fmt.Errorf("bundle is missing %s", attachment.ArchivePath)
			}
			data, err := readZipFileLimited(zf, MaxImageBytes)
			if err != nil {
				return ImportResult{}, err
			}
			sum := sha256.Sum256(data)
			if fmt.Sprintf("%x", sum[:]) != strings.ToLower(attachment.SHA256) {
				return ImportResult{}, fmt.Errorf("checksum mismatch for %s", attachment.Filename)
			}
			actualBytes += int64(len(data))
			if actualBytes > maxBundleUncompressedBytes {
				return ImportResult{}, fmt.Errorf("bundle exceeds 100 MiB uncompressed limit")
			}
			item.images = append(item.images, ImageInput{Filename: attachment.Filename, MediaType: attachment.MediaType, Data: data})
		}
		if _, err := prepareImages(item.images, 0); err != nil {
			return ImportResult{}, err
		}
		validated = append(validated, item)
	}

	legacyCompatible := env
	legacyCompatible.SchemaVersion = ShareSchemaVersion
	manifestForImport, _ := json.Marshal(legacyCompatible)
	result, err := e.importSharedQuoteJSON(ctx, manifestForImport)
	if err != nil {
		return ImportResult{}, err
	}
	for _, item := range validated {
		if !item.replace {
			continue
		}
		row, err := e.store.GetQuoteByGlobalID(item.entry.GlobalID)
		if err != nil {
			return ImportResult{}, err
		}
		if err := e.replaceImportedAttachments(row.ID, item.images); err != nil {
			return ImportResult{}, err
		}
	}
	return result, nil
}

func readZipFileLimited(f *zip.File, limit int64) ([]byte, error) {
	if f.UncompressedSize64 > uint64(limit) {
		return nil, fmt.Errorf("bundle entry %s exceeds size limit", f.Name)
	}
	r, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return io.ReadAll(io.LimitReader(r, limit+1))
}

func (e *Engine) replaceImportedAttachments(quoteID int64, images []ImageInput) error {
	current, err := e.store.ListQuoteAttachments(quoteID)
	if err != nil {
		return err
	}
	prepared, err := prepareImages(images, 0)
	if err != nil {
		return err
	}
	for _, row := range current {
		if err := e.store.DeleteQuoteAttachment(row.ID); err != nil {
			return err
		}
		_ = os.Remove(filepath.Join(e.attachmentRoot, filepath.Base(row.StoragePath)))
	}
	for _, image := range prepared {
		row, err := e.writePreparedImage(quoteID, image)
		if err != nil {
			return err
		}
		if err := e.store.InsertQuoteAttachment(row); err != nil {
			_ = os.Remove(filepath.Join(e.attachmentRoot, row.StoragePath))
			return err
		}
	}
	return nil
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
