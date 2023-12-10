package deploy

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"code.smartsheep.studio/goatworks/roadsign/pkg/cmd/rds/conn"
	"code.smartsheep.studio/goatworks/roadsign/pkg/sign"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mholt/archiver/v4"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
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
			cleanup := true
			workdir, _ := os.Getwd()
			var filename string
			if ctx.Args().Len() < 3 || !strings.HasSuffix(ctx.Args().Get(3), ".zip") {
				log.Info().Msg("Preparing file to upload, please stand by...")

				filelist, err := archiver.FilesFromDisk(nil, map[string]string{
					lo.Ternary(ctx.Args().Len() > 3, ctx.Args().Get(4), "."): "",
				})
				if err != nil {
					return fmt.Errorf("failed to prepare file: %q", err)
				}

				filename = filepath.Join(workdir, fmt.Sprintf("rds-deploy-cache-%s.zip", uuid.NewString()))
				out, err := os.Create(filename)
				if err != nil {
					return fmt.Errorf("failed to prepare file: %q", err)
				}
				defer out.Close()

				if err := (archiver.Zip{}).Archive(context.Background(), out, filelist); err != nil {
					return fmt.Errorf("failed to prepare file: %q", err)
				}
			} else if ctx.Args().Len() > 3 {
				cleanup = false
				filename = ctx.Args().Get(3)
			}

			// Send request
			log.Info().Msg("Now publishing to remote server...")

			url := fmt.Sprintf("/webhooks/publish/%s/%s?mimetype=%s", ctx.Args().Get(1), ctx.Args().Get(2), "application/zip")
			client := fiber.Put(server.Url+url).
				SendFile(filename, "attachments").
				MultipartForm(nil).
				BasicAuth("RoadSign CLI", server.Credential)

			var mistake error
			if status, data, err := client.Bytes(); len(err) > 0 {
				mistake = fmt.Errorf("failed to publish to remote: %q", err)
			} else if status != 200 {
				mistake = fmt.Errorf("server rejected request, status code %d, response %s", status, string(data))
			}

			// Cleanup
			if cleanup {
				log.Info().Msg("Cleaning up...")
				os.Remove(filename)
			}

			if mistake != nil {
				return mistake
			}

			log.Info().Msg("Well done! Your site is successfully published! 🎉")

			return nil
		},
	},
	{
		Name:      "sync",
		Aliases:   []string{"sc"},
		ArgsUsage: "<server> <site> <configuration path>",
		Action: func(ctx *cli.Context) error {
			if ctx.Args().Len() < 3 {
				return fmt.Errorf("must have three arguments: <server> <site> <configuration path>")
			}

			server, ok := conn.GetConnection(ctx.Args().Get(0))
			if !ok {
				return fmt.Errorf("server was not found, use \"rds connect\" add one first")
			}

			var site sign.SiteConfig
			if file, err := os.Open(ctx.Args().Get(2)); err != nil {
				return err
			} else {
				raw, _ := io.ReadAll(file)
				yaml.Unmarshal(raw, &site)
			}

			url := fmt.Sprintf("/webhooks/sync/%s", ctx.Args().Get(1))
			client := fiber.Put(server.Url+url).
				JSON(site).
				BasicAuth("RoadSign CLI", server.Credential)

			if status, data, err := client.Bytes(); len(err) > 0 {
				return fmt.Errorf("failed to sync to remote: %q", err)
			} else if status != 200 {
				return fmt.Errorf("server rejected request, status code %d, response %s", status, string(data))
			}

			log.Info().Msg("Well done! Your site configuration is up-to-date! 🎉")

			return nil
		},
	},
}
