package roadsign

import (
	"runtime/debug"
)

func init() {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				AppVersion += "#" + setting.Value
			}
		}
	}
}

var AppVersion = "1.2.1"
