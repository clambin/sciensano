package demographics

// TODO: API to get these figures online?

func GetAgeGroupFigures() (groups map[string]int) {
	return map[string]int{
		"0-17":  2569322,
		"18-34": 2151014,
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
