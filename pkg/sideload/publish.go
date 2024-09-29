package sideload

import (
	"context"
	"git.solsynth.dev/goatworks/roadsign/pkg/warden"
	"os"
	"path/filepath"

	"git.solsynth.dev/goatworks/roadsign/pkg/navi"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/saracen/fastzip"
)

func doPublish(c *fiber.Ctx) error {
	var workdir string
	var destination *navi.Destination
	var application *warden.Application
	for _, item := range navi.R.Regions {
		if item.ID == c.Params("site") {
			for _, location := range item.Locations {
				for _, dest := range location.Destinations {
					if dest.ID == c.Params("slug") {
						destination = &dest
						workdir, _ = dest.GetRawUri()
						break
					}
				}
			}
			for _, app := range item.Applications {
				if app.ID == c.Params("slug") {
					application = &app
					workdir = app.Workdir
					break
				}
			}
			break
		}
	}

	var instance *warden.AppInstance
	if application != nil {
		if instance = warden.GetFromPool(application.ID); instance != nil {
			instance.Stop()
		}
	} else if destination != nil && destination.GetType() != navi.DestinationStaticFile {
		return fiber.ErrUnprocessableEntity
	} else {
		return fiber.ErrNotFound
	}

	if c.Query("overwrite", "yes") == "yes" {
		files, _ := filepath.Glob(filepath.Join(workdir, "*"))
		for _, file := range files {
			_ = os.Remove(file)
		}
	}

	if form, err := c.MultipartForm(); err == nil {
		files := form.File["attachments"]
		for _, file := range files {
			mimetype := lo.Ternary(len(c.Query("mimetype")) > 0, c.Query("mimetype"), file.Header["Content-Type"][0])
			switch mimetype {
			case "application/zip":
				dst := filepath.Join(os.TempDir(), uuid.NewString()+".zip")
				if err := c.SaveFile(file, dst); err != nil {
					return err
				} else {
					if ex, err := fastzip.NewExtractor(dst, workdir); err != nil {
						return err
					} else if err = ex.Extract(context.Background()); err != nil {
						defer ex.Close()
						return err
					}
				}
			default:
				dst := filepath.Join(workdir, file.Filename)
				if err := c.SaveFile(file, dst); err != nil {
					return err
				}
			}
		}
	}

	if instance != nil {
		instance.Wake()
	}

	return c.SendStatus(fiber.StatusOK)
}
