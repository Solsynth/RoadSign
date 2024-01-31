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

var AppVersion = "2.0.0-delta2"
