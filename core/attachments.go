package core

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gigol/irecall/core/db"
	"github.com/google/uuid"
	_ "golang.org/x/image/webp"
)

const (
	MaxQuoteImages     = 5
	MaxImageBytes      = 10 << 20
	MaxQuoteImageBytes = 25 << 20
)

var supportedImageTypes = map[string]string{
	"image/png":  ".png",
	"image/jpeg": ".jpg",
	"image/webp": ".webp",
	"image/gif":  ".gif",
}

type preparedImage struct {
	input  ImageInput
	id     string
	ext    string
	width  int
	height int
	sha    string
}

func prepareImages(inputs []ImageInput, existingCount int) ([]preparedImage, error) {
	if existingCount+len(inputs) > MaxQuoteImages {
		return nil, fmt.Errorf("a quote can have at most %d images", MaxQuoteImages)
	}
	total := int64(0)
	out := make([]preparedImage, 0, len(inputs))
	for _, in := range inputs {
		size := int64(len(in.Data))
		if size == 0 {
			return nil, fmt.Errorf("image %q is empty", in.Filename)
		}
		if size > MaxImageBytes {
			return nil, fmt.Errorf("image %q exceeds the 10 MiB limit", in.Filename)
		}
		total += size
		if total > MaxQuoteImageBytes {
			return nil, fmt.Errorf("quote images exceed the 25 MiB total limit")
		}
		mediaType := http.DetectContentType(in.Data[:min(len(in.Data), 512)])
		ext, ok := supportedImageTypes[mediaType]
		if !ok {
			return nil, fmt.Errorf("image %q has unsupported media type %q", in.Filename, mediaType)
		}
		cfg, format, err := image.DecodeConfig(bytes.NewReader(in.Data))
		if err != nil {
			return nil, fmt.Errorf("decode image %q: %w", in.Filename, err)
		}
		if cfg.Width < 1 || cfg.Height < 1 || format == "" {
			return nil, fmt.Errorf("image %q has invalid dimensions", in.Filename)
		}
		name := strings.TrimSpace(filepath.Base(in.Filename))
		if name == "" || name == "." {
			name = "image" + ext
		}
		sum := sha256.Sum256(in.Data)
		in.Filename = name
		in.MediaType = mediaType
		out = append(out, preparedImage{
			input: in, id: uuid.NewString(), ext: ext,
			width: cfg.Width, height: cfg.Height, sha: hex.EncodeToString(sum[:]),
		})
	}
	return out, nil
}

func (e *Engine) SetAttachmentRoot(root string) error {
	root = strings.TrimSpace(root)
	if root == "" {
		return fmt.Errorf("attachment root is empty")
	}
	abs, err := filepath.Abs(root)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(abs, 0o700); err != nil {
		return fmt.Errorf("create attachment directory: %w", err)
	}
	e.attachmentRoot = abs
	return e.cleanupOrphanAttachments()
}

func (e *Engine) cleanupOrphanAttachments() error {
	referenced, err := e.store.ListAttachmentStoragePaths()
	if err != nil {
		return err
	}
	entries, err := os.ReadDir(e.attachmentRoot)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasPrefix(entry.Name(), ".staging-") || !referenced[entry.Name()] {
			if err := os.Remove(filepath.Join(e.attachmentRoot, entry.Name())); err != nil && !os.IsNotExist(err) {
				return err
			}
		}
	}
	return nil
}

func (e *Engine) writePreparedImage(quoteID int64, p preparedImage) (db.QuoteAttachmentRow, error) {
	if e.attachmentRoot == "" {
		return db.QuoteAttachmentRow{}, fmt.Errorf("attachment storage is not configured")
	}
	name := p.id + p.ext
	finalPath := filepath.Join(e.attachmentRoot, name)
	tmp, err := os.CreateTemp(e.attachmentRoot, ".staging-*")
	if err != nil {
		return db.QuoteAttachmentRow{}, err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)
	if err := tmp.Chmod(0o600); err != nil {
		tmp.Close()
		return db.QuoteAttachmentRow{}, err
	}
	if _, err := tmp.Write(p.input.Data); err != nil {
		tmp.Close()
		return db.QuoteAttachmentRow{}, err
	}
	if err := tmp.Close(); err != nil {
		return db.QuoteAttachmentRow{}, err
	}
	if err := os.Rename(tmpName, finalPath); err != nil {
		return db.QuoteAttachmentRow{}, err
	}
	return db.QuoteAttachmentRow{
		ID: p.id, QuoteID: quoteID, Filename: p.input.Filename, MediaType: p.input.MediaType,
		Size: int64(len(p.input.Data)), Width: p.width, Height: p.height, SHA256: p.sha,
		StoragePath: name, CreatedAt: time.Now().Unix(),
	}, nil
}

func attachmentFromRow(a db.QuoteAttachmentRow) QuoteAttachment {
	return QuoteAttachment{ID: a.ID, Filename: a.Filename, MediaType: a.MediaType, Size: a.Size,
		Width: a.Width, Height: a.Height, CreatedAt: time.Unix(a.CreatedAt, 0)}
}

func (e *Engine) attachmentData(row db.QuoteAttachmentRow) ([]byte, error) {
	if e.attachmentRoot == "" {
		return nil, fmt.Errorf("attachment storage is not configured")
	}
	return os.ReadFile(filepath.Join(e.attachmentRoot, filepath.Base(row.StoragePath)))
}

func (e *Engine) GetQuoteAttachmentData(ctx context.Context, id string) (QuoteAttachment, []byte, error) {
	_ = ctx
	row, err := e.store.GetQuoteAttachment(strings.TrimSpace(id))
	if err != nil {
		return QuoteAttachment{}, nil, err
	}
	data, err := e.attachmentData(row)
	return attachmentFromRow(row), data, err
}

func (e *Engine) enrichQuoteAttachments(quotes []Quote) error {
	for i := range quotes {
		rows, err := e.store.ListQuoteAttachments(quotes[i].ID)
		if err != nil {
			return err
		}
		quotes[i].Attachments = make([]QuoteAttachment, len(rows))
		for j, row := range rows {
			quotes[i].Attachments[j] = attachmentFromRow(row)
		}
	}
	return nil
}
