package navi

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

var R *RoadApp

func ReadInConfig(root string) error {
	instance := &RoadApp{
		Regions: make([]*Region, 0),
		Metrics: &RoadMetrics{
			Traces:       make([]RoadTrace, 0),
			Traffic:      make(map[string]int64),
			TrafficFrom:  make(map[string]int64),
			TotalTraffic: 0,
		},
	}

	if err := filepath.Walk(root, func(fp string, info os.FileInfo, _ error) error {
		var region Region
		if info.IsDir() {
			return nil
		} else if !strings.HasSuffix(info.Name(), ".toml") {
			return nil
		} else if file, err := os.OpenFile(fp, os.O_RDONLY, 0755); err != nil {
			return err
		} else if data, err := io.ReadAll(file); err != nil {
			return err
		} else if err := toml.Unmarshal(data, &region); err != nil {
			return err
		} else {
			defer file.Close()

			if region.Disabled {
				return nil
			}

			instance.Regions = append(instance.Regions, &region)
		}

		return nil
	}); err != nil {
		return err
	}

	R = instance

	return nil
}
