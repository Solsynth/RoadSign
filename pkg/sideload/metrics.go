package sideload

import (
	"git.solsynth.dev/goatworks/roadsign/pkg/navi"
	"github.com/gofiber/fiber/v2"
)

func getTraffic(c *fiber.Ctx) error {
	return c.JSON(navi.R.Metrics)
}

func getTraces(c *fiber.Ctx) error {
	return c.JSON(navi.R.Metrics.ReadTrace())
}
