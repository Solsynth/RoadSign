package sign

import (
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
)

func makeHypertextResponse(ctx *fiber.Ctx, upstream *UpstreamConfig) error {
	timeout := time.Duration(viper.GetInt64("performance.network_timeout")) * time.Millisecond
	return proxy.Do(ctx, upstream.MakeURI(ctx), &fasthttp.Client{
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
	})
}

func makeFileResponse(ctx *fiber.Ctx, upstream *UpstreamConfig) error {
	uri, queries := upstream.GetRawURI()
	root := http.Dir(uri)

	method := ctx.Method()

	// We only serve static assets for GET and HEAD methods
	if method != fiber.MethodGet && method != fiber.MethodHead {
		return ctx.Next()
	}

	// Strip prefix
	prefix := ctx.Route().Path
	path := strings.TrimPrefix(ctx.Path(), prefix)
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	// Add prefix
	if queries.Get("prefix") != "" {
		path = queries.Get("prefix") + path
	}

	if len(path) > 1 {
		path = utils.TrimRight(path, '/')
	}

	file, err := root.Open(path)
	if err != nil && errors.Is(err, fs.ErrNotExist) {
		if queries.Get("suffix") != "" {
			file, err = root.Open(path + queries.Get("suffix"))
		}
		if err != nil && queries.Get("fallback") != "" {
			file, err = root.Open(queries.Get("fallback"))
		}
	}

	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return ctx.Status(fiber.StatusNotFound).Next()
		}
		return fmt.Errorf("failed to open: %w", err)
	}

	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat: %w", err)
	}

	// Serve index if path is directory
	if stat.IsDir() {
		indexPath := utils.TrimRight(path, '/') + queries.Get("index")
		index, err := root.Open(indexPath)
		if err == nil {
			indexStat, err := index.Stat()
			if err == nil {
				file = index
				stat = indexStat
			}
		}
	}

	ctx.Status(fiber.StatusOK)

	modTime := stat.ModTime()
	contentLength := int(stat.Size())

	// Set Content-Type header
	if queries.Get("charset") == "" {
		ctx.Type(filepath.Ext(stat.Name()))
	} else {
		ctx.Type(filepath.Ext(stat.Name()), queries.Get("charset"))
	}

	// Set Last-Modified header
	if !modTime.IsZero() {
		ctx.Set(fiber.HeaderLastModified, modTime.UTC().Format(http.TimeFormat))
	}

	// Set Cache-Control header
	if method == fiber.MethodGet {
		maxAge, err := strconv.Atoi(queries.Get("maxAge"))
		if lo.Ternary(err != nil, maxAge, 0) > 0 {
			ctx.Set(fiber.HeaderCacheControl, "public, max-age="+queries.Get("maxAge"))
		}
		ctx.Response().SetBodyStream(file, contentLength)
		return nil
	}
	if method == fiber.MethodHead {
		ctx.Request().ResetBody()
		ctx.Response().SkipBody = true
		ctx.Response().Header.SetContentLength(contentLength)
		if err := file.Close(); err != nil {
			return fmt.Errorf("failed to close: %w", err)
		}
		return nil
	}

	return ctx.Next()
}
