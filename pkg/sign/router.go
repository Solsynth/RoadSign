package sign

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type AppConfig struct {
	Sites []SiteConfig `json:"sites"`
}

func (v *AppConfig) Forward(ctx *fiber.Ctx, site SiteConfig) error {
	if len(site.Upstreams) == 0 {
		return errors.New("invalid configuration")
	}

	idx := rand.Intn(len(site.Upstreams))
	upstream := &site.Upstreams[idx]

	switch upstream.GetType() {
	case UpstreamTypeHypertext:
		timeout := time.Duration(viper.GetInt64("performance.network_timeout")) * time.Millisecond
		return proxy.Do(ctx, upstream.MakeURI(ctx), &fasthttp.Client{
			ReadTimeout:  timeout,
			WriteTimeout: timeout,
		})
	case UpstreamTypeFile:
		uri, queries := upstream.GetRawURI()
		fs := filesystem.New(filesystem.Config{
			Root:               http.Dir(uri),
			ContentTypeCharset: queries.Get("charset"),
			Index: func() string {
				val := queries.Get("index")
				return lo.Ternary(len(val) > 0, val, "index.html")
			}(),
			NotFoundFile: func() string {
				val := queries.Get("fallback")
				return lo.Ternary(len(val) > 0, val, "404.html")
			}(),
			Browse: func() bool {
				browse, err := strconv.ParseBool(queries.Get("browse"))
				return lo.Ternary(err == nil, browse, false)
			}(),
			MaxAge: func() int {
				age, err := strconv.Atoi(queries.Get("maxAge"))
				return lo.Ternary(err == nil, age, 3600)
			}(),
		})
		return fs(ctx)
	default:
		return fiber.ErrBadGateway
	}
}

type SiteConfig struct {
	ID string `json:"id"`

	Name         string                     `json:"name"`
	Rules        []RouterRuleConfig         `json:"rules"`
	Transformers []RequestTransformerConfig `json:"transformers"`
	Upstreams    []UpstreamConfig           `json:"upstreams"`
}

type RouterRuleConfig struct {
	Host    []string            `json:"host"`
	Path    []string            `json:"path"`
	Queries map[string]string   `json:"query"`
	Headers map[string][]string `json:"headers"`
}
