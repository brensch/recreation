package main

import (
	"context"
	"net/http"
	"sync"
	"time"

	"cloud.google.com/go/firestore"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"go.uber.org/zap"
)

func HandleAvailabilitySync(log *zap.Logger, fs *firestore.Client, ifdb influxdb2.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		concurrentCHAN := make(chan struct{}, 4)
		var wg sync.WaitGroup

		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		now := time.Now()
		for _, groundID := range campgroundIDs {
			time.Sleep(100 * time.Millisecond)

			wg.Add(1)
			go func(groundID string) {
				concurrentCHAN <- struct{}{}
				defer func() { <-concurrentCHAN }()
				defer wg.Done()
				err := DoAvailabilitySync(ctx, log, fs, ifdb, now, groundID)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					cancel()
					return
				}
			}(groundID)
		}

		wg.Wait()

		w.WriteHeader(http.StatusOK)
	}
}
