package navi

import (
	"fmt"
	roadsign "git.solsynth.dev/goatworks/roadsign/pkg"
	"git.solsynth.dev/goatworks/roadsign/pkg/navi/transformers"
	"github.com/spf13/viper"

	"github.com/gofiber/fiber/v2"
)

type RoadApp struct {
	Regions []*Region    `json:"regions"`
	Metrics *RoadMetrics `json:"metrics"`
}

func (v *RoadApp) Forward(c *fiber.Ctx, dest *Destination) error {
	// Add reserve proxy headers
	ip := c.IP()
	scheme := c.Protocol()
	protocol := string(c.Request().Header.Protocol())
	c.Request().Header.Set(fiber.HeaderXForwardedFor, ip)
	c.Request().Header.Set(fiber.HeaderXForwardedHost, ip)
	c.Request().Header.Set(fiber.HeaderXForwardedProto, scheme)
	c.Request().Header.Set(
		fiber.HeaderVia,
		fmt.Sprintf("%s %s", protocol, viper.GetString("central")),
	)
	c.Request().Header.Set(
		fiber.HeaderForwarded,
		fmt.Sprintf("by=%s; for=%s; host=%s; proto=%s", c.IP(), c.IP(), c.Get(fiber.HeaderHost), scheme),
	)

	// Response body
	var err error
	switch dest.GetType() {
	case DestinationHypertext:
		err = makeUnifiedResponse(c, dest)
	case DestinationStaticFile:
		err = makeFileResponse(c, dest)
	default:
		err = fiber.ErrBadGateway
	}

	// Apply helmet
	if dest.Helmet != nil {
		dest.Helmet.Apply(c)
	}

	// Apply watermark
	c.Response().Header.Set(fiber.HeaderServer, "RoadSign")
	c.Response().Header.Set(fiber.HeaderXPoweredBy, fmt.Sprintf("RoadSign %s", roadsign.AppVersion))

	return err
}

type RequestTransformerConfig = transformers.TransformerConfig
