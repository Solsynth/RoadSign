package sideload

import (
	"git.solsynth.dev/goatworks/roadsign/pkg/navi"
	"git.solsynth.dev/goatworks/roadsign/pkg/warden"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

func getApplications(c *fiber.Ctx) error {
	applications := lo.FlatMap(navi.R.Regions, func(item *navi.Region, idx int) []warden.ApplicationInfo {
		return lo.Map(item.Applications, func(item warden.Application, index int) warden.ApplicationInfo {
			return warden.ApplicationInfo{
				Application: item,
				Status:      warden.GetFromPool(item.ID).Status,
			}
		})
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

func letApplicationStart(c *fiber.Ctx) error {
	if instance, ok := lo.Find(warden.InstancePool, func(item *warden.AppInstance) bool {
		return item.Manifest.ID == c.Params("id")
	}); !ok {
		return fiber.NewError(fiber.StatusNotFound)
	} else {
		if err := instance.Wake(); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.SendStatus(fiber.StatusOK)
	}
}

func letApplicationStop(c *fiber.Ctx) error {
	if instance, ok := lo.Find(warden.InstancePool, func(item *warden.AppInstance) bool {
		return item.Manifest.ID == c.Params("id")
	}); !ok {
		return fiber.NewError(fiber.StatusNotFound)
	} else {
		if err := instance.Stop(); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.SendStatus(fiber.StatusOK)
	}
}

func letApplicationRestart(c *fiber.Ctx) error {
	if instance, ok := lo.Find(warden.InstancePool, func(item *warden.AppInstance) bool {
		return item.Manifest.ID == c.Params("id")
	}); !ok {
		return fiber.NewError(fiber.StatusNotFound)
	} else {
		if err := instance.Stop(); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		if err := instance.Start(); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.SendStatus(fiber.StatusOK)
	}
}
