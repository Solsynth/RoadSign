package configurator

import (
	"math/rand"
)

type AppConfig struct {
	Sites []SiteConfig `json:"sites"`
}

func (v *AppConfig) LoadBalance(site SiteConfig) *UpstreamConfig {
	idx := rand.Intn(len(site.Upstreams))

	switch site.Upstreams[idx].GetType() {
	case UpstreamTypeHypertext:
		return &site.Upstreams[idx]
	case UpstreamTypeFile:
		// TODO Make this into hypertext configuration
		return &site.Upstreams[idx]
	default:
		// Give him the null value when this configuration is invalid.
		// Then we can print a log in the console to warm him.
		return nil
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
