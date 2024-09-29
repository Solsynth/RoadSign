package navi

import "git.solsynth.dev/goatworks/roadsign/pkg/warden"

func InitializeWarden(regions []*Region) {
	pool := make([]*warden.AppInstance, 0)

	for _, region := range regions {
		for _, application := range region.Applications {
			pool = append(pool, &warden.AppInstance{
				Manifest: application,
			})
		}
	}

	// Hot swap
	warden.InstancePool = pool

	for _, instance := range warden.InstancePool {
		instance.Wake()
	}
}
