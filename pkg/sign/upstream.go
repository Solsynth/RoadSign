package sign

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

const (
	UpstreamTypeFile      = "file"
	UpstreamTypeHypertext = "hypertext"
	UpstreamTypeUnknown   = "unknown"
)

type UpstreamInstance struct {
	ID  string `json:"id" yaml:"id"`
	URI string `json:"uri" yaml:"uri"`
}

func (v *UpstreamInstance) GetType() string {
	protocol := strings.SplitN(v.URI, "://", 2)[0]
	switch protocol {
	case "file", "files":
		return UpstreamTypeFile
	case "http", "https":
		return UpstreamTypeHypertext
	}

	return UpstreamTypeUnknown
}

func (v *UpstreamInstance) GetRawURI() (string, url.Values) {
	uri := strings.SplitN(v.URI, "://", 2)[1]
	data := strings.SplitN(uri, "?", 2)
	data = append(data, " ") // Make data array least have two element
	qs, _ := url.ParseQuery(data[0])

	return data[0], qs
}

func (v *UpstreamInstance) MakeURI(ctx *fiber.Ctx) string {
	var queries []string
	for k, v := range ctx.Queries() {
		parsed, _ := url.QueryUnescape(v)
		value := url.QueryEscape(parsed)
		queries = append(queries, fmt.Sprintf("%s=%s", k, value))
	}

	path := string(ctx.Request().URI().Path())
	hash := string(ctx.Request().URI().Hash())

	return v.URI + path +
		lo.Ternary(len(queries) > 0, "?"+strings.Join(queries, "&"), "") +
		lo.Ternary(len(hash) > 0, "#"+hash, "")
}
