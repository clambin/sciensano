package demographics

import (
	log "github.com/sirupsen/logrus"
	"time"
)

// Update refreshes the demographics cache
func (store *Store) Update() {
	start := time.Now()
	if byAge, byRegion, err := store.refresh(); err == nil {
		// only lock when we have retrieved the data (as this can take over a minute on RPI)
		// so that handlers that need demographics data are never blocked while we refresh
		store.lock.Lock()
		store.byAge = byAge
		store.byRegion = byRegion
		store.lock.Unlock()
		log.Infof("loaded demographics in %s", time.Now().Sub(start))
	} else {
		log.WithError(err).Warning("failed to retrieve demographics")
	}
	return
}

func (store *Store) refresh() (byAge map[Bracket]int, byRegion map[string]int, err error) {
	datafile := DataFile{
		TempDirectory: store.TempDirectory,
		URL:           store.URL,
	}
	defer datafile.Remove()

	if err = datafile.Download(); err == nil {
		var byRegionRaw, byAgeRaw map[string]int
		byRegionRaw, byAgeRaw, err = groupPopulation(datafile.filename)

		if err == nil {
			byAge = groupPopulationByAge(byAgeRaw, store.AgeBrackets)
			byRegion = groupPopulationByRegion(byRegionRaw)
		}
	}
	return
}
