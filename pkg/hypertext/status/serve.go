package status

import (
	"errors"
	"fmt"
	roadsign "git.solsynth.dev/goatworks/roadsign/pkg"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

type ErrorPayload struct {
	Title   string `json:"title"`
	Message string `json:"message"`
	Version string `json:"version"`
}

func StatusPageHandler(c *fiber.Ctx, err error) error {
	var reqErr *fiber.Error
	var status = fiber.StatusInternalServerError
	if errors.As(err, &reqErr) {
		status = reqErr.Code
	}

	c.Status(status)

	payload := ErrorPayload{
		Version: roadsign.AppVersion,
	}

	switch status {
	case fiber.StatusNotFound:
		payload.Title = "Not Found"
		payload.Message = fmt.Sprintf("no resource for \"%s\"", c.OriginalURL())
		return c.Render("views/not-found", payload)
	case fiber.StatusTooManyRequests:
		payload.Title = "Request Too Fast"
		payload.Message = fmt.Sprintf("you have sent over %d request(s) in a second", viper.GetInt("hypertext.limitation.max_qps"))
		return c.Render("views/too-many-requests", payload)
	case fiber.StatusRequestEntityTooLarge:
		payload.Title = "Request Too Large"
		payload.Message = fmt.Sprintf("you have sent a request over %d bytes", viper.GetInt("hypertext.limitation.max_body_size"))
		return c.Render("views/request-too-large", payload)
	case fiber.StatusBadGateway:
		payload.Title = "Backend Down"
		payload.Message = fmt.Sprintf("all destnations configured to handle your request are down: %s", err.Error())
		return c.Render("views/bad-gateway", payload)
	case fiber.StatusGatewayTimeout:
		payload.Title = "Backend Took Too Long To Response"
		payload.Message = fmt.Sprintf("the destnation took too long to response your request: %s", err.Error())
		return c.Render("views/gateway-timeout", payload)
	default:
		payload.Title = "Oops"
		payload.Message = err.Error()
		return c.Render("views/fallback", payload)
	}
}
