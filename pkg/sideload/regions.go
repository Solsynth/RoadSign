package sideload

import (
	"fmt"
	"os"
	"path/filepath"

	"git.solsynth.dev/goatworks/roadsign/pkg/navi"
	"git.solsynth.dev/goatworks/roadsign/pkg/warden"
	"github.com/gofiber/fiber/v2"
	"github.com/pelletier/go-toml/v2"
	"github.com/samber/lo"
	"github.com/spf13/viper"
)

func getRegions(c *fiber.Ctx) error {
	return c.JSON(navi.R.Regions)
}

func getRegionConfig(c *fiber.Ctx) error {
	fp := filepath.Join(viper.GetString("paths.configs"), c.Params("id"))

	var err error
	var data []byte
	if data, err = os.ReadFile(fp + ".toml"); err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.Type("toml").SendString(string(data))
}

func doSync(c *fiber.Ctx) error {
	req := string(c.Body())

	id := c.Params("slug")
	path := filepath.Join(viper.GetString("paths.configs"), fmt.Sprintf("%s.toml", id))

	if file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755); err != nil {
		return fiber.NewError(fiber.ErrInternalServerError.Code, err.Error())
	} else {
		var testOut map[string]any
		if err := toml.Unmarshal([]byte(req), &testOut); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid configuration: %v", err))
		}
		_, _ = file.Write([]byte(req))
		defer file.Close()
	}

	var stopQueue, startQueue []*warden.AppInstance
	// Getting things need to stop
	if region, ok := lo.Find(navi.R.Regions, func(item *navi.Region) bool {
		return item.ID == id
	}); ok {
		for _, application := range region.Applications {
			if instance := warden.GetFromPool(application.ID); instance != nil {
				stopQueue = append(stopQueue, instance)
			}
		}
	}

	// Reload
	_ = navi.ReadInConfig(viper.GetString("paths.configs"))

	// Getting things need to start
	if region, ok := lo.Find(navi.R.Regions, func(item *navi.Region) bool {
		return item.ID == id
	}); ok {
		for _, application := range region.Applications {
			if instance := warden.GetFromPool(application.ID); instance != nil {
				startQueue = append(startQueue, instance)
			}
		}
	}

	// Reboot
	for _, instance := range stopQueue {
		_ = instance.Stop()
	}
	for _, instance := range startQueue {
		_ = instance.Wake()
	}

	return c.SendStatus(fiber.StatusOK)
}
