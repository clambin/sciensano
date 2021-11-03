package demographics

import (
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

func (store *Store) update() (byAge map[Bracket]int, byRegion map[string]int, err error) {
	if store.once == nil || time.Now().After(store.expiry) {
		store.once = &sync.Once{}
		store.expiry = time.Now().Add(store.Retention)
	}

	store.once.Do(func() {
		store.byAge, store.byRegion, err = store.refresh()
	})

	if err != nil {
		log.WithError(err).Warning("failed to retrieve demographics")
	}

	return store.byAge, store.byRegion, err
}

func (store *Store) refresh() (byAge map[Bracket]int, byRegion map[string]int, err error) {
	start := time.Now()
	datafile := DataFile{
		TempDirectory: store.TempDirectory,
		URL:           store.URL,
	}
	defer datafile.Remove()

	err = datafile.Download()
	if err != nil {
		return
	}

	var byRegionRaw, byAgeRaw map[string]int
	byRegionRaw, byAgeRaw, err = groupPopulation(datafile.filename)

	if err != nil {
		return
	}

	byAge = groupPopulationByAge(byAgeRaw, store.AgeBrackets)
	byRegion = groupPopulationByRegion(byRegionRaw)

	log.Infof("loaded demographics in %s", time.Now().Sub(start))
	return
}
