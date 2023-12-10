package conn

import (
	"encoding/json"

	"github.com/samber/lo"
	"github.com/spf13/viper"
)

func GetConnection(id string) (CliConnection, bool) {
	var servers []CliConnection
	raw, _ := json.Marshal(viper.Get("servers"))
	_ = json.Unmarshal(raw, &servers)
	return lo.Find(servers, func(item CliConnection) bool {
		return item.ID == id
	})
}
