package hypertext

import (
	"math/rand"
	"regexp"

	"code.smartsheep.studio/goatworks/roadsign/pkg/navi"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

func UseProxies(app *fiber.App) {
	app.All("/*", func(ctx *fiber.Ctx) error {
		host := ctx.Hostname()
		path := ctx.Path()
		queries := ctx.Queries()
		headers := ctx.GetReqHeaders()

		// Filtering sites
		for _, region := range navi.R.Regions {
			// Matching rules
			for _, location := range region.Locations {
				if !lo.Contains(location.Host, host) {
					continue
				}

				if !func() bool {
					flag := false
					for _, pattern := range location.Path {
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
				for rk, rv := range location.Queries {
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
				for rk, rv := range location.Headers {
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

				idx := rand.Intn(len(location.Destinations))
				dest := location.Destinations[idx]

				// Passing all the rules means the site is what we are looking for.
				// Let us respond to our client!
				return makeResponse(ctx, &dest)
			}
		}

		// There is no site available for this request.
		// Just ignore it and give our client a not found status.
		// Do not care about the user experience, we can do it in custom error handler.
		return fiber.ErrNotFound
	})
}

func makeResponse(ctx *fiber.Ctx, dest *navi.Destination) error {
	// Modify request
	for _, transformer := range dest.Transformers {
		if err := transformer.TransformRequest(ctx); err != nil {
			return err
		}
	}

	// Forward
	err := navi.R.Forward(ctx, dest)

	// Modify response
	for _, transformer := range dest.Transformers {
		if err := transformer.TransformResponse(ctx); err != nil {
			return err
		}
	}

	return err
}
