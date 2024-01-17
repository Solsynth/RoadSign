package sideload

import (
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

func getStatistics(c *fiber.Ctx) error {
	upstreams := lo.FlatMap(sign.App.Sites, func(item *sign.SiteConfig, idx int) []*sign.UpstreamInstance {
		return item.Upstreams
	})
	processes := lo.FlatMap(sign.App.Sites, func(item *sign.SiteConfig, idx int) []*sign.ProcessInstance {
		return item.Processes
	})
	unhealthy := lo.FlatMap(sign.App.Sites, func(item *sign.SiteConfig, idx int) []*sign.ProcessInstance {
		return lo.Filter(item.Processes, func(item *sign.ProcessInstance, idx int) bool {
			return item.Status != navi.ProcessStarted
		})
	})

	return c.JSON(fiber.Map{
		"sites":     len(sign.App.Sites),
		"upstreams": len(upstreams),
		"processes": len(processes),
		"status":    len(unhealthy) == 0,
	})
}
