package load

type CDNTarget struct {
	Name string
	URL  string
}

type MultiCDN struct {
	Targets []CDNTarget
}

func (m *MultiCDN) Pick(i int) CDNTarget {
	return m.Targets[i%len(m.Targets)]
}
