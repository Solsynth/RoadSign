package sideload

import (
	roadsign "code.smartsheep.studio/goatworks/roadsign/pkg"
	"github.com/gofiber/fiber/v2"
)

func responseConnectivity(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"server":  "RoadSign",
		"version": roadsign.AppVersion,
	})
}
