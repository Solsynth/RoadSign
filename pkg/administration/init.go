package administration

import (
	roadsign "code.smartsheep.studio/goatworks/roadsign/pkg"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
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

	app.Use(basicauth.New(basicauth.Config{
		Realm: fmt.Sprintf("RoadSign v%s", roadsign.AppVersion),
		Authorizer: func(_, password string) bool {
			return password == viper.GetString("security.credential")
		},
	}))

	webhooks := app.Group("/webhooks").Name("WebHooks")
	{
		webhooks.Put("/publish/:site/:upstream", doPublish)
	}

	return app
}
