package recreation

import (
	"context"
	"net/http"
	"sync"
	"time"
)

// This file contains all handlerfuncs being used as cloud functions.
// Have gone with cloud functions for this build due to low cost and ability to easily deploy to multiple regions.
// Note these functions rely on global variables in the package and an init.
// I don't love this but there's no way to pass in state otherwise since the function signature is set for a cloud function.

// HandleAvailabilitiesSync executes one run of the availability sync
func HandleAvailabilitiesSyncConcurrent(w http.ResponseWriter, r *http.Request) {

	now := time.Now()
	var genWG, lisWG sync.WaitGroup
	errCHAN := make(chan error)
	var err error

	ctx, cancel := context.WithCancel(r.Context())

	// if we get an error, cancel the context to abort all the other checks
	lisWG.Add(1)
	go func() {
		defer lisWG.Done()
		for receivedErr := range errCHAN {
			err = receivedErr

		}
		cancel()
	}()

	for i := 0; i < 3; i++ {
		genWG.Add(1)

		go func(i int) {
			defer genWG.Done()
			targetTime := time.Date(now.Year(), now.Month()+time.Month(i), now.Day(), 0, 0, 0, 0, time.UTC)
			err := AvailabilitiesSync(ctx, log, fs, targetTime, now)
			if err != nil {
				errCHAN <- err
			}

		}(i)
	}
	genWG.Wait()
	close(errCHAN)
	lisWG.Wait()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func HandleAvailabilitiesSync(w http.ResponseWriter, r *http.Request) {

	now := time.Now()

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	for i := 0; i < 3; i++ {
		// not async to spread out request load.
		targetTime := time.Date(now.Year(), now.Month()+time.Month(i), now.Day(), 0, 0, 0, 0, time.UTC)
		err := AvailabilitiesSync(ctx, log, fs, targetTime, now)
		if err != nil {
			cancel()
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

// HandleProxyRequest is how I'm obfuscating my IP. Any request sent to this function gets relayed to
// recreation.gov. It is effective at avoiding the cloudflare rate limiting being imposed.
func HandleProxyRequest(w http.ResponseWriter, r *http.Request) {
	proxy.ServeHTTP(w, r)
}
