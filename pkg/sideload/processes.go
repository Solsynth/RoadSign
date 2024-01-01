package sideload

import (
	"code.smartsheep.studio/goatworks/roadsign/pkg/sign"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

func getProcesses(c *fiber.Ctx) error {
	processes := lo.FlatMap(sign.App.Sites, func(item *sign.SiteConfig, idx int) []*sign.ProcessInstance {
		return item.Processes
	})

	return c.JSON(processes)
}

func getProcessLog(c *fiber.Ctx) error {
	processes := lo.FlatMap(sign.App.Sites, func(item *sign.SiteConfig, idx int) []*sign.ProcessInstance {
		return item.Processes
	})

	if target, ok := lo.Find(processes, func(item *sign.ProcessInstance) bool {
		return item.ID == c.Params("id")
	}); !ok {
		return fiber.NewError(fiber.StatusNotFound)
	} else {
		return c.SendString(target.GetLogs())
	}
}
