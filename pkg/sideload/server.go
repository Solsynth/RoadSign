package sideload

import (
	"fmt"
	roadsign "git.solsynth.dev/goatworks/roadsign/pkg"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/logger"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func InitSideload() *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:               "RoadSign Sideload",
		ServerHeader:          "RoadSign Sideload",
		DisableStartupMessage: true,
		EnableIPValidation:    true,
		JSONDecoder:           jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal,
		JSONEncoder:           jsoniter.ConfigCompatibleWithStandardLibrary.Marshal,
		ProxyHeader:           fiber.HeaderXForwardedFor,
		EnablePrintRoutes:     viper.GetBool("debug.print_routes"),
		TrustedProxies:        viper.GetStringSlice("sideload.trusted_proxies"),
		BodyLimit:             viper.GetInt("hypertext.limitation.max_body_size"),
	})

	if viper.GetBool("telemetry.request_logging") {
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
		cgi.Get("/metadata", getMetadata)
		cgi.Get("/traffic", getTraffic)
		cgi.Get("/traces", getTraces)
		cgi.Get("/stats", getStats)
		cgi.Get("/regions", getRegions)
		cgi.Get("/regions/cfg/:id", getRegionConfig)
		cgi.Get("/applications", getApplications)
		cgi.Get("/applications/:id/logs", getApplicationLogs)

		cgi.Post("/reload", doReload)
	}

	webhooks := app.Group("/webhooks").Name("WebHooks")
	{
		webhooks.Put("/publish/:site/:slug", doPublish)
		webhooks.Put("/sync/:slug", doSync)
	}

	return app
}
