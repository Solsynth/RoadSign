package conn

import (
	"encoding/json"
	"fmt"
	"strings"

	roadsign "code.smartsheep.studio/goatworks/roadsign/pkg"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type CliConnection struct {
	ID         string `json:"id"`
	Url        string `json:"url"`
	Credential string `json:"credential"`
}

func (v CliConnection) CheckConnectivity() error {
	client := fiber.Get(v.Url + "/cgi/connectivity")
	client.BasicAuth("RoadSign CLI", v.Credential)

	if status, data, err := client.Bytes(); len(err) > 0 {
		return fmt.Errorf("couldn't connect to server: %q", err)
	} else if status != 200 {
		return fmt.Errorf("server rejected request, may cause by invalid credential")
	} else {
		var resp fiber.Map
		if err := json.Unmarshal(data, &resp); err != nil {
			return err
		} else if resp["server"] != "RoadSign" {
			return fmt.Errorf("remote server isn't roadsign")
		} else if resp["version"] != roadsign.AppVersion {
			if strings.Contains(roadsign.AppVersion, "#") {
				return fmt.Errorf("remote server version mismatch client version, update or downgrade client required")
			} else {
				log.Warn().Msg("RoadSign CLI didn't complied with vcs information, compatibility was disabled. To enable it, reinstall cli with -buildvcs flag.")
			}
		}
	}
	return nil
}
