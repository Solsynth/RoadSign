package hypertext

import (
	"math/rand"
	"regexp"

	"github.com/spf13/viper"

	"git.solsynth.dev/goatworks/roadsign/pkg/navi"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

func ProxiesHandler(ctx *fiber.Ctx) error {
	host := ctx.Hostname()
	path := ctx.Path()
	queries := ctx.Queries()
	headers := ctx.GetReqHeaders()

	// Filtering sites
	for _, region := range navi.R.Regions {
		// Matching rules
		for _, location := range region.Locations {
			if !lo.Contains(location.Hosts, host) {
				continue
			}

			if !func() bool {
				flag := false
				for _, pattern := range location.Paths {
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
			return makeResponse(ctx, region, &location, &dest)
		}
	}

	// There is no site available for this request.
	// Just ignore it and give our client a not found status.
	// Do not care about the user experience, we can do it in custom error handler.
	return fiber.ErrNotFound
}

func makeResponse(c *fiber.Ctx, region *navi.Region, location *navi.Location, dest *navi.Destination) error {
	uri := c.Request().URI().String()

	// Modify request
	for _, transformer := range dest.Transformers {
		if err := transformer.TransformRequest(c); err != nil {
			return err
		}
	}

	// Forward
	err := navi.R.Forward(c, dest)

	// Modify response
	for _, transformer := range dest.Transformers {
		if err := transformer.TransformResponse(c); err != nil {
			return err
		}
	}

	// Collect trace
	if viper.GetBool("telemetry.capture_traces") {
		var message string
		if err != nil {
			message = err.Error()
		}

		go navi.R.Metrics.AddTrace(navi.RoadTrace{
			Region:      region.ID,
			Location:    location.ID,
			Destination: dest.ID,
			Uri:         uri,
			IpAddress:   c.IP(),
			UserAgent:   c.Get(fiber.HeaderUserAgent),
			Error: navi.RoadTraceError{
				IsNull:  err == nil,
				Message: message,
			},
		})
	}

	return err
}
