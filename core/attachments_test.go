package core

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/gigol/irecall/core/db"
)

func TestQuoteImageLifecycleAndBundleRoundTrip(t *testing.T) {
	ctx := context.Background()
	source := newAttachmentTestEngine(t, filepath.Join(t.TempDir(), "source"), "source-user")
	imageData := testPNG(t)

	created, err := source.AddQuoteWithImages(ctx, "diagram of a blue system", []ImageInput{{Filename: "diagram.png", Data: imageData}})
	if err != nil {
		t.Fatalf("AddQuoteWithImages() error = %v", err)
	}
	if len(created.Quote.Attachments) != 1 {
		t.Fatalf("attachments = %#v, want one", created.Quote.Attachments)
	}
	attachment := created.Quote.Attachments[0]
	if attachment.Width != 3 || attachment.Height != 2 || attachment.MediaType != "image/png" {
		t.Fatalf("attachment metadata = %#v", attachment)
	}
	_, stored, err := source.GetQuoteAttachmentData(ctx, attachment.ID)
	if err != nil {
		t.Fatalf("GetQuoteAttachmentData() error = %v", err)
	}
	if !bytes.Equal(stored, imageData) {
		t.Fatal("stored image differs from input")
	}

	bundle, err := source.ExportQuoteBundle(ctx, []int64{created.Quote.ID})
	if err != nil {
		t.Fatalf("ExportQuoteBundle() error = %v", err)
	}
	if !bytes.HasPrefix(bundle, []byte{'P', 'K', 3, 4}) {
		t.Fatal("bundle is not a ZIP archive")
	}

	destination := newAttachmentTestEngine(t, filepath.Join(t.TempDir(), "destination"), "destination-user")
	result, err := destination.ImportSharedQuotes(ctx, bundle)
	if err != nil {
		t.Fatalf("ImportSharedQuotes(bundle) error = %v", err)
	}
	if result.Inserted != 1 {
		t.Fatalf("import result = %#v, want one insert", result)
	}
	quotes, err := destination.ListQuotes(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(quotes) != 1 || len(quotes[0].Attachments) != 1 {
		t.Fatalf("imported quotes = %#v", quotes)
	}
	_, imported, err := destination.GetQuoteAttachmentData(ctx, quotes[0].Attachments[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(imported, imageData) {
		t.Fatal("bundle image differs after import")
	}

	updated, err := source.UpdateQuoteWithImages(ctx, created.Quote.ID, "text only now", nil, nil)
	if err != nil {
		t.Fatalf("UpdateQuoteWithImages(remove) error = %v", err)
	}
	if len(updated.Quote.Attachments) != 0 {
		t.Fatalf("attachments after removal = %#v", updated.Quote.Attachments)
	}
}

func TestQuoteImagesValidateTypeAndCount(t *testing.T) {
	engine := newAttachmentTestEngine(t, t.TempDir(), "user")
	if _, err := engine.AddQuoteWithImages(context.Background(), "bad image", []ImageInput{{Filename: "fake.png", Data: []byte("not an image")}}); err == nil {
		t.Fatal("unsupported image should fail")
	}
	images := make([]ImageInput, MaxQuoteImages+1)
	for i := range images {
		images[i] = ImageInput{Filename: "image.png", Data: testPNG(t)}
	}
	if _, err := engine.AddQuoteWithImages(context.Background(), "too many", images); err == nil {
		t.Fatal("too many images should fail")
	}
}

func newAttachmentTestEngine(t *testing.T, root, userID string) *Engine {
	t.Helper()
	if err := os.MkdirAll(root, 0o700); err != nil {
		t.Fatal(err)
	}
	store, err := db.Open(filepath.Join(root, "irecall.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = store.Close() })
	cfg := DefaultSettings()
	cfg.Debug.MockLLM = true
	engine := New(store, cfg)
	if err := engine.SetAttachmentRoot(filepath.Join(root, "attachments")); err != nil {
		t.Fatal(err)
	}
	profile := &UserProfile{UserID: userID, DisplayName: userID}
	if err := engine.SaveUserProfile(context.Background(), profile); err != nil {
		t.Fatal(err)
	}
	engine.UpdateUserProfile(profile)
	return engine
}

func testPNG(t *testing.T) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 3, 2))
	for y := 0; y < 2; y++ {
		for x := 0; x < 3; x++ {
			img.Set(x, y, color.RGBA{B: 255, A: 255})
		}
	}
	var out bytes.Buffer
	if err := png.Encode(&out, img); err != nil {
		t.Fatal(err)
	}
	return out.Bytes()
}
