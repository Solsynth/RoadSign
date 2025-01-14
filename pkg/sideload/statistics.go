package sideload

import (
	"time"

	"git.solsynth.dev/goatworks/roadsign/pkg/navi"
	"git.solsynth.dev/goatworks/roadsign/pkg/warden"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

func getStats(c *fiber.Ctx) error {
	locations := lo.FlatMap(navi.R.Regions, func(item *navi.Region, idx int) []navi.Location {
		return item.Locations
	})
	destinations := lo.FlatMap(locations, func(item navi.Location, idx int) []navi.Destination {
		return item.Destinations
	})
	applications := lo.FlatMap(navi.R.Regions, func(item *navi.Region, idx int) []warden.Application {
		return item.Applications
	})

	return c.JSON(fiber.Map{
		"regions":      len(navi.R.Regions),
		"locations":    len(locations),
		"destinations": len(destinations),
		"applications": len(applications),
		"uptime":       time.Since(navi.R.Metrics.StartupAt).Milliseconds(),
		"traffic": fiber.Map{
			"total": navi.R.Metrics.TotalTraffic,
		},
	})
}
