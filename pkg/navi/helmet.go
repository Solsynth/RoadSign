package navi

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type HelmetConfig struct {
	XSSProtection             string `json:"xss_protection" toml:"xss_protection"`
	ContentTypeNosniff        string `json:"content_type_nosniff" toml:"content_type_nosniff"`
	XFrameOptions             string `json:"x_frame_options" toml:"x_frame_options"`
	HSTSMaxAge                int    `json:"hsts_max_age" toml:"hsts_max_age"`
	HSTSExcludeSubdomains     bool   `json:"hsts_exclude_subdomains" toml:"hsts_exclude_subdomains"`
	ContentSecurityPolicy     string `json:"content_security_policy" toml:"content_security_policy"`
	CSPReportOnly             bool   `json:"csp_report_only" toml:"csp_report_only"`
	HSTSPreloadEnabled        bool   `json:"hsts_preload_enabled" toml:"hsts_preload_enabled"`
	ReferrerPolicy            string `json:"referrer_policy" toml:"referrer_policy"`
	PermissionPolicy          string `json:"permission_policy" toml:"permission_policy"`
	CrossOriginEmbedderPolicy string `json:"cross_origin_embedder_policy" toml:"cross_origin_embedder_policy"`
	CrossOriginOpenerPolicy   string `json:"cross_origin_opener_policy" toml:"cross_origin_opener_policy"`
	CrossOriginResourcePolicy string `json:"cross_origin_resource_policy" toml:"cross_origin_resource_policy"`
	OriginAgentCluster        string `json:"origin_agent_cluster" toml:"origin_agent_cluster"`
	XDNSPrefetchControl       string `json:"xdns_prefetch_control" toml:"xdns_prefetch_control"`
	XDownloadOptions          string `json:"x_download_options" toml:"x_download_options"`
	XPermittedCrossDomain     string `json:"x_permitted_cross_domain" toml:"x_permitted_cross_domain"`
}

func (cfg HelmetConfig) Apply(c *fiber.Ctx) {
	// Apply other headers
	if cfg.XSSProtection != "" {
		c.Set(fiber.HeaderXXSSProtection, cfg.XSSProtection)
	}
	if cfg.ContentTypeNosniff != "" {
		c.Set(fiber.HeaderXContentTypeOptions, cfg.ContentTypeNosniff)
	}
	if cfg.XFrameOptions != "" {
		c.Set(fiber.HeaderXFrameOptions, cfg.XFrameOptions)
	}
	if cfg.CrossOriginEmbedderPolicy != "" {
		c.Set("Cross-Origin-Embedder-Policy", cfg.CrossOriginEmbedderPolicy)
	}
	if cfg.CrossOriginOpenerPolicy != "" {
		c.Set("Cross-Origin-Opener-Policy", cfg.CrossOriginOpenerPolicy)
	}
	if cfg.CrossOriginResourcePolicy != "" {
		c.Set("Cross-Origin-Resource-Policy", cfg.CrossOriginResourcePolicy)
	}
	if cfg.OriginAgentCluster != "" {
		c.Set("Origin-Agent-Cluster", cfg.OriginAgentCluster)
	}
	if cfg.ReferrerPolicy != "" {
		c.Set("Referrer-Policy", cfg.ReferrerPolicy)
	}
	if cfg.XDNSPrefetchControl != "" {
		c.Set("X-DNS-Prefetch-Control", cfg.XDNSPrefetchControl)
	}
	if cfg.XDownloadOptions != "" {
		c.Set("X-Download-Options", cfg.XDownloadOptions)
	}
	if cfg.XPermittedCrossDomain != "" {
		c.Set("X-Permitted-Cross-Domain-Policies", cfg.XPermittedCrossDomain)
	}

	// Handle HSTS headers
	if c.Protocol() == "https" && cfg.HSTSMaxAge != 0 {
		subdomains := ""
		if !cfg.HSTSExcludeSubdomains {
			subdomains = "; includeSubDomains"
		}
		if cfg.HSTSPreloadEnabled {
			subdomains = fmt.Sprintf("%s; preload", subdomains)
		}
		c.Set(fiber.HeaderStrictTransportSecurity, fmt.Sprintf("max-age=%d%s", cfg.HSTSMaxAge, subdomains))
	}

	// Handle Content-Security-Policy headers
	if cfg.ContentSecurityPolicy != "" {
		if cfg.CSPReportOnly {
			c.Set(fiber.HeaderContentSecurityPolicyReportOnly, cfg.ContentSecurityPolicy)
		} else {
			c.Set(fiber.HeaderContentSecurityPolicy, cfg.ContentSecurityPolicy)
		}
	}

	// Handle Permissions-Policy headers
	if cfg.PermissionPolicy != "" {
		c.Set(fiber.HeaderPermissionsPolicy, cfg.PermissionPolicy)
	}
}
