// Package hooks wires PocketBase record hooks that are not handled by
// declarative collection rules alone.
package hooks

import (
	"fmt"
	"path/filepath"
	"strings"

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
// Note: only the first uploaded file is compressed. Multi-image upload is
// still an open question (see README "複数画像のアップロードの方法を検討する"),
// so this intentionally stays single-file until that's decided.
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

// replaceExt returns name with its extension replaced by ext (which
// includes the leading dot, e.g. ".webp").
func replaceExt(name, ext string) string {
	base := strings.TrimSuffix(name, filepath.Ext(name))
	return base + ext
}
