package administration

import (
	"os"
	"path/filepath"

	"code.smartsheep.studio/goatworks/roadsign/pkg/filesystem"
	"code.smartsheep.studio/goatworks/roadsign/pkg/sign"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

func doPublish(c *fiber.Ctx) error {
	var upstream *sign.UpstreamConfig
	var site *sign.SiteConfig
	for _, item := range sign.App.Sites {
		if item.ID == c.Params("site") {
			site = item
			for _, stream := range item.Upstreams {
				if stream.ID == c.Params("upstream") {
					upstream = stream
					break
				}
			}
			break
		}
	}

	if upstream == nil {
		return fiber.ErrNotFound
	} else if upstream.GetType() != sign.UpstreamTypeFile {
		return fiber.ErrUnprocessableEntity
	}

	for _, process := range site.Processes {
		process.StopProcess()
	}

	workdir, _ := upstream.GetRawURI()

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
					_ = filesystem.Unzip(dst, workdir)
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
