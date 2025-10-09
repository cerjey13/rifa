package mymiddlewares

import (
	"bytes"
	"context"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"

	"rifa/backend/internal/repository"
	"rifa/backend/pkg/logx"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
)

type bodyCaptureCtx struct {
	middleContext huma.Context
	w             io.Writer
}

func (c bodyCaptureCtx) Operation() *huma.Operation {
	return c.middleContext.Operation()
}
func (c bodyCaptureCtx) Context() context.Context {
	return c.middleContext.Context()
}
func (c bodyCaptureCtx) TLS() *tls.ConnectionState {
	return c.middleContext.TLS()
}
func (c bodyCaptureCtx) Version() huma.ProtoVersion {
	return c.middleContext.Version()
}
func (c bodyCaptureCtx) Method() string {
	return c.middleContext.Method()
}
func (c bodyCaptureCtx) Host() string {
	return c.middleContext.Host()
}
func (c bodyCaptureCtx) RemoteAddr() string {
	return c.middleContext.RemoteAddr()
}
func (c bodyCaptureCtx) URL() url.URL {
	return c.middleContext.URL()
}
func (c bodyCaptureCtx) Param(name string) string {
	return c.middleContext.Param(name)
}
func (c bodyCaptureCtx) Query(name string) string {
	return c.middleContext.Query(name)
}
func (c bodyCaptureCtx) Header(name string) string {
	return c.middleContext.Header(name)
}
func (c bodyCaptureCtx) EachHeader(cb func(name, value string)) {
	c.middleContext.EachHeader(cb)
}
func (c bodyCaptureCtx) BodyReader() io.Reader {
	return c.middleContext.BodyReader()
}
func (c bodyCaptureCtx) GetMultipartForm() (*multipart.Form, error) {

	return c.middleContext.GetMultipartForm()
}
func (c bodyCaptureCtx) SetReadDeadline(t time.Time) error {
	return c.middleContext.SetReadDeadline(t)
}
func (c bodyCaptureCtx) SetStatus(code int) { c.middleContext.SetStatus(code) }
func (c bodyCaptureCtx) Status() int {
	return c.middleContext.Status()
}
func (c bodyCaptureCtx) SetHeader(name, value string) {
	c.middleContext.SetHeader(name, value)
}
func (c bodyCaptureCtx) AppendHeader(name, value string) {
	c.middleContext.AppendHeader(name, value)
}
func (c bodyCaptureCtx) BodyWriter() io.Writer {
	return c.w
}

const maxFileBytes int64 = 1 * 1024 * 1024

func IdempotencyMiddleware(
	api huma.API,
	repo repository.IdempotencyRepository,
	logger logx.Logger,
) func(ctx huma.Context, next func(ctx huma.Context)) {
	return func(ctx huma.Context, next func(ctx huma.Context)) {
		key := ctx.Header("Idempotency-Key")
		if key == "" {
			next(ctx)
			return
		}

		contentType := ctx.Header("Content-Type")
		if !strings.HasPrefix(contentType, "multipart/form-data") {
			logger.Warn(
				ctx.Context(),
				"Invalid Content-Type for idempotent endpoint",
				"type",
				contentType,
			)
			_ = huma.WriteErr(
				api,
				ctx,
				http.StatusBadRequest,
				"Invalid Content-Type, expected multipart/form-data",
				errors.New("invalid content-type"),
			)
			return
		}

		boundary := boundaryFromContentType(contentType)
		if boundary == "" {
			logger.Warn(ctx.Context(), "Missing multipart boundary")
			_ = huma.WriteErr(
				api,
				ctx,
				http.StatusBadRequest,
				"Missing multipart boundary",
				errors.New("missing boundary"),
			)
			return
		}

		r, _ := humachi.Unwrap(ctx)
		bodyBytes, err := io.ReadAll(io.LimitReader(r.Body, maxFileBytes))
		if err != nil {
			logger.Error(
				ctx.Context(),
				"Failed to read request body",
				"error",
				err,
			)
			_ = huma.WriteErr(
				api,
				ctx,
				http.StatusBadRequest,
				"Invalid request body",
				err,
			)
			return
		}

		formReader := multipart.NewReader(bytes.NewReader(bodyBytes), boundary)
		formValues := map[string]string{}

		for {
			part, err := formReader.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				logger.Error(
					ctx.Context(),
					"Failed to read multipart body",
					"error",
					err,
				)
				_ = huma.WriteErr(
					api,
					ctx,
					http.StatusBadRequest,
					"Invalid multipart body",
					err,
				)
				return
			}

			name := part.FormName()
			if part.FileName() != "" {
				closeErr := part.Close()
				if closeErr != nil {
					logger.Warn(
						ctx.Context(),
						"Failed to close file part",
						"field",
						name,
						"error",
						closeErr,
					)
				}
				continue
			}

			data, readErr := io.ReadAll(part)
			if closeErr := part.Close(); closeErr != nil {
				logger.Warn(
					ctx.Context(),
					"Failed to close form field part",
					"field",
					name,
					"error",
					closeErr,
				)
			}
			if readErr != nil {
				logger.Error(
					ctx.Context(),
					"Failed to read form field",
					"field",
					name,
					"error",
					readErr,
				)
				_ = huma.WriteErr(
					api,
					ctx,
					http.StatusBadRequest,
					"Invalid form field",
					readErr,
				)
				return
			}

			formValues[part.FormName()] = string(data)
		}

		cleanBody, err := json.Marshal(formValues)
		if err != nil {
			logger.Error(
				ctx.Context(),
				"Failed to re-marshal cleaned JSON",
				"error",
				err,
			)
			_ = huma.WriteErr(
				api,
				ctx,
				http.StatusInternalServerError,
				"Failed to hash request",
				err,
			)
			return
		}

		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		sum := sha256.Sum256(cleanBody)
		bodyHash := hex.EncodeToString(sum[:])

		rec, err := repo.Get(ctx.Context(), key)
		if err == nil {
			if rec.BodyHash != bodyHash {
				logger.Warn(
					ctx.Context(),
					"Idempotency key reused with different payload",
					"key",
					key,
				)
				_ = huma.WriteErr(
					api,
					ctx,
					http.StatusConflict,
					"Idempotency key reused with different payload",
					errors.New("payload mismatch"),
				)
				return
			}
			logger.Info(ctx.Context(), "Idempotent replay detected", "key", key)
			ctx.SetStatus(rec.StatusCode)
			_, writeErr := ctx.BodyWriter().Write(rec.ResponseBody)
			if writeErr != nil {
				logger.Error(
					ctx.Context(),
					"Failed to write cached response body",
					"key",
					key,
					"error",
					writeErr,
				)
			}
			return
		}

		buf := new(bytes.Buffer)
		mw := io.MultiWriter(ctx.BodyWriter(), buf)
		captureCtx := bodyCaptureCtx{ctx, mw}

		next(captureCtx)

		status := captureCtx.Status()
		saveErr := repo.Save(
			captureCtx.Context(),
			repository.IdempotencyRecord{
				Key:          key,
				BodyHash:     bodyHash,
				StatusCode:   status,
				ResponseBody: buf.Bytes(),
			},
		)
		if saveErr != nil {
			logger.Error(
				ctx.Context(),
				"Failed to persist idempotency record",
				"key",
				key,
				"error",
				saveErr,
			)
		} else {
			logger.Info(
				ctx.Context(),
				"Stored new idempotency record",
				"key",
				key,
				"status",
				status,
				"bytes",
				buf.Len(),
			)
		}
	}
}

func boundaryFromContentType(contentType string) string {
	_, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return ""
	}
	return params["boundary"]
}
