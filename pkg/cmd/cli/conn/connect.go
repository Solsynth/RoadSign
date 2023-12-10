package conn

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type CliConnection struct {
	Url        string `json:"url"`
	Credential string `json:"credential"`
}

func (v CliConnection) GetConnectivity() error {
	client := fiber.Get(v.Url + "/cgi/connectivity")
	client.BasicAuth("RoadSign CLI", v.Credential)

	if status, _, err := client.String(); len(err) > 0 {
		return fmt.Errorf("couldn't connect to server: %q", err)
	} else if status != 200 {
		return fmt.Errorf("server rejected request, may cause by invalid credential")
	}

	return nil
}
