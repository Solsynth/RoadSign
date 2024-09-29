package hypertext

import (
	"crypto/tls"
	"git.solsynth.dev/goatworks/roadsign/pkg/hypertext/status"
	jsoniter "github.com/json-iterator/go"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func InitServer() *fiber.App {
	views := html.NewFileSystem(http.FS(status.FS), ".gohtml")
	app := fiber.New(fiber.Config{
		ViewsLayout:           "views/index",
		AppName:               "RoadSign",
		ServerHeader:          "RoadSign",
		DisableStartupMessage: true,
		EnableIPValidation:    true,
		Views:                 views,
		ErrorHandler:          status.StatusPageHandler,
		JSONDecoder:           jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal,
		JSONEncoder:           jsoniter.ConfigCompatibleWithStandardLibrary.Marshal,
		ProxyHeader:           fiber.HeaderXForwardedFor,
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
			LimitReached: func(c *fiber.Ctx) error {
				return fiber.ErrTooManyRequests
			},
		}))
	}

	app.All("/*", ProxiesHandler)

	return app
}

type CertificateConfig struct {
	Key string `json:"key"`
	Pem string `json:"pem"`
}

func RunServer(app *fiber.App, ports []string, securedPorts []string) {
	var certs []CertificateConfig
	raw, _ := jsoniter.Marshal(viper.Get("hypertext.certificate"))
	_ = jsoniter.Unmarshal(raw, &certs)

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

		log.Info().Msgf("Listening for %s... http://0.0.0.0%s", app.Config().AppName, port)
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

		log.Info().Msgf("Listening for %s... https://0.0.0.0%s", app.Config().AppName, port)
	}
}
