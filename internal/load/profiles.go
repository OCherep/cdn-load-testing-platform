package load

type Profile struct {
	Targets map[string]string `json:"targets"`

	MinRPS int `json:"min_rps"`
	MaxRPS int `json:"max_rps"`
	Step   int `json:"step"`

	GeoDistribution map[string]int `json:"geo_distribution"`
}
