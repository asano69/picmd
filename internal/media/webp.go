package media

import (
	"bytes"
	"fmt"
	"github.com/chai2010/webp"
	"golang.org/x/image/draw"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"image/png"
	"io"
)

const (
	maxWidth   = 1920
	maxHeight  = 1080
	maxBytes   = 500 * 1024 // 500 KB
	qualityMin = 10.0
	qualityMax = 85.0
)

// WebPCompressor resizes the image to fit within maxWidth x maxHeight,
// then encodes it as WebP, targeting maxBytes via binary search on quality.
type WebPCompressor struct{}

func (WebPCompressor) Compress(r io.Reader, originalSize int64) (*Result, error) {
	src, format, err := image.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	src = resize(src)

	if format == "png" && hasAlpha(src) {
		// Keep PNGs with transparency as PNG: chai2010/webp cannot reliably encode transparency.
		var buf bytes.Buffer
		if err := png.Encode(&buf, src); err != nil {
			return nil, fmt.Errorf("png encode: %w", err)
		}
		data := buf.Bytes()
		return &Result{
			Data:           data,
			Extension:      ".png",
			OriginalSize:   originalSize,
			CompressedSize: int64(len(data)),
		}, nil
	}

	// Fall through to WebP encoding for non-transparent PNGs and all other formats.
	data, err := encodeAtTargetSize(src)
	if err != nil {
		return nil, err
	}
	return &Result{
		Data:           data,
		Extension:      ".webp",
		OriginalSize:   originalSize,
		CompressedSize: int64(len(data)),
	}, nil
}

// hasAlpha reports whether img contains any pixel with non-opaque alpha.
func hasAlpha(img image.Image) bool {
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a < 0xffff {
				return true
			}
		}
	}
	return false
}

// toNRGBA returns the image as *image.NRGBA.
// If the image is already *image.NRGBA it is returned as-is.
func toNRGBA(img image.Image) *image.NRGBA {
	if v, ok := img.(*image.NRGBA); ok {
		return v
	}
	bounds := img.Bounds()
	dst := image.NewNRGBA(bounds)
	// Use draw.Over instead of draw.Src to handle paletted images correctly.
	draw.Draw(dst, bounds, img, bounds.Min, draw.Over)
	return dst
}

// resize scales img down so it fits within maxWidth x maxHeight, preserving aspect ratio.
// Images already within bounds are returned unchanged.
func resize(img image.Image) image.Image {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	if w <= maxWidth && h <= maxHeight {
		return img
	}

	scale := 1.0
	if float64(w)/float64(maxWidth) > scale {
		scale = float64(w) / float64(maxWidth)
	}
	if float64(h)/float64(maxHeight) > scale {
		scale = float64(h) / float64(maxHeight)
	}

	newW := int(float64(w) / scale)
	newH := int(float64(h) / scale)

	// Use NRGBA (non-premultiplied) to preserve semi-transparent edges correctly.
	dst := image.NewNRGBA(image.Rect(0, 0, newW, newH))
	draw.BiLinear.Scale(dst, dst.Bounds(), img, bounds, draw.Src, nil)
	return dst
}

// encodeAtTargetSize encodes img as lossy WebP using binary search over quality
// to find the highest quality that still fits within maxBytes.
func encodeAtTargetSize(img image.Image) ([]byte, error) {
	lo, hi := qualityMin, qualityMax

	var best []byte
	for hi-lo > 1 {
		mid := (lo + hi) / 2
		data, err := encodeWebP(img, float32(mid), false)
		if err != nil {
			return nil, err
		}
		if len(data) <= maxBytes {
			best = data
			lo = mid // try higher quality
		} else {
			hi = mid // too large, lower quality
		}
	}

	// If no quality produced a small-enough result, use the minimum.
	if best == nil {
		var err error
		best, err = encodeWebP(img, qualityMin, false)
		if err != nil {
			return nil, err
		}
	}
	return best, nil
}

func encodeWebP(img image.Image, quality float32, lossless bool) ([]byte, error) {
	var buf bytes.Buffer
	if err := webp.Encode(&buf, img, &webp.Options{Lossless: lossless, Quality: quality}); err != nil {
		return nil, fmt.Errorf("webp encode: %w", err)
	}
	return buf.Bytes(), nil
}
