package sideload

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

func getSites(c *fiber.Ctx) error {
	return c.JSON(sign.App.Sites)
}

func getSiteConfig(c *fiber.Ctx) error {
	fp := filepath.Join(viper.GetString("paths.configs"), c.Params("id"))

	var err error
	var data []byte
	if data, err = os.ReadFile(fp + ".yml"); err != nil {
		if data, err = os.ReadFile(fp + ".yaml"); err != nil {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
	}

	return c.Type("yaml").SendString(string(data))
}

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
	if site, ok := lo.Find(sign.App.Sites, func(item *sign.SiteConfig) bool {
		return item.ID == id
	}); ok {
		for _, process := range site.Processes {
			process.StopProcess()
		}
	}

	// Reload
	sign.ReadInConfig(viper.GetString("paths.configs"))
	sign.App.PreheatProcesses()

	return c.SendStatus(fiber.StatusOK)
}
