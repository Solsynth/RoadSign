package sideload

import (
	"context"
	"os"
	"path/filepath"

	"code.smartsheep.studio/goatworks/roadsign/pkg/sign"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/saracen/fastzip"
)

func doPublish(c *fiber.Ctx) error {
	var workdir string
	var site *sign.SiteConfig
	var upstream *sign.UpstreamConfig
	var process *sign.ProcessConfig
	for _, item := range sign.App.Sites {
		if item.ID == c.Params("site") {
			site = item
			for _, stream := range item.Upstreams {
				if stream.ID == c.Params("slug") {
					upstream = stream
					workdir, _ = stream.GetRawURI()
					break
				}
			}
			for _, proc := range item.Processes {
				if proc.ID == c.Params("slug") {
					process = proc
					workdir = proc.Workdir
					break
				}
			}
			break
		}
	}

	if upstream == nil && process == nil {
		return fiber.ErrNotFound
	} else if upstream != nil && upstream.GetType() != sign.UpstreamTypeFile {
		return fiber.ErrUnprocessableEntity
	}

	for _, process := range site.Processes {
		process.StopProcess()
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

	return c.SendStatus(fiber.StatusOK)
}
