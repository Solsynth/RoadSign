package deploy

import (
	"fmt"
	"io"
	"os"
	"strings"

	jsoniter "github.com/json-iterator/go"

	"code.smartsheep.studio/goatworks/roadsign/pkg/cmd/rds/conn"
	"code.smartsheep.studio/goatworks/roadsign/pkg/navi"
	"github.com/gofiber/fiber/v2"
	"github.com/pelletier/go-toml/v2"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var DeployCommands = []*cli.Command{
	{
		Name:      "deploy",
		Aliases:   []string{"dp"},
		ArgsUsage: "<server> <site> <upstream> [path]",
		Action: func(ctx *cli.Context) error {
			if ctx.Args().Len() < 4 {
				return fmt.Errorf("must have four arguments: <server> <site> <upstream> <path>")
			}

			if !strings.HasSuffix(ctx.Args().Get(3), ".zip") {
				return fmt.Errorf("input file must be a zip file and ends with .zip")
			}

			server, ok := conn.GetConnection(ctx.Args().Get(0))
			if !ok {
				return fmt.Errorf("server was not found, use \"rds connect\" add one first")
			} else if err := server.CheckConnectivity(); err != nil {
				return fmt.Errorf("couldn't connect server: %s", err.Error())
			}

			// Send request
			log.Info().Msg("Now publishing to remote server...")

			url := fmt.Sprintf("/webhooks/publish/%s/%s?mimetype=%s", ctx.Args().Get(1), ctx.Args().Get(2), "application/zip")
			client := fiber.Put(server.Url+url).
				SendFile(ctx.Args().Get(3), "attachments").
				MultipartForm(nil).
				BasicAuth("RoadSign CLI", server.Credential)

			if status, data, err := client.Bytes(); len(err) > 0 {
				return fmt.Errorf("failed to publish to remote: %q", err)
			} else if status != 200 {
				return fmt.Errorf("server rejected request, status code %d, response %s", status, string(data))
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
			} else if err := server.CheckConnectivity(); err != nil {
				return fmt.Errorf("couldn't connect server: %s", err.Error())
			}

			var site navi.SiteConfig
			if file, err := os.Open(ctx.Args().Get(2)); err != nil {
				return err
			} else {
				raw, _ := io.ReadAll(file)
				toml.Unmarshal(raw, &site)
			}

			url := fmt.Sprintf("/webhooks/sync/%s", ctx.Args().Get(1))
			client := fiber.Put(server.Url+url).
				JSONEncoder(jsoniter.ConfigCompatibleWithStandardLibrary.Marshal).
				JSONDecoder(jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal).
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
