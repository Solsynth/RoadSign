package sideload

import (
	"fmt"
	"net/http"

	"code.smartsheep.studio/goatworks/roadsign/pkg/sideload/view"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	jsoniter "github.com/json-iterator/go"

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
		ServerHeader:          "RoadSign Sideload",
		DisableStartupMessage: true,
		EnableIPValidation:    true,
		JSONDecoder:           jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal,
		JSONEncoder:           jsoniter.ConfigCompatibleWithStandardLibrary.Marshal,
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

	app.Use("/", filesystem.New(filesystem.Config{
		Root:         http.FS(view.FS),
		PathPrefix:   "dist",
		Index:        "index.html",
		NotFoundFile: "index.html",
	}))

	cgi := app.Group("/cgi").Name("CGI")
	{
		cgi.Get("/metadata", getMetadata)
		cgi.Get("/statistics", getStatistics)
		cgi.Get("/sites", getRegions)
		cgi.Get("/sites/cfg/:id", getRegionConfig)
		cgi.Get("/processes", getApplications)
		cgi.Get("/processes/logs/:id", getApplicationLogs)
	}

	webhooks := app.Group("/webhooks").Name("WebHooks")
	{
		webhooks.Put("/publish/:site/:slug", doPublish)
		webhooks.Put("/sync/:slug", doSync)
	}

	return app
}
