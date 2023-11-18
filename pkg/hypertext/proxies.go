package hypertext

import (
	"code.smartsheep.studio/goatworks/roadsign/pkg/sign"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"regexp"
)

func UseProxies(app *fiber.App) {
	app.All("/*", func(ctx *fiber.Ctx) error {
		host := ctx.Hostname()
		path := ctx.Path()
		queries := ctx.Queries()
		headers := ctx.GetReqHeaders()

		// Filtering sites
		for _, site := range sign.App.Sites {
			// Matching rules
			for _, rule := range site.Rules {
				if !lo.Contains(rule.Host, host) {
					continue
				}

				if !func() bool {
					flag := false
					for _, pattern := range rule.Path {
						if ok, _ := regexp.MatchString(pattern, path); ok {
							flag = true
							break
						}
					}
					return flag
				}() {
					continue
				}

				// Filter query strings
				flag := true
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

				// Passing all the rules means the site is what we are looking for.
				// Let us respond to our client!
				return makeResponse(ctx, site)
			}
		}

		// There is no site available for this request.
		// Just ignore it and give our client a not found status.
		// Do not care about the user experience, we can do it in custom error handler.
		return fiber.ErrNotFound
	})
}

func makeResponse(ctx *fiber.Ctx, site sign.SiteConfig) error {
	// Modify request
	for _, transformer := range site.Transformers {
		transformer.TransformRequest(ctx)
	}

	// Forward
	err := sign.App.Forward(ctx, site)

	// Modify response
	for _, transformer := range site.Transformers {
		transformer.TransformResponse(ctx)
	}

	return err
}
