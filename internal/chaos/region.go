package chaos

type RegionChaos struct {
	Region string
	Config Config
}

func ApplyRegion(region string, configs []RegionChaos) {
	for _, rc := range configs {
		if rc.Region == region && rc.Config.Enabled {
			Apply(rc.Config)
		}
	}
}
