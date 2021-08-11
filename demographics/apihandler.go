package demographics

// GetAgeGroupFigures is used by apiHandler to get all age group figures
func (server *Server) GetAgeGroupFigures() (groups map[string]int) {
	groups = make(map[string]int)

	for _, bracket := range server.GetAgeBrackets() {
		count, _ := server.GetByAge(bracket)
		groups[bracket.String()] = count
	}

	return
}

// GetRegionFigures is used by apiHandler to get all region figures
func (server *Server) GetRegionFigures() (groups map[string]int) {
	groups = map[string]int{
		"Ostbelgien": 77949,
	}

	for _, region := range server.GetRegions() {
		count, _ := server.GetByRegion(region)

		if region == "Wallonia" {
			count -= groups["Ostbelgien"]
		}

		groups[region] = count
	}

	return
}
