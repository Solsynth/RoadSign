package hypertext

import (
	roadsign "code.smartsheep.studio/goatworks/roadsign/pkg"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"time"
)

func InitServer() *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:               "RoadSign",
		ServerHeader:          fmt.Sprintf("RoadSign v%s", roadsign.AppVersion),
		DisableStartupMessage: true,
		EnableIPValidation:    true,
		EnablePrintRoutes:     viper.GetBool("debug.print_routes"),
		Prefork:               viper.GetBool("performance.prefork"),
		BodyLimit:             viper.GetInt("hypertext.limitation.max_body_size"),
	})

	if viper.GetInt("hypertext.limitation.max_qps") > 0 {
		app.Use(limiter.New(limiter.Config{
			Max:        viper.GetInt("hypertext.limitation.max_qps"),
			Expiration: 1 * time.Second,
		}))
	}

	UseProxies(app)

	return app
}

func RunServer(app *fiber.App) {
	for _, port := range viper.GetStringSlice("hypertext.ports") {
		port := port
		go func() {
			if err := app.Listen(port); err != nil {
				log.Panic().Err(err).Msg("An error occurred when listening hypertext tls ports.")
			}
		}()
	}

	for _, port := range viper.GetStringSlice("hypertext.secured_ports") {
		port := port
		pem := viper.GetString("hypertext.certificate.pem")
		key := viper.GetString("hypertext.certificate.key")
		go func() {
			if err := app.ListenTLS(port, pem, key); err != nil {
				log.Panic().Err(err).Msg("An error occurred when listening hypertext tls ports.")
			}
		}()
	}
}
