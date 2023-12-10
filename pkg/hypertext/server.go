package hypertext

import (
	"fmt"
	"time"

	roadsign "code.smartsheep.studio/goatworks/roadsign/pkg"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func InitServer() *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:               "RoadSign",
		ServerHeader:          fmt.Sprintf("RoadSign v%s", roadsign.AppVersion),
		DisableStartupMessage: true,
		EnableIPValidation:    true,
		Prefork:               viper.GetBool("performance.prefork"),
		BodyLimit:             viper.GetInt("hypertext.limitation.max_body_size"),
	})

	if viper.GetBool("performance.request_logging") {
		app.Use(logger.New(logger.Config{
			Output: log.Logger,
			Format: "[Proxies] [${time}] ${status} - ${latency} ${method} ${path}\n",
		}))
	}

	if viper.GetInt("hypertext.limitation.max_qps") > 0 {
		app.Use(limiter.New(limiter.Config{
			Max:        viper.GetInt("hypertext.limitation.max_qps"),
			Expiration: 1 * time.Second,
		}))
	}

	UseProxies(app)

	return app
}

func RunServer(app *fiber.App, ports []string, securedPorts []string, pem string, key string) {
	for _, port := range ports {
		port := port
		go func() {
			if err := app.Listen(port); err != nil {
				log.Panic().Err(err).Msg("An error occurred when listening hypertext tls ports.")
			}
		}()
	}

	for _, port := range securedPorts {
		port := port
		go func() {
			if err := app.ListenTLS(port, pem, key); err != nil {
				log.Panic().Err(err).Msg("An error occurred when listening hypertext tls ports.")
			}
		}()
	}
}
