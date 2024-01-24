package navi

import "code.smartsheep.studio/goatworks/roadsign/pkg/warden"

func InitializeWarden(regions []*Region) {
	for _, region := range regions {
		for _, application := range region.Applications {
			warden.InstancePool = append(warden.InstancePool, &warden.AppInstance{
				Manifest: application,
			})
		}
	}
	
	for _, instance := range warden.InstancePool {
		instance.Wake()
	}
}
