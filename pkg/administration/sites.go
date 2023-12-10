package administration

import (
	"fmt"
	"os"
	"path/filepath"

	"code.smartsheep.studio/goatworks/roadsign/pkg/sign"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

func doSyncSite(c *fiber.Ctx) error {
	var req sign.SiteConfig

	if err := c.BodyParser(&req); err != nil {
		return err
	}

	id := c.Params("slug")
	path := filepath.Join(viper.GetString("paths.configs"), fmt.Sprintf("%s.yaml", id))

	if file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755); err != nil {
		return fiber.NewError(fiber.ErrInternalServerError.Code, err.Error())
	} else {
		raw, _ := yaml.Marshal(req)
		file.Write(raw)
		defer file.Close()
	}

	sign.App.Sites = lo.Map(sign.App.Sites, func(item *sign.SiteConfig, idx int) *sign.SiteConfig {
		if item.ID == id {
			return &req
		} else {
			return item
		}
	})

	return c.SendStatus(fiber.StatusOK)
}
