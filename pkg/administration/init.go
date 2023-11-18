package administration

import (
	roadsign "code.smartsheep.studio/goatworks/roadsign/pkg"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func InitAdministration() *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:               "RoadSign Administration",
		ServerHeader:          fmt.Sprintf("RoadSign Administration v%s", roadsign.AppVersion),
		DisableStartupMessage: true,
		EnableIPValidation:    true,
		EnablePrintRoutes:     viper.GetBool("debug.print_routes"),
		TrustedProxies:        viper.GetStringSlice("security.administration_trusted_proxies"),
	})

	webhooks := app.Group("/webhooks").Name("WebHooks")
	{
		webhooks.Put("/publish/:site/:upstream", doPublish)
	}

	return app
}
