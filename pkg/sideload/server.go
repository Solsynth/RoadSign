package sideload

import (
	"fmt"

	roadsign "code.smartsheep.studio/goatworks/roadsign/pkg"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func InitSideload() *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:               "RoadSign Sideload",
		ServerHeader:          fmt.Sprintf("RoadSign Sideload v%s", roadsign.AppVersion),
		DisableStartupMessage: true,
		EnableIPValidation:    true,
		EnablePrintRoutes:     viper.GetBool("debug.print_routes"),
		TrustedProxies:        viper.GetStringSlice("security.sideload_trusted_proxies"),
		BodyLimit:             viper.GetInt("hypertext.limitation.max_body_size"),
	})

	if viper.GetBool("performance.request_logging") {
		app.Use(logger.New(logger.Config{
			Output: log.Logger,
			Format: "[Sideload] [${time}] ${status} - ${latency} ${method} ${path}\n",
		}))
	}

	app.Use(basicauth.New(basicauth.Config{
		Realm: fmt.Sprintf("RoadSign v%s", roadsign.AppVersion),
		Authorizer: func(_, password string) bool {
			return password == viper.GetString("security.credential")
		},
	}))

	cgi := app.Group("/cgi").Name("CGI")
	{
		cgi.All("/connectivity", responseConnectivity)
	}

	webhooks := app.Group("/webhooks").Name("WebHooks")
	{
		webhooks.Put("/publish/:site/:slug", doPublish)
		webhooks.Put("/sync/:slug", doSyncSite)
	}

	return app
}
