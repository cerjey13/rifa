package utils

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	_ "image/png"

	"github.com/disintegration/imaging"
)

const (
	_MaxImageWidth      int = 800        // px
	_JPEGQuality        int = 70         // 0-100
	_MaxScreenshotBytes int = 100 * 1024 // 100 KB
)

func CompressToJPG(screenshot []byte) ([]byte, error) {
	if len(screenshot) == 0 {
		return nil, errors.New("no screenshot data")
	}

	img, _, err := image.Decode(bytes.NewReader(screenshot))
	if err != nil {
		return nil, err
	}

	// Only resize if wider than MaxImageWidth
	if img.Bounds().Dx() > _MaxImageWidth {
		img = imaging.Resize(img, _MaxImageWidth, 0, imaging.Lanczos)
	}

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: _JPEGQuality})
	if err != nil {
		return nil, err
	}

	compressedScreenshot := buf.Bytes()
	if len(compressedScreenshot) > _MaxScreenshotBytes {
		return nil, err
	}
	return compressedScreenshot, nil
}
