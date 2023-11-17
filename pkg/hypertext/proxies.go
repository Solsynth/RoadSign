package hypertext

import (
	"code.smartsheep.studio/goatworks/roadsign/pkg/configurator"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"regexp"
	"time"
)

func UseProxies(app *fiber.App) {
	app.All("/*", func(ctx *fiber.Ctx) error {
		host := ctx.Hostname()
		path := ctx.Path()
		queries := ctx.Queries()
		headers := ctx.GetReqHeaders()

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

func makeResponse(ctx *fiber.Ctx, site configurator.SiteConfig) error {
	// Load balance
	upstream := configurator.C.LoadBalance(site)
	if upstream == nil {
		log.Warn().Str("id", site.ID).Msg("There is no available upstream for this request.")
		return fiber.ErrBadGateway
	}

	// Modify request
	for _, transformer := range site.Transformers {
		transformer.TransformRequest(ctx)
	}

	// Perform forward
	timeout := time.Duration(viper.GetInt64("performance.network_timeout")) * time.Millisecond
	err := proxy.Do(ctx, upstream.MakeURI(ctx), &fasthttp.Client{
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
	})

	// Modify response
	for _, transformer := range site.Transformers {
		transformer.TransformResponse(ctx)
	}

	return err
}
