package deploy

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"code.smartsheep.studio/goatworks/roadsign/pkg/cmd/rds/conn"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mholt/archiver/v4"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/urfave/cli/v2"
)

var DeployCommands = []*cli.Command{
	{
		Name:      "deploy",
		Aliases:   []string{"dp"},
		ArgsUsage: "<server> <site> <upstream> [path]",
		Action: func(ctx *cli.Context) error {
			if ctx.Args().Len() < 3 {
				return fmt.Errorf("must have three arguments: <server> <site> <upstream> [path]")
			}

			server, ok := conn.GetConnection(ctx.Args().Get(0))
			if !ok {
				return fmt.Errorf("server was not found, use \"rds connect\" add one first")
			}

			// Prepare file to upload
			log.Info().Msg("Preparing file to upload, please stand by...")

			filelist, err := archiver.FilesFromDisk(&archiver.FromDiskOptions{FollowSymlinks: true}, map[string]string{
				lo.Ternary(ctx.Args().Len() > 3, ctx.Args().Get(4), "."): "",
			})
			if err != nil {
				return fmt.Errorf("failed to prepare file: %q", err)
			}

			workdir, _ := os.Getwd()
			filename := filepath.Join(workdir, fmt.Sprintf("rds-deploy-cache-%s.zip", uuid.NewString()))
			out, err := os.Create(filename)
			if err != nil {
				return fmt.Errorf("failed to prepare file: %q", err)
			}
			defer out.Close()

			if err := (archiver.Zip{}).Archive(context.Background(), out, filelist); err != nil {
				return fmt.Errorf("failed to prepare file: %q", err)
			}

			// Send request
			log.Info().Msg("Now publishing to remote server...")

			url := fmt.Sprintf("/webhooks/publish/%s/%s?mimetype=%s", ctx.Args().Get(1), ctx.Args().Get(2), "application/zip")
			client := fiber.Put(server.Url+url).
				SendFile(filename, "attachments").
				MultipartForm(nil).
				BasicAuth("RoadSign CLI", server.Credential)

			var mistake error
			if status, _, err := client.Bytes(); len(err) > 0 {
				mistake = fmt.Errorf("failed to publish to remote: %q", err)
			} else if status != 200 {
				mistake = fmt.Errorf("server rejected request, may cause by invalid credential")
			}

			// Cleanup
			log.Info().Msg("Cleaning up...")
			os.Remove(filename)

			if mistake != nil {
				return mistake
			}

			// Well done!
			log.Info().Msg("Well done! Your site is successfully published! 🎉")

			return nil
		},
	},
}
