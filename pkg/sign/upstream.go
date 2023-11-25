package sign

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"net/url"
	"strings"
)

const (
	UpstreamTypeFile      = "file"
	UpstreamTypeHypertext = "hypertext"
	UpstreamTypeUnknown   = "unknown"
)

type UpstreamConfig struct {
	ID  string `json:"id"`
	URI string `json:"uri"`
}

func (v *UpstreamConfig) GetType() string {
	protocol := strings.SplitN(v.URI, "://", 2)[0]
	switch protocol {
	case "file":
	case "files":
		return UpstreamTypeFile
	case "http":
	case "https":
		return UpstreamTypeHypertext
	}

	return UpstreamTypeUnknown
}

func (v *UpstreamConfig) GetRawURI() (string, url.Values) {
	uri := strings.SplitN(v.URI, "://", 2)[1]
	data := strings.SplitN(uri, "?", 2)
	qs, _ := url.ParseQuery(uri)

	return data[0], qs
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
