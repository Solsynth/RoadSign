package administration

import (
	"code.smartsheep.studio/goatworks/roadsign/pkg/fs"
	"code.smartsheep.studio/goatworks/roadsign/pkg/sign"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"os"
	"path/filepath"
)

func doPublish(ctx *fiber.Ctx) error {
	var upstream *sign.UpstreamConfig
	for _, item := range sign.App.Sites {
		if item.ID == ctx.Params("site") {
			for _, stream := range item.Upstreams {
				if stream.ID == ctx.Params("upstream") {
					upstream = &stream
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

	workdir, _ := upstream.GetRawURI()

	if ctx.Query("overwrite", "yes") == "yes" {
		files, _ := filepath.Glob(filepath.Join(workdir, "*"))
		for _, file := range files {
			_ = os.Remove(file)
		}
	}

	if form, err := ctx.MultipartForm(); err == nil {
		files := form.File["attachments"]
		for _, file := range files {
			mimetype := file.Header["Content-Type"][0]
			switch mimetype {
			case "application/zip":
				dst := filepath.Join(os.TempDir(), uuid.NewString()+".zip")
				if err := ctx.SaveFile(file, dst); err != nil {
					return err
				} else {
					_ = fs.Unzip(dst, workdir)
				}
			default:
				dst := filepath.Join(workdir, file.Filename)
				if err := ctx.SaveFile(file, dst); err != nil {
					return err
				}
			}
		}
	}

	return ctx.SendStatus(fiber.StatusOK)
}
