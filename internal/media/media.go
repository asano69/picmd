// Package media provides image compression for uploaded files.
package media

import "io"

// Result holds the output of a compression operation.
type Result struct {
	Data           []byte
	Extension      string // e.g. ".webp"
	OriginalSize   int64
	CompressedSize int64
}

// Compressor converts and compresses an image from r.
// Implementations are free to choose format, quality, and resize strategy.
type Compressor interface {
	Compress(r io.Reader, originalSize int64) (*Result, error)
}
