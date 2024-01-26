package sideload

import (
	"code.smartsheep.studio/goatworks/roadsign/pkg/navi"
	"code.smartsheep.studio/goatworks/roadsign/pkg/warden"
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
	})
}
