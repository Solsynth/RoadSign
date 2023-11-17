package configurator

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"strings"
)

const (
	UpstreamTypeFile      = "file"
	UpstreamTypeHypertext = "hypertext"
	UpstreamTypeUnknown   = "unknown"
)

type UpstreamConfig struct {
	Name string `json:"name"`
	URI  string `json:"uri"`
}

func (v *UpstreamConfig) GetType() string {
	protocol := strings.SplitN(v.URI, "://", 2)[0]
	switch protocol {
	case "file":
		return UpstreamTypeFile
	case "http":
	case "https":
		return UpstreamTypeHypertext
	}

	return UpstreamTypeUnknown
}

func (v *UpstreamConfig) MakeURI(ctx *fiber.Ctx) string {
	var queries []string
	for k, v := range ctx.Queries() {
		queries = append(queries, fmt.Sprintf("%s=%s", k, v))
	}

	path := string(ctx.Request().URI().Path())
	hash := string(ctx.Request().URI().Hash())

	return v.URI + path +
		lo.Ternary(len(queries) > 0, "?"+strings.Join(queries, "&"), "") +
		lo.Ternary(len(hash) > 0, "#"+hash, "")
}
