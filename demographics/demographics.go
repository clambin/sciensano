package demographics

// TODO: API to get these figures online?

func GetAgeGroupFigures() (groups map[string]int) {
	return map[string]int{
		"0-11":  1525064,
		"12-15": 533441,
		"16-17": 253617,
		"18-24": 922803,
		"25-34": 1486440,
		"35-44": 1485793,
		"45-54": 1558559,
		"55-64": 1523475,
		"65-74": 1170399,
		"75-84": 698940,
		"85+":   335139,
	}
}

func GetRegionFigures() (groups map[string]int) {
	return map[string]int{
		"Flanders":   6629143,
		"Wallonia":   3645243,
		"Brussels":   1218255,
		"Ostbelgien": 77949,
	}
}
