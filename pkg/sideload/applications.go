package sideload

import (
	"git.solsynth.dev/goatworks/roadsign/pkg/navi"
	"git.solsynth.dev/goatworks/roadsign/pkg/warden"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

func getApplications(c *fiber.Ctx) error {
	applications := lo.FlatMap(navi.R.Regions, func(item *navi.Region, idx int) []warden.Application {
		return item.Applications
	})

	return c.JSON(applications)
}

func getApplicationLogs(c *fiber.Ctx) error {
	if instance, ok := lo.Find(warden.InstancePool, func(item *warden.AppInstance) bool {
		return item.Manifest.ID == c.Params("id")
	}); !ok {
		return fiber.NewError(fiber.StatusNotFound)
	} else {
		return c.SendString(instance.Logs())
	}
}
