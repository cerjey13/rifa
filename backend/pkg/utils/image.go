package utils

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"math"
	"net/http"

	"github.com/disintegration/imaging"
)

const (
	_MaxImageWidth      int = 800       // px
	_MinImageWidth      int = 320       // px (don’t go smaller than this)
	_JPEGQualityStart   int = 75        // start quality
	_JPEGQualityMin     int = 35        // don’t go below this
	_MaxScreenshotBytes int = 80 * 1024 // 80 KB
)

// CompressToJPG takes any image format (JPEG, PNG, GIF, WEBP, etc.)
// and compresses it into a normalized JPEG under _MaxScreenshotBytes.
func CompressToJPG(screenshot []byte) ([]byte, error) {
	if len(screenshot) == 0 {
		return nil, errors.New("no screenshot data")
	}

	mime := http.DetectContentType(screenshot)
	if mime == "application/octet-stream" {
		return nil, fmt.Errorf("unsupported or unrecognized image data")
	}

	img, format, err := image.Decode(bytes.NewReader(screenshot))
	if err != nil {
		return nil, fmt.Errorf("decode (%s): %w", mime, err)
	}

	w := img.Bounds().Dx()
	if w > _MaxImageWidth {
		img = imaging.Resize(img, _MaxImageWidth, 0, imaging.Lanczos)
	}

	if out, ok := tryQualities(img); ok {
		return out, nil
	}

	width := img.Bounds().Dx()
	for width > _MinImageWidth {
		// shrink by ~10% each iteration
		width = int(math.Max(float64(width)*0.9, float64(_MinImageWidth)))
		img = imaging.Resize(img, width, 0, imaging.Lanczos)

		if out, ok := tryQualities(img); ok {
			return out, nil
		}
	}

	return nil, fmt.Errorf(
		"unable to compress image under %d bytes (input format: %s)",
		_MaxScreenshotBytes,
		format,
	)
}

// tryQualities tries descending JPEG quality levels until the image fits the byte limit.
func tryQualities(img image.Image) ([]byte, bool) {
	for q := _JPEGQualityStart; q >= _JPEGQualityMin; q -= 5 {
		var buf bytes.Buffer
		if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: q}); err != nil {
			// encoding error — bail on this quality but keep trying smaller q
			continue
		}
		if buf.Len() <= _MaxScreenshotBytes {
			return buf.Bytes(), true
		}
	}
	return nil, false
}
