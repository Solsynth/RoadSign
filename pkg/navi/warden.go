package navi

import (
	"git.solsynth.dev/goatworks/roadsign/pkg/warden"
	"github.com/rs/zerolog/log"
)

func InitializeWarden(regions []*Region) {
	pool := make([]*warden.AppInstance, 0)

	log.Info().Msg("Starting Warden applications...")

	for _, region := range regions {
		for _, application := range region.Applications {
			pool = append(pool, &warden.AppInstance{
				Manifest: application,
			})
		}
	}

	// Hot swap
	warden.InstancePool = pool
	errs := warden.StartPool()

	log.Info().Any("errs", errs).Msg("Warden applications has been started.")
}
