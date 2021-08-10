package demographics

// TODO: API to get these figures online?

func GetAgeGroupFigures() (groups map[string]int) {
	return map[string]int{
		"00-11": 1525064,
		"12-15": 533441,
		"16-17": 253617,
		"18-24": 922803,
		"25-34": 1486440,
		"35-44": 1493556,
		"45-54": 1538889,
		"55-64": 1538307,
		"65-74": 1194091,
		"75-84": 703107,
		"85+":   331923,
	}
}

func GetRegionFigures() (groups map[string]int) {
	return map[string]int{
		"Flanders":   6653062,
		"Wallonia":   3570257,
		"Brussels":   1219970,
		"Ostbelgien": 77949,
	}
}
