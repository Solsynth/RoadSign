package sideload

import (
	"code.smartsheep.studio/goatworks/roadsign/pkg/navi"
	"github.com/gofiber/fiber/v2"
)

func getTraces(c *fiber.Ctx) error {
	return c.JSON(navi.R.Traces)
}