package navi

import (
	"errors"
	"math/rand"

	"code.smartsheep.studio/goatworks/roadsign/pkg/navi/transformers"

	"github.com/gofiber/fiber/v2"
)

type RoadApp struct {
	Sites []*SiteConfig `json:"sites"`
}

func (v *RoadApp) Forward(ctx *fiber.Ctx, site *SiteConfig) error {
	if len(site.Upstreams) == 0 {
		return errors.New("invalid configuration")
	}

	// Do forward
	idx := rand.Intn(len(site.Upstreams))
	upstream := site.Upstreams[idx]

	switch upstream.GetType() {
	case UpstreamTypeHypertext:
		return makeHypertextResponse(ctx, upstream)
	case UpstreamTypeFile:
		return makeFileResponse(ctx, upstream)
	default:
		return fiber.ErrBadGateway
	}
}

type RequestTransformerConfig = transformers.RequestTransformerConfig

type SiteConfig struct {
	ID           string                      `json:"id"`
	Rules        []*RouterRule               `json:"rules" yaml:"rules"`
	Transformers []*RequestTransformerConfig `json:"transformers" yaml:"transformers"`
	Upstreams    []*UpstreamInstance         `json:"upstreams" yaml:"upstreams"`
}

type RouterRule struct {
	Host    []string            `json:"host" yaml:"host"`
	Path    []string            `json:"path" yaml:"path"`
	Queries map[string]string   `json:"queries" yaml:"queries"`
	Headers map[string][]string `json:"headers" yaml:"headers"`
}
