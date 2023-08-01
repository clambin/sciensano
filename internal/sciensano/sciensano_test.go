package sciensano_test

var (
	wantAgeGroups = []string{"(unknown)", "0-9", "10-19", "20-29", "30-39", "40-49", "50-59", "60-69", "70-79", "80-89", "90+"}
	wantRegions   = []string{"(unknown)", "Brussels", "Flanders", "Wallonia"}
	wantProvinces = []string{"(unknown)", "Antwerpen", "BrabantWallon", "Brussels", "Hainaut", "Limburg", "Liège", "Luxembourg", "Namur", "OostVlaanderen", "VlaamsBrabant", "WestVlaanderen"}
)

/*
var (
	regions       = []string{"(unknown)", "Brussels", "Flanders", "Ostbelgien", "Wallonia"}
	provinces     = []string{"foo", "bar", "snafu"}
	ageGroups     = []string{"00-04", "05-11", "12-15", "16-17", "18-24", "25-34", "35-44", "45-54", "55-64", "65-74", "75-84", "85+"}
	manufacturers = []string{"AstraZeneca-Oxford", "Johnson&Johnson", "Moderna", "Novavax", "Other", "Pfizer-BioNTech"}
	doses         = []sciensano.DoseType{sciensano.Partial, sciensano.Full, sciensano.SingleDose, sciensano.Booster, sciensano.Booster2, sciensano.Booster3}
)

func makeResponse[T any](count int, maker func(timestamp time.Time, region, province, ageGroup, manufacturer string, dose sciensano.DoseType) *T) []*T {
	response := make([]*T, 0)
	timestamp := time.Date(2022, 11, 22, 0, 0, 0, 0, time.UTC)

	metaValue := reflect.ValueOf(new(T)).Elem()
	actualRegions := regions
	if metaValue.FieldByName("Region") == (reflect.Value{}) {
		actualRegions = []string{""}
	}
	actualProvinces := provinces
	if metaValue.FieldByName("Province") == (reflect.Value{}) {
		actualProvinces = []string{""}
	}
	actualAgeGroups := ageGroups
	if metaValue.FieldByName("AgeGroup") == (reflect.Value{}) {
		actualAgeGroups = []string{""}
	}
	actualManufacturers := manufacturers
	if metaValue.FieldByName("Manufacturer") == (reflect.Value{}) {
		actualManufacturers = []string{""}
	}
	actualDoses := doses
	if metaValue.FieldByName("Dose") == (reflect.Value{}) {
		actualDoses = []sciensano.DoseType{sciensano.Partial}
	}

	for i := 0; i < count; i++ {
		for _, region := range actualRegions {
			for _, province := range actualProvinces {
				for _, ageGroup := range actualAgeGroups {
					for _, manufacturer := range actualManufacturers {
						for _, dose := range actualDoses {
							response = append(response, maker(timestamp, region, province, ageGroup, manufacturer, dose))
						}
					}
				}
			}
		}
		timestamp = timestamp.Add(24 * time.Hour)
	}
	return response
}


*/
