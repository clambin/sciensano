package demographics

import (
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

func (store *Store) update() (err error) {
	var (
		byAge    map[Bracket]int
		byRegion map[string]int
	)
	store.lock.Lock()
	if store.once == nil || time.Now().After(store.expiry) {
		store.once = &sync.Once{}
		store.expiry = time.Now().Add(store.Retention)
	}
	store.lock.Unlock()

	store.once.Do(func() {
		start := time.Now()
		byAge, byRegion, err = store.refresh()
		if err == nil {
			store.lock.Lock()
			store.byAge = byAge
			store.byRegion = byRegion
			store.expiry = time.Now().Add(store.Retention)
			store.lock.Unlock()
			log.Infof("loaded demographics in %s", time.Now().Sub(start))
		} else {
			log.WithError(err).Warning("failed to retrieve demographics")
		}
	})
	return
}

func (store *Store) refresh() (byAge map[Bracket]int, byRegion map[string]int, err error) {
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

	if err == nil {
		byAge = groupPopulationByAge(byAgeRaw, store.AgeBrackets)
		byRegion = groupPopulationByRegion(byRegionRaw)
	}
	return
}
