package utils

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"math/rand/v2"
	"testing"
)

// --- helpers ---

func mustPNGEncode(t *testing.T, img image.Image) []byte {
	t.Helper()
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("png.Encode: %v", err)
	}
	return buf.Bytes()
}

func mustJPEGEncode(t *testing.T, img image.Image, q int) []byte {
	t.Helper()
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: q}); err != nil {
		t.Fatalf("jpeg.Encode: %v", err)
	}
	return buf.Bytes()
}

func decodeJPG(t *testing.T, jpg []byte) image.Image {
	t.Helper()
	img, _, err := image.Decode(bytes.NewReader(jpg))
	if err != nil {
		t.Fatalf("decodeJPG: %v", err)
	}
	return img
}

func mkSolid(w, h int, c color.Color) *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := range h {
		for x := 0; x < w; x++ {
			img.Set(x, y, c)
		}
	}
	return img
}

func mkGradient(w, h int) *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := range h {
		for x := range w {
			img.SetNRGBA(x, y, color.NRGBA{
				R: uint8((x * 255) / (w + 1)),
				G: uint8((y * 255) / (h + 1)),
				B: uint8(((x + y) * 255) / (w + h + 2)),
				A: 255,
			})
		}
	}
	return img
}

func mkNoise(w, h int, seed uint64) *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	r := rand.New(rand.NewPCG(seed, seed^0x9e3779b97f4a7c15))
	for y := range h {
		for x := range w {
			img.SetNRGBA(x, y, color.NRGBA{
				R: uint8(r.IntN(256)),
				G: uint8(r.IntN(256)),
				B: uint8(r.IntN(256)),
				A: 255,
			})
		}
	}
	return img
}

func isJPEGMagic(b []byte) bool {
	return len(b) >= 2 && b[0] == 0xFF && b[1] == 0xD8
}

// min avoid slice bounds panic in error messages
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// --- tests ---
func TestCompressToJPG_Errors(t *testing.T) {
	tests := []struct {
		name string
		in   []byte
	}{
		{"empty_input", nil},
		{"undecodable_bytes", []byte("not an image")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := CompressToJPG(tt.in)
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if out != nil {
				t.Fatalf("expected nil output on error, got %d bytes", len(out))
			}
		})
	}
}

func TestCompressToJPG_PNGCases(t *testing.T) {
	type exp struct {
		maxWidth         int  // final width must be <= this (0 = ignore)
		minWidth         int  // final width must be >= this (0 = ignore)
		noUpscaleFromW   int  // ensure final width <= original width (0 = ignore)
		enforceSizeLimit bool // ensure <= _MaxScreenshotBytes
	}

	tests := []struct {
		name       string
		srcW, srcH int
		exp        exp
	}{
		{
			name: "small_png_no_upscale",
			srcW: 200, srcH: 200,
			exp: exp{
				noUpscaleFromW:   200,
				minWidth:         _MinImageWidth,
				enforceSizeLimit: true,
			},
		},
		{
			name: "large_png_resizes_and_compresses",
			srcW: 2400, srcH: 1600,
			exp: exp{
				maxWidth:         _MaxImageWidth,
				minWidth:         _MinImageWidth,
				enforceSizeLimit: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := mkGradient(tt.srcW, tt.srcH)
			pngBytes := mustPNGEncode(t, src)

			jpg, err := CompressToJPG(pngBytes)
			if err != nil {
				t.Fatalf("CompressToJPG error: %v", err)
			}
			if !isJPEGMagic(jpg) {
				n := min(4, len(jpg))
				t.Fatalf(
					"output is not JPEG (missing SOI), first bytes: %v",
					jpg[:n],
				)
			}
			if tt.exp.enforceSizeLimit && len(jpg) > _MaxScreenshotBytes {
				t.Fatalf(
					"compressed size = %d, want <= %d",
					len(jpg),
					_MaxScreenshotBytes,
				)
			}

			out := decodeJPG(t, jpg)
			w := out.Bounds().Dx()

			// no-upscale check
			if tt.exp.noUpscaleFromW > 0 && w > tt.exp.noUpscaleFromW {
				t.Fatalf(
					"width upscaled: got %d, want <= %d",
					w,
					tt.exp.noUpscaleFromW,
				)
			}
			// max/min width checks
			if tt.exp.maxWidth > 0 && w > tt.exp.maxWidth {
				t.Fatalf("expected width <= %d, got %d", tt.exp.maxWidth, w)
			}
			if tt.exp.minWidth > 0 && w < tt.exp.minWidth &&
				tt.srcW >= _MinImageWidth {
				t.Fatalf("expected width >= %d, got %d", tt.exp.minWidth, w)
			}
		})
	}
}

