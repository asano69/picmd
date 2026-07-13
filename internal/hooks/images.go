// Package hooks wires PocketBase record hooks that are not handled by
// declarative collection rules alone.
package hooks

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/filesystem"
	"github.com/sirupsen/logrus"

	"github.com/asano69/picmd2/internal/media"
)

// compressor is the Compressor implementation used for all uploaded images.
var compressor media.Compressor = media.WebPCompressor{}

// RegisterImageCompression intercepts uploads to the "images" collection's
// "image" field, replacing the uploaded file with a resized/compressed
// version before the record is persisted. The "filename" and "filesize"
// fields are updated to reflect the compressed output.
//
// Note: each "images" record holds exactly one file. Multi-file selection
// on the frontend is handled by creating one record per file, so this
// hook only ever needs to deal with a single upload per call.
func RegisterImageCompression(app core.App) {
	app.OnRecordCreateRequest("images").BindFunc(func(e *core.RecordRequestEvent) error {
		files, err := e.FindUploadedFiles("image")
		if err != nil {
			return err
		}
		if len(files) == 0 {
			return e.Next()
		}

		original := files[0]

		r, err := original.Reader.Open()
		if err != nil {
			return fmt.Errorf("open uploaded file: %w", err)
		}
		defer r.Close()

		result, err := compressor.Compress(r, original.Size)
		if err != nil {
			return fmt.Errorf("compress image: %w", err)
		}

		name := replaceExt(original.OriginalName, result.Extension)

		compressed, err := filesystem.NewFileFromBytes(result.Data, name)
		if err != nil {
			return fmt.Errorf("wrap compressed file: %w", err)
		}

		e.Record.Set("image", compressed)
		e.Record.Set("filename", name)
		e.Record.Set("filesize", result.CompressedSize)

		logrus.WithFields(logrus.Fields{
			"name":       name,
			"original":   result.OriginalSize,
			"compressed": result.CompressedSize,
		}).Info("images: compressed upload")

		return e.Next()
	})
}

// RegisterUUIDAssignment stamps each new "images" record with a UUIDv7,
// used to build a stable public URL independent of PocketBase's own
// collection/record-id scheme (see the "/i/{uuid}" route in serve.go).
func RegisterUUIDAssignment(app core.App) {
	app.OnRecordCreateRequest("images").BindFunc(func(e *core.RecordRequestEvent) error {
		id, err := uuid.NewV7()
		if err != nil {
			return fmt.Errorf("generate uuid: %w", err)
		}
		e.Record.Set("uuid", id.String())
		return e.Next()
	})
}

// replaceExt returns name with its extension replaced by ext (which
// includes the leading dot, e.g. ".webp").
func replaceExt(name, ext string) string {
	base := strings.TrimSuffix(name, filepath.Ext(name))
	return base + ext
}
