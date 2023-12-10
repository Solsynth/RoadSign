package administration

import (
	"fmt"

	roadsign "code.smartsheep.studio/goatworks/roadsign/pkg"
	"github.com/gofiber/fiber/v2"
)

func responseConnectivity(c *fiber.Ctx) error {
	return c.
		Status(fiber.StatusOK).
		SendString(fmt.Sprintf("Hello from RoadSign v%s", roadsign.AppVersion))
}