func TestCompressToJPG_AlreadyJPEGInputs(t *testing.T) {
	type exp struct {
		maxWidth         int  // final width must be <= this (0 = ignore)
		minWidth         int  // final width must be >= this (0 = ignore)
		noUpscaleFromW   int  // final width must be <= this (0 = ignore)
		enforceSizeLimit bool // output must be <= _MaxScreenshotBytes
	}

	tests := []struct {
		name       string
		srcW, srcH int
		quality    int
		exp        exp
	}{
		{
			name: "small_jpeg_no_upscale_q90",
			srcW: 320, srcH: 200, quality: 90,
			exp: exp{
				noUpscaleFromW:   320,
				minWidth:         _MinImageWidth,
				enforceSizeLimit: true,
			},
		},
		{
			name: "medium_jpeg_may_or_may_not_resize_q80",
			srcW: 800, srcH: 512, quality: 80,
			exp: exp{
				maxWidth:         _MaxImageWidth,
				minWidth:         _MinImageWidth,
				noUpscaleFromW:   800,
				enforceSizeLimit: true,
			},
		},
		{
			name: "large_jpeg_resizes_q90",
			srcW: 2048, srcH: 1024, quality: 90,
			exp: exp{
				maxWidth:         _MaxImageWidth,
				minWidth:         _MinImageWidth,
				enforceSizeLimit: true,
			},
		},
		{
			name: "very_tall_jpeg_resizes_q75",
			srcW: 900, srcH: 2400, quality: 75,
			exp: exp{
				maxWidth:         _MaxImageWidth,
				minWidth:         _MinImageWidth,
				enforceSizeLimit: true,
			},
		},
		{
			name: "very_wide_jpeg_resizes_q60",
			srcW: 3000, srcH: 400, quality: 60,
			exp: exp{
				maxWidth:         _MaxImageWidth,
				minWidth:         _MinImageWidth,
				enforceSizeLimit: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := mkSolid(
				tt.srcW,
				tt.srcH,
				color.NRGBA{R: 200, G: 50, B: 50, A: 255},
			)
			jpegSrc := mustJPEGEncode(t, src, tt.quality)

			jpg, err := CompressToJPG(jpegSrc)
			if err != nil {
				t.Fatalf("CompressToJPG error: %v", err)
			}
			if !isJPEGMagic(jpg) {
				t.Fatalf("output is not JPEG")
			}
			if tt.exp.enforceSizeLimit && len(jpg) > _MaxScreenshotBytes {
				t.Fatalf(
					"compressed size = %d, want <= %d",
					len(jpg),
					_MaxScreenshotBytes,
				)
			}

			out := decodeJPG(t, jpg)
			w := out.Bounds().Dx()

			// no-upscale check
			if tt.exp.noUpscaleFromW > 0 && w > tt.exp.noUpscaleFromW {
				t.Fatalf(
					"width upscaled: got %d, want <= %d",
					w,
					tt.exp.noUpscaleFromW,
				)
			}
			// width caps
			if tt.exp.maxWidth > 0 && w > tt.exp.maxWidth {
				t.Fatalf("expected width <= %d, got %d", tt.exp.maxWidth, w)
			}
			if tt.exp.minWidth > 0 && w < tt.exp.minWidth &&
				tt.srcW >= _MinImageWidth {
				t.Fatalf("expected width >= %d, got %d", tt.exp.minWidth, w)
			}
		})
	}
}

// --- benchmarks ---

func BenchmarkCompressToJPG_Black_512(b *testing.B) {
	src := mkSolid(512, 512, color.NRGBA{0, 0, 0, 255})
	in := mustPNGEncodeForBench(src, b)
	benchCompress(b, in)
}

func BenchmarkCompressToJPG_Colorful_1600x1200(b *testing.B) {
	src := mkGradient(1600, 1200)
	in := mustPNGEncodeForBench(src, b)
	benchCompress(b, in)
}

func BenchmarkCompressToJPG_Noisy_800x4000(b *testing.B) {
	src := mkNoise(800, 4000, 42)
	in := mustPNGEncodeForBench(src, b)
	benchCompress(b, in)
}

func mustPNGEncodeForBench(img image.Image, b *testing.B) []byte {
	b.Helper()
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		b.Fatalf("png.Encode: %v", err)
	}
	return buf.Bytes()
}

func benchCompress(b *testing.B, in []byte) {
	b.ReportAllocs()
	// Warm-up once (optional)
	if _, err := CompressToJPG(in); err != nil {
		b.Fatalf("warmup failed: %v", err)
	}
	b.ResetTimer()
	for b.Loop() {
		out, err := CompressToJPG(in)
		if err != nil {
			b.Fatalf("CompressToJPG: %v", err)
		}
		if len(out) == 0 {
			b.Fatalf("empty output")
		}
	}
}
