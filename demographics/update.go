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

	var err1, err2 error
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		byAge, err1 = datafile.ParseByAge(store.AgeBrackets)

		if err1 != nil {
			log.WithError(err1).Error("unable to parse demographics file")
		}
		wg.Done()
	}()
	go func() {
		byRegion, err2 = datafile.ParseByRegion()
		if err2 != nil {
			log.WithError(err2).Error("unable to parse demographics file")
		}
		wg.Done()
	}()

	wg.Wait()
	if err2 != nil {
		err = err2
	}
	if err1 != nil {
		err = err1
	}
	log.Infof("loaded demographics in %s", time.Now().Sub(start))
	return
}
