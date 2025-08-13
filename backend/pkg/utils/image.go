package utils

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"math"

	"github.com/disintegration/imaging"
)

const (
	_MaxImageWidth      int = 800       // px
	_MinImageWidth      int = 320       // px (don’t go smaller than this)
	_JPEGQualityStart   int = 75        // start quality
	_JPEGQualityMin     int = 35        // don’t go below this
	_MaxScreenshotBytes int = 80 * 1024 // 80 KB
)

func CompressToJPG(screenshot []byte) ([]byte, error) {
	if len(screenshot) == 0 {
		return nil, errors.New("no screenshot data")
	}

	img, _, err := image.Decode(bytes.NewReader(screenshot))
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
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

	return nil, fmt.Errorf("unable to compress image under %d bytes", _MaxScreenshotBytes)
}

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
