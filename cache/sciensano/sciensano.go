package sciensano

const BaseURL = "https://epistat.sciensano.be"

var Routes = map[string]string{
	"cases":            "/Data/COVID19BE_CASES_AGESEX.json",
	"hospitalisations": "/Data/COVID19BE_HOSP.json",
	"mortalities":      "/Data/COVID19BE_MORT.json",
	"testResults":      "/Data/COVID19BE_tests.json",
	"vaccinations":     "/Data/COVID19BE_VACC.json",
}
