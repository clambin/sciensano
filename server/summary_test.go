package server_test

/*
func TestSummaryHandler_Query(t *testing.T) {
	targets := []struct {
		sciensanoType string
		summaryTypes  set.Set[sciensano.SummaryColumn]
	}{
		{sciensanoType: "cases", summaryTypes: sciensano.CasesValidSummaryModes()},
		{sciensanoType: "hospitalisations", summaryTypes: sciensano.HospitalisationsValidSummaryModes()},
		{sciensanoType: "mortalities", summaryTypes: sciensano.MortalitiesValidSummaryModes()},
		{sciensanoType: "testResults", summaryTypes: sciensano.TestResultsValidSummaryModes()},
		{sciensanoType: "vaccinations", summaryTypes: sciensano.VaccinationsValidSummaryModes()},
	}

	type summarizer interface {
		Summarize(summaryColumn sciensano.SummaryColumn) (*tabulator.Tabulator, error)
	}

	for _, target := range targets {
		var records summarizer
		var summaryTypes set.Set[sciensano.SummaryColumn]
		switch target.sciensanoType {
		case "cases":
			records = testutil.Cases()
			summaryTypes = sciensano.CasesValidSummaryModes()
		case "hospitalisations":
			records = testutil.Hospitalisations()
			summaryTypes = sciensano.HospitalisationsValidSummaryModes()
		case "mortalities":
			records = testutil.Mortalities()
			summaryTypes = sciensano.MortalitiesValidSummaryModes()
		case "testResults":
			records = testutil.TestResults()
			summaryTypes = sciensano.TestResultsValidSummaryModes()
		case "vaccinations":
			records = testutil.Vaccinations()
			summaryTypes = sciensano.VaccinationsValidSummaryModes()
		}

		for _, summaryType := range summaryTypes.List() {
			t.Run(target.sciensanoType+"-"+summaryType.String(), func(t *testing.T) {
				s := mocks.NewReportsStorer(t)
				report, _ := records.Summarize(summaryType)
				expectedColumns := 1 + len(report.GetColumns())
				s.EXPECT().Get(target.sciensanoType+"-"+summaryType.String()).Return(report, nil).Once()

				r := server.SummaryHandler{ReportsStore: s}

				req := grafanaJSONServer.QueryRequest{
					Targets: []grafanaJSONServer.QueryRequestTarget{{
						Target:  target.sciensanoType,
						Payload: []byte(fmt.Sprintf(`{ "summary": "%s", "accumulate": "no" }`, summaryType.String())),
					}},
					Range: grafanaJSONServer.Range{From: time.Now().Add(-24 * time.Hour)},
				}

				resp, err := r.Query(context.Background(), target.sciensanoType, req)
				require.NoError(t, err)
				assert.Len(t, resp.(grafanaJSONServer.TableResponse).Columns, expectedColumns)
			})
		}
	}
}


*/
