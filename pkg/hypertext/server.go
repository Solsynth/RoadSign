package hypertext

import (
	"crypto/tls"
	jsoniter "github.com/json-iterator/go"
	"net"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func InitServer() *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:               "RoadSign",
		ServerHeader:          "RoadSign",
		DisableStartupMessage: true,
		EnableIPValidation:    true,
		JSONDecoder:           jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal,
		JSONEncoder:           jsoniter.ConfigCompatibleWithStandardLibrary.Marshal,
		Prefork:               viper.GetBool("performance.prefork"),
		BodyLimit:             viper.GetInt("hypertext.limitation.max_body_size"),
	})

	if viper.GetBool("hypertext.force_https") {
		app.Use(func(c *fiber.Ctx) error {
			if !c.Secure() {
				return c.Redirect(
					strings.Replace(c.Request().URI().String(), "http", "https", 1),
					fiber.StatusMovedPermanently,
				)
			}

			return c.Next()
		})
	}

	if viper.GetBool("telemetry.request_logging") {
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

type CertificateConfig struct {
	Key string `json:"key"`
	Pem string `json:"pem"`
}

func RunServer(app *fiber.App, ports []string, securedPorts []string) {
	var certs []CertificateConfig
	raw, _ := jsoniter.Marshal(viper.Get("hypertext.certificate"))
	jsoniter.Unmarshal(raw, &certs)

	tlsCfg := &tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{},
	}

	for _, info := range certs {
		cert, err := tls.LoadX509KeyPair(info.Pem, info.Key)
		if err != nil {
			log.Error().Err(err).
				Str("pem", info.Pem).
				Str("key", info.Key).
				Msg("An error occurred when loading certificate.")
		} else {
			tlsCfg.Certificates = append(tlsCfg.Certificates, cert)
		}
	}

	for _, port := range ports {
		port := port
		go func() {
			if viper.GetBool("hypertext.redirect_to_https") {
				redirector := fiber.New(fiber.Config{
					AppName:               "RoadSign",
					ServerHeader:          "RoadSign",
					DisableStartupMessage: true,
					EnableIPValidation:    true,
				})
				redirector.All("/", func(c *fiber.Ctx) error {
					return c.Redirect(strings.ReplaceAll(string(c.Request().URI().FullURI()), "http", "https"))
				})
				if err := redirector.Listen(port); err != nil {
					log.Panic().Err(err).Msg("An error occurred when listening hypertext non-tls ports.")
				}
			} else {
				if err := app.Listen(port); err != nil {
					log.Panic().Err(err).Msg("An error occurred when listening hypertext non-tls ports.")
				}
			}
		}()
	}

	for _, port := range securedPorts {
		port := port
		go func() {
			listener, err := net.Listen("tcp", port)
			if err != nil {
				log.Panic().Err(err).Msg("An error occurred when listening hypertext tls ports.")
			}
			if err := app.Listener(tls.NewListener(listener, tlsCfg)); err != nil {
				log.Panic().Err(err).Msg("An error occurred when listening hypertext tls ports.")
			}
		}()
	}
}
