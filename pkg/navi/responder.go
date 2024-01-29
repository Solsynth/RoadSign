package navi

import (
	"errors"
	"fmt"
	"github.com/fasthttp/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/valyala/fasthttp"
	"io/fs"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func makeUnifiedResponse(c *fiber.Ctx, dest *Destination) error {
	if websocket.FastHTTPIsWebSocketUpgrade(c.Context()) {
		// Handle websocket
		return makeWebsocketResponse(c, dest)
	} else {
		// TODO Impl SSE with https://github.com/gofiber/recipes/blob/master/sse/main.go
		// Handle normal http request
		return makeHypertextResponse(c, dest)
	}
}

func makeHypertextResponse(c *fiber.Ctx, dest *Destination) error {
	_, queries := dest.GetRawUri()
	raw := lo.Ternary(len(queries.Get("timeout")) > 0, queries.Get("timeout"), "5000")
	num, err := strconv.Atoi(raw)
	if err != nil {
		num = 5000
	}

	limit := time.Duration(num) * time.Millisecond
	uri := dest.MakeUri(c)
	return proxy.Do(c, uri, &fasthttp.Client{
		ReadTimeout:  limit,
		WriteTimeout: limit,
	})
}

var wsUpgrader = websocket.FastHTTPUpgrader{}

func makeWebsocketResponse(c *fiber.Ctx, dest *Destination) error {
	uri := dest.MakeWebsocketUri(c)

	// Upgrade connection
	return wsUpgrader.Upgrade(c.Context(), func(conn *websocket.Conn) {
		// Dial the destination
		remote, _, err := websocket.DefaultDialer.Dial(uri, nil)
		if err != nil {
			return
		}
		defer remote.Close()

		// Read messages from remote
		disconnect := make(chan struct{})
		signal := make(chan struct {
			head int
			data []byte
		})
		go func() {
			defer close(disconnect)
			for {
				mode, message, err := remote.ReadMessage()
				if err != nil {
					log.Warn().Err(err).Msg("An error occurred during the websocket proxying...")
					return
				} else {
					signal <- struct {
						head int
						data []byte
					}{head: mode, data: message}
				}
			}
		}()

		// Relay the destination websocket to client
		for {
			select {
			case <-disconnect:
			case val := <-signal:
				if err := conn.WriteMessage(val.head, val.data); err != nil {
					return
				}
			default:
				if head, data, err := conn.ReadMessage(); err != nil {
					return
				} else {
					remote.WriteMessage(head, data)
				}
			}
		}
	})
}

func makeFileResponse(c *fiber.Ctx, dest *Destination) error {
	uri, queries := dest.GetRawUri()
	root := http.Dir(uri)

	method := c.Method()

	// We only serve static assets for GET and HEAD methods
	if method != fiber.MethodGet && method != fiber.MethodHead {
		return c.Next()
	}

	// Strip prefix
	prefix := c.Route().Path
	path := strings.TrimPrefix(c.Path(), prefix)
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
			return fiber.ErrNotFound
		}
		return fmt.Errorf("failed to open: %w", err)
	}

	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat: %w", err)
	}

	// Serve index if path is directory
	if stat.IsDir() {
		indexFile := lo.Ternary(len(queries.Get("index")) > 0, queries.Get("index"), "index.html")
		indexPath := utils.TrimRight(path, '/') + indexFile
		index, err := root.Open(indexPath)
		if err == nil {
			indexStat, err := index.Stat()
			if err == nil {
				file = index
				stat = indexStat
			}
		}
	}

	c.Status(fiber.StatusOK)

	modTime := stat.ModTime()
	contentLength := int(stat.Size())

	// Set Content-Type header
	if queries.Get("charset") == "" {
		c.Type(filepath.Ext(stat.Name()))
	} else {
		c.Type(filepath.Ext(stat.Name()), queries.Get("charset"))
	}

	// Set Last-Modified header
	if !modTime.IsZero() {
		c.Set(fiber.HeaderLastModified, modTime.UTC().Format(http.TimeFormat))
	}

	if method == fiber.MethodGet {
		maxAge, err := strconv.Atoi(queries.Get("maxAge"))
		if lo.Ternary(err != nil, maxAge, 0) > 0 {
			c.Set(fiber.HeaderCacheControl, "public, max-age="+queries.Get("maxAge"))
		}
		c.Response().SetBodyStream(file, contentLength)
		return nil
	}
	if method == fiber.MethodHead {
		c.Request().ResetBody()
		c.Response().SkipBody = true
		c.Response().Header.SetContentLength(contentLength)
		if err := file.Close(); err != nil {
			return fmt.Errorf("failed to close: %w", err)
		}
		return nil
	}

	return fiber.ErrNotFound
}
