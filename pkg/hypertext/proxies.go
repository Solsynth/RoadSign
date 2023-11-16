package hypertext

import (
	"code.smartsheep.studio/goatworks/roadsign/pkg/configurator"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

func UseProxies(app *fiber.App) {
	app.All("/", func(ctx *fiber.Ctx) error {
		host := ctx.Hostname()
		path := ctx.Path()
		queries := ctx.Queries()
		headers := ctx.GetReqHeaders()

		log.Debug().
			Any("host", host).
			Any("path", path).
			Any("queries", queries).
			Any("headers", headers).
			Msg("A new request received")

		// Filtering sites
		for _, site := range configurator.C.Sites {
			// Matching rules
			for _, rule := range site.Rules {
				if !lo.Contains(rule.Host, host) {
					continue
				} else if !lo.ContainsBy(rule.Path, func(item string) bool {
					matched, err := regexp.MatchString(item, path)
					return matched && err == nil
				}) {
					continue
				}

				flag := true

				// Filter query strings
				for rk, rv := range rule.Queries {
					for ik, iv := range queries {
						if rk != ik && rv != iv {
							flag = false
							break
						}
					}
					if !flag {
						break
					}
				}
				if !flag {
					continue
				}

				// Filter headers
				for rk, rv := range rule.Headers {
					for ik, iv := range headers {
						if rk == ik {
							for _, ov := range iv {
								if !lo.Contains(rv, ov) {
									flag = false
									break
								}
							}
						}
						if !flag {
							break
						}
					}
					if !flag {
						break
					}
				}
				if !flag {
					continue
				}
			}

			// Passing all the rules means the site is what we are looking for.
			// Let us respond to our client!
			return makeResponse(ctx, site)
		}

		// There is no site available for this request.
		// Just ignore it and give our client a not found status.
		// Do not care about the user experience, we can do it in custom error handler.
		return fiber.ErrNotFound
	})
}

func doLoadBalance(site configurator.SiteConfig) *configurator.UpstreamConfig {
	idx := rand.Intn(len(site.Upstreams))

	switch site.Upstreams[idx].GetType() {
	case configurator.UpstreamTypeHypertext:
		return &site.Upstreams[idx]
	case configurator.UpstreamTypeFile:
		// TODO Make this into hypertext configuration
		return &site.Upstreams[idx]
	default:
		// Give him the null value when this configuration is invalid.
		// Then we can print a log in the console to warm him.
		return nil
	}
}

func makeRequestUri(ctx *fiber.Ctx, upstream configurator.UpstreamConfig) string {
	var queries []string
	for k, v := range ctx.Queries() {
		queries = append(queries, fmt.Sprintf("%s=%s", k, v))
	}

	hash := string(ctx.Request().URI().Hash())

	return upstream.URI +
		lo.Ternary(len(queries) > 0, "?"+strings.Join(queries, "&"), "") +
		lo.Ternary(len(hash) > 0, "#"+hash, "")
}

func makeResponse(ctx *fiber.Ctx, site configurator.SiteConfig) error {
	upstream := doLoadBalance(site)
	if upstream == nil {
		log.Warn().Str("id", site.ID).Msg("There is no available upstream for this request.")
		return fiber.ErrBadGateway
	}

	timeout := time.Duration(viper.GetInt64("performance.network_timeout")) * time.Millisecond
	return proxy.Do(ctx, makeRequestUri(ctx, *upstream), &fasthttp.Client{
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
	})
}
