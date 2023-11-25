package sign

import (
	"errors"
	"math/rand"

	"github.com/gofiber/fiber/v2"
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
		return makeHypertextResponse(ctx, upstream)
	case UpstreamTypeFile:
		return makeFileResponse(ctx, upstream)
	default:
		return fiber.ErrBadGateway
	}
}

type SiteConfig struct {
	ID           string                     `json:"id"`
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
