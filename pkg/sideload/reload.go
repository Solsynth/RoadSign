package sideload

import (
	"git.solsynth.dev/goatworks/roadsign/pkg/navi"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func doReload(c *fiber.Ctx) error {
	if err := navi.ReadInConfig(viper.GetString("paths.configs")); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if c.QueryBool("warden", false) {
		navi.InitializeWarden(navi.R.Regions)
	}

	return c.SendStatus(fiber.StatusOK)
}
