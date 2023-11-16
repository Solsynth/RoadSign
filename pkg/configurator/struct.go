package configurator

import "strings"

type AppConfig struct {
	Sites []SiteConfig `json:"sites"`
}

type SiteConfig struct {
	ID string `json:"id"`

	Name      string             `json:"name"`
	Rules     []RouterRuleConfig `json:"rules"`
	Upstreams []UpstreamConfig   `json:"upstreams"`
}

type RouterRuleConfig struct {
	Host    []string            `json:"host"`
	Path    []string            `json:"path"`
	Queries map[string]string   `json:"query"`
	Headers map[string][]string `json:"headers"`
}

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
