package navi

import (
	"fmt"
	"net/url"
	"strings"

	"code.smartsheep.studio/goatworks/roadsign/pkg/navi/transformers"
	"code.smartsheep.studio/goatworks/roadsign/pkg/warden"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

type Region struct {
	ID           string               `json:"id" toml:"id"`
	Disabled     bool                 `json:"disabled" toml:"disabled"`
	Locations    []Location           `json:"locations" toml:"locations"`
	Applications []warden.Application `json:"applications" toml:"applications"`
}

type Location struct {
	ID           string              `json:"id" toml:"id"`
	Host         []string            `json:"host" toml:"host"`
	Path         []string            `json:"path" toml:"path"`
	Queries      map[string]string   `json:"queries" toml:"queries"`
	Headers      map[string][]string `json:"headers" toml:"headers"`
	Destinations []Destination       `json:"destinations" toml:"destinations"`
}

type DestinationType = int8

const (
	DestinationHypertext = DestinationType(iota)
	DestinationStaticFile
	DestinationUnknown
)

type Destination struct {
	ID           string                           `json:"id" toml:"id"`
	Uri          string                           `json:"uri" toml:"uri"`
	Transformers []transformers.TransformerConfig `json:"transformers" toml:"transformers"`
}

func (v *Destination) GetType() DestinationType {
	protocol := strings.SplitN(v.Uri, "://", 2)[0]
	switch protocol {
	case "http", "https":
		return DestinationHypertext
	case "file", "files":
		return DestinationStaticFile
	}
	return DestinationUnknown
}

func (v *Destination) GetRawUri() (string, url.Values) {
	uri := strings.SplitN(v.Uri, "://", 2)[1]
	data := strings.SplitN(uri, "?", 2)
	data = append(data, " ") // Make data array least have two element
	qs, _ := url.ParseQuery(data[0])

	return data[0], qs
}

func (v *Destination) MakeUri(ctx *fiber.Ctx) string {
	var queries []string
	for k, v := range ctx.Queries() {
		parsed, _ := url.QueryUnescape(v)
		value := url.QueryEscape(parsed)
		queries = append(queries, fmt.Sprintf("%s=%s", k, value))
	}

	path := string(ctx.Request().URI().Path())
	hash := string(ctx.Request().URI().Hash())

	return v.Uri + path +
		lo.Ternary(len(queries) > 0, "?"+strings.Join(queries, "&"), "") +
		lo.Ternary(len(hash) > 0, "#"+hash, "")
}

func (v *Destination) MakeWebsocketUri(ctx *fiber.Ctx) string {
	return strings.Replace(v.MakeUri(ctx), "http", "ws", 1)
}
