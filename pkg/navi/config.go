package navi

import (
	"io"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

var R *RoadApp

func ReadInConfig(root string) error {
	instance := &RoadApp{
		Regions: make([]*Region, 0),
	}

	if err := filepath.Walk(root, func(fp string, info os.FileInfo, _ error) error {
		var region Region
		if info.IsDir() {
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
