package sign

import (
	"errors"
	"math/rand"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type AppConfig struct {
	Sites []*SiteConfig `json:"sites"`
}

func (v *AppConfig) Forward(ctx *fiber.Ctx, site *SiteConfig) error {
	if len(site.Upstreams) == 0 {
		return errors.New("invalid configuration")
	}

	// Boot processes
	for _, process := range site.Processes {
		if err := process.BootProcess(); err != nil {
			log.Warn().Err(err).Msgf("An error occurred when booting process (%s) for %s", process.ID, site.ID)
			return fiber.ErrBadGateway
		} else {
			log.Debug().Msg("process is alive!")
		}
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

type SiteConfig struct {
	ID           string                      `json:"id"`
	Rules        []*RouterRuleConfig         `json:"rules" yaml:"rules"`
	Transformers []*RequestTransformerConfig `json:"transformers" yaml:"transformers"`
	Upstreams    []*UpstreamConfig           `json:"upstreams" yaml:"upstreams"`
	Processes    []*ProcessConfig            `json:"processes" yaml:"processes"`
}

type RouterRuleConfig struct {
	Host    []string            `json:"host" yaml:"host"`
	Path    []string            `json:"path" yaml:"path"`
	Queries map[string]string   `json:"queries" yaml:"queries"`
	Headers map[string][]string `json:"headers" yaml:"headers"`
}
