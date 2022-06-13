package recreation

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/brensch/recreation/api"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	CollectionAvailability       = "availability"
	CollectionAvailabilityDeltas = "availability_deltas"
)

// GetAvailabilityRef gives us a consistent ID to work with for our documents
func GetAvailabilityRef(groundID string, targetTime time.Time) string {
	targetTime = api.GetStartOfMonth(targetTime)
	return fmt.Sprintf("%s-%s", groundID, targetTime.Format(time.RFC3339))
}

// CheckForAvailabilityChange gets old and new states of availability and returns the new availability and deltas to old availability.
func CheckForAvailabilityChange(ctx context.Context, log *zap.Logger, baseURI string, fs *firestore.Client, targetTime, now time.Time, targetGround string) (api.Availability, []CampsiteDelta, error) {

	// get new availability from API
	newAvailability, err := api.GetAvailability(ctx, log, baseURI, targetGround, targetTime)
	if err != nil {
		return api.Availability{}, nil, err
	}

	// get old availability from firestore
	// NotFound errors ignored since the document not existing just results in an empty object, as intended
	oldAvailabilitySnap, err := fs.Collection(CollectionAvailability).Doc(GetAvailabilityRef(targetGround, targetTime)).Get(ctx)
	if err != nil && status.Code(err) != codes.NotFound {
		return api.Availability{}, nil, err
	}
	var oldAvailability api.Availability
	err = oldAvailabilitySnap.DataTo(&oldAvailability)
	if err != nil && status.Code(err) != codes.NotFound {
		return api.Availability{}, nil, err
	}

	// compare the old and new availabilities
	deltas, err := FindAvailabilityDeltas(oldAvailability, newAvailability, targetGround, now)
	return newAvailability, deltas, err
}

func GetAllAvailabilities(ctx context.Context, log *zap.Logger, targetTime, now time.Time) ([]AvailabilityDetailed, error) {

	// start is used to measure duration and as the change time in the delta object
	start := time.Now()
	syncID := uuid.New().String()

	log = log.With(
		zap.String("sync_id", syncID),
		zap.Time("target_time", targetTime),
	)
	log.Info("started getting all availabilities")

	// groundChunks := ChunkGroundsUp(proxies, campgroundIDs)

	var allDeltas []CampsiteDelta
	var allNewAvails []AvailabilityDetailed
	// allNewAvails := make(map[string]api.Availability)

	// allow up to five concurrent requests
	concurrencyLimiter := make(chan struct{}, 5)
	var mu sync.Mutex
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	rand.Seed(time.Now().UnixMicro())

	for _, targetGround := range campgroundIDs {
		wg.Add(1)
		go func(targetGround string) {
			defer wg.Done()
			// write to the buffered chan. as soon as it fills up, this will become blocking until other goroutines end.
			concurrencyLimiter <- struct{}{}
			defer func() { <-concurrencyLimiter }()

			randRegion := proxies[rand.Intn(len(proxies))]
			log := log.With(
				zap.String("proxy", randRegion),
				zap.String("ground_id", targetGround),
			)

			log.Debug("checking campground availability")
			proxyURI := fmt.Sprintf("https://%s-campr-app.cloudfunctions.net/HandleProxyRequest", randRegion)

			newAvailability, err := api.GetAvailability(ctx, log, proxyURI, targetGround, targetTime)
			if err != nil {
				log.Error("failed to check availability change", zap.Error(err))
				cancel()
				return
			}

			mu.Lock()
			allDeltas = append(allDeltas, deltas...)
			allNewAvails = append(allNewAvails, AvailabilityDetailed{
				Availability: newAvailability,
				Month:        api.GetStartOfMonth(targetTime),
				GroundID:     targetGround,
			})
			mu.Unlock()
		}(targetGround)
	}

	wg.Wait()

	// if there are no deltas at all, end here
	if len(allDeltas) == 0 {
		log.Info("checked availabilities, no deltas",
			zap.Duration("duration", time.Since(start)),
			zap.Int("deltas", len(allDeltas)),
		)
		return nil
	}

	updateStart := time.Now()
	log.Debug("finished scraping, uploading availabilities that have changed")

	// update new availabilities all at once
	for ref, availabiltity := range allNewAvails {
		_, err := fs.Collection(CollectionAvailability).Doc(ref).Set(ctx, availabiltity)
		if err != nil {
			log.Error("couldn't write new availability to firestore", zap.Error(err))
		}
	}

	commitEnd := time.Now()

	// update deltas in firestore.
	// add all deltas to the same document since campsites are globally unique so cheaper for consumers to grab the whole
	// document and search through for campsites you want.
	// breaking it up into chunks of 500 sites to avoid max document size limit
	chunkSize := 500
	startIndex := 0
	for startIndex < len(allDeltas) {
		endIndex := startIndex + chunkSize
		if endIndex > len(allDeltas) {
			endIndex = len(allDeltas)
		}

		checkDelta := CheckDelta{
			Deltas:    allDeltas[startIndex:endIndex],
			CheckTime: start,
		}
		ref, _, err := fs.Collection(CollectionAvailabilityDeltas).Add(
			ctx,
			checkDelta,
		)
		if err != nil {
			log.Error("couldn't add availability deltas to firestore", zap.Error(err))
			return err
		}
		log.Debug("synced chunk of deltas", zap.String("firestore_document", ref.ID))
		startIndex += chunkSize
	}

	log.Warn("synced deltas to firestore",
		zap.Int("delta_count", len(allDeltas)),
		zap.Duration("commit_duration", commitEnd.Sub(updateStart)),
		zap.Duration("delta_update_duration", time.Since(commitEnd)),
		zap.Duration("delta_retrieve_duration", updateStart.Sub(start)),
	)

	return nil
}

// func AvailabilitiesSyncConcurrent(ctx context.Context, log *zap.Logger, fs *firestore.Client, targetTime, now time.Time) error {

// 	// start is used to measure duration and as the change time in the delta object
// 	start := time.Now()
// 	syncID := uuid.New().String()

// 	log = log.With(
// 		zap.String("sync_id", syncID),
// 		zap.Time("target_time", targetTime),
// 	)
// 	log.Info("started availability sync routine")

// 	groundChunks := ChunkGroundsUp(proxies, campgroundIDs)

// 	type chunkOutput struct {
// 		deltas            []CampsiteDelta
// 		newAvailabilities map[string]api.Availability
// 	}
// 	var allDeltas []CampsiteDelta
// 	allNewAvails := make(map[string]api.Availability)
// 	var chunkWG, deltasWG sync.WaitGroup
// 	chunkOutputCHAN := make(chan chunkOutput)

// 	// batch := fs.Batch()

// 	deltasWG.Add(1)
// 	go func() {
// 		defer deltasWG.Done()
// 		for chunkOutput := range chunkOutputCHAN {
// 			allDeltas = append(allDeltas, chunkOutput.deltas...)
// 			// add all updates to the batch
// 			for ref, availability := range chunkOutput.newAvailabilities {
// 				allNewAvails[ref] = availability
// 			}
// 		}
// 	}()

// 	for i, chunk := range groundChunks {
// 		chunkWG.Add(1)
// 		go func(i int, chunk []string) {
// 			defer chunkWG.Done()
// 			log := log.With(zap.String("proxy", proxies[i]))
// 			proxy := fmt.Sprintf("https://%s-campr-app.cloudfunctions.net/HandleProxyRequest", proxies[i])
// 			newAvailabilities, chunkDeltas, err := CheckChunkForDeltas(ctx, log, fs, targetTime, now, chunk, proxy)
// 			if err != nil {
// 				log.Error("failed to check chunk", zap.Error(err))
// 				return
// 			}
// 			chunkOutputCHAN <- chunkOutput{
// 				deltas:            chunkDeltas,
// 				newAvailabilities: newAvailabilities,
// 			}
// 		}(i, chunk)
// 	}

// 	chunkWG.Wait()
// 	close(chunkOutputCHAN)
// 	deltasWG.Wait()

// 	// if there are no deltas at all, end here
// 	if len(allDeltas) == 0 {
// 		log.Info("checked availabilities, no deltas",
// 			zap.Duration("duration", time.Since(start)),
// 			zap.Int("deltas", len(allDeltas)),
// 		)
// 		return nil
// 	}
// 	updateStart := time.Now()

// 	// update new availabilities all at once
// 	for ref, availabiltity := range allNewAvails {
// 		_, err := fs.Collection(CollectionAvailability).Doc(ref).Set(ctx, availabiltity)
// 		if err != nil {
// 			log.Error("couldn't write new availability to firestore", zap.Error(err))
// 		}
// 	}

// 	commitEnd := time.Now()

// 	// update deltas in firestore.
// 	// add all deltas to the same document since campsites are globally unique so cheaper for consumers to grab the whole
// 	// document and search through for campsites you want.
// 	// breaking it up into chunks of 500 sites to avoid max document size limit
// 	chunkSize := 500
// 	startIndex := 0
// 	for startIndex < len(allDeltas) {
// 		endIndex := startIndex + chunkSize
// 		if endIndex > len(allDeltas) {
// 			endIndex = len(allDeltas)
// 		}

// 		checkDelta := CheckDelta{
// 			Deltas:    allDeltas[startIndex:endIndex],
// 			CheckTime: start,
// 		}
// 		ref, _, err := fs.Collection(CollectionAvailabilityDeltas).Add(
// 			ctx,
// 			checkDelta,
// 		)
// 		if err != nil {
// 			log.Error("couldn't add availability deltas to firestore", zap.Error(err))
// 			return err
// 		}
// 		log.Debug("synced chunk of deltas", zap.String("firestore_document", ref.ID))
// 		startIndex += chunkSize

// 	}
// 	log.Warn("synced deltas to firestore",
// 		zap.Int("delta_count", len(allDeltas)),
// 		zap.Duration("commit_duration", commitEnd.Sub(updateStart)),
// 		zap.Duration("delta_update_duration", time.Since(commitEnd)),
// 		zap.Duration("delta_retrieve_duration", updateStart.Sub(start)),
// 	)

// 	return nil
// }

type AvailabilityDetailed struct {
	Availability api.Availability
	Month        time.Time
	GroundID     string
}

func AvailabilitiesSync(ctx context.Context, log *zap.Logger, fs *firestore.Client, targetTime, now time.Time) ([]AvailabilityDetailed, error) {

	// start is used to measure duration and as the change time in the delta object
	start := time.Now()
	syncID := uuid.New().String()

	log = log.With(
		zap.String("sync_id", syncID),
		zap.Time("target_time", targetTime),
	)
	log.Info("started availability sync routine")

	// groundChunks := ChunkGroundsUp(proxies, campgroundIDs)

	var allDeltas []CampsiteDelta
	var allNewAvails []AvailabilityDetailed
	// allNewAvails := make(map[string]api.Availability)

	// allow up to five concurrent requests
	concurrencyLimiter := make(chan struct{}, 5)
	var mu sync.Mutex
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	rand.Seed(time.Now().UnixMicro())

	for _, targetGround := range campgroundIDs {
		wg.Add(1)
		go func(targetGround string) {
			defer wg.Done()
			// write to the buffered chan. as soon as it fills up, this will become blocking until other goroutines end.
			concurrencyLimiter <- struct{}{}
			defer func() { <-concurrencyLimiter }()

			randRegion := proxies[rand.Intn(len(proxies))]
			log := log.With(
				zap.String("proxy", randRegion),
				zap.String("ground_id", targetGround),
			)

			log.Debug("checking campground availability")
			proxyURI := fmt.Sprintf("https://%s-campr-app.cloudfunctions.net/HandleProxyRequest", randRegion)

			newAvailability, deltas, err := CheckForAvailabilityChange(ctx, log, proxyURI, fs, targetTime, now, targetGround)
			if err != nil {
				log.Error("failed to check availability change", zap.Error(err))
				cancel()
				return
			}

			// if there are no deltas, continue
			if len(deltas) == 0 {
				log.Debug("found no deltas")
				return
			}

			// if there are deltas, add the new availability and deltas to the returned lists
			log.Debug("deltas found", zap.Int("delta_count", len(deltas)))
			mu.Lock()
			allDeltas = append(allDeltas, deltas...)
			allNewAvails = append(allNewAvails, AvailabilityDetailed{
				Availability: newAvailability,
				Month:        api.GetStartOfMonth(targetTime),
				GroundID:     targetGround,
			})
			mu.Unlock()
		}(targetGround)
	}

	wg.Wait()

	// if there are no deltas at all, end here
	if len(allDeltas) == 0 {
		log.Info("checked availabilities, no deltas",
			zap.Duration("duration", time.Since(start)),
			zap.Int("deltas", len(allDeltas)),
		)
		return nil
	}

	updateStart := time.Now()
	log.Debug("finished scraping, uploading availabilities that have changed")

	// update new availabilities all at once
	for ref, availabiltity := range allNewAvails {
		_, err := fs.Collection(CollectionAvailability).Doc(ref).Set(ctx, availabiltity)
		if err != nil {
			log.Error("couldn't write new availability to firestore", zap.Error(err))
		}
	}

	commitEnd := time.Now()

	// update deltas in firestore.
	// add all deltas to the same document since campsites are globally unique so cheaper for consumers to grab the whole
	// document and search through for campsites you want.
	// breaking it up into chunks of 500 sites to avoid max document size limit
	chunkSize := 500
	startIndex := 0
	for startIndex < len(allDeltas) {
		endIndex := startIndex + chunkSize
		if endIndex > len(allDeltas) {
			endIndex = len(allDeltas)
		}

		checkDelta := CheckDelta{
			Deltas:    allDeltas[startIndex:endIndex],
			CheckTime: start,
		}
		ref, _, err := fs.Collection(CollectionAvailabilityDeltas).Add(
			ctx,
			checkDelta,
		)
		if err != nil {
			log.Error("couldn't add availability deltas to firestore", zap.Error(err))
			return err
		}
		log.Debug("synced chunk of deltas", zap.String("firestore_document", ref.ID))
		startIndex += chunkSize
	}

	log.Warn("synced deltas to firestore",
		zap.Int("delta_count", len(allDeltas)),
		zap.Duration("commit_duration", commitEnd.Sub(updateStart)),
		zap.Duration("delta_update_duration", time.Since(commitEnd)),
		zap.Duration("delta_retrieve_duration", updateStart.Sub(start)),
	)

	return nil
}

// func ChunkGroundsUp(proxies, campgroundIDs []string) [][]string {

// 	// shuffle this every time for extra obfuscation from cloudflare
// 	rand.Seed(time.Now().UnixNano())
// 	rand.Shuffle(len(campgroundIDs), func(i, j int) { campgroundIDs[i], campgroundIDs[j] = campgroundIDs[j], campgroundIDs[i] })

// 	chunks := make([][]string, len(proxies))
// 	for i, campgroundID := range campgroundIDs {
// 		index := i % len(proxies)
// 		chunks[index] = append(chunks[index], campgroundID)
// 	}

// 	return chunks
// }

// // Doing this async to minimise function operation time
// func CheckChunkForDeltasAsync(ctx context.Context, log *zap.Logger, fs *firestore.Client, targetTime, now time.Time, targetIDs []string, proxy string) (map[string]api.Availability, []CampsiteDelta, error) {

// 	// do it with a mutex for simplicity here
// 	var mu sync.Mutex
// 	var wg sync.WaitGroup

// 	// iterate over targets, looking for changes in availability for the given time
// 	var allDeltas []CampsiteDelta
// 	newAvailabilities := make(map[string]api.Availability)
// 	for _, targetGround := range targetIDs {
// 		wg.Add(1)
// 		go func(targetGround string) {
// 			defer wg.Done()
// 			log := log.With(zap.String("ground_id", targetGround))
// 			log.Debug("checking campground availability")
// 			newAvailability, deltas, err := CheckForAvailabilityChange(ctx, log, proxy, fs, targetTime, now, targetGround)
// 			if err != nil {
// 				log.Error("failed to check availability change", zap.Error(err))
// 				return
// 			}

// 			// if there are no deltas, continue
// 			if len(deltas) == 0 {
// 				log.Debug("found no deltas")
// 				return
// 			}

// 			// if there are deltas, add the new availability and deltas to the returned lists
// 			log.Debug("deltas found", zap.Int("delta_count", len(deltas)))
// 			mu.Lock()
// 			allDeltas = append(allDeltas, deltas...)
// 			newAvailabilities[GetAvailabilityRef(targetGround, targetTime)] = newAvailability
// 			mu.Unlock()
// 		}(targetGround)
// 	}

// 	wg.Wait()

// 	return newAvailabilities, allDeltas, nil
// }

// Doing this async to minimise function operation time
func CheckChunkForDeltas(ctx context.Context, log *zap.Logger, fs *firestore.Client, targetTime, now time.Time, targetIDs []string, proxy string) (map[string]api.Availability, []CampsiteDelta, error) {

	// iterate over targets, looking for changes in availability for the given time
	var allDeltas []CampsiteDelta
	newAvailabilities := make(map[string]api.Availability)
	for _, targetGround := range targetIDs {

		log := log.With(zap.String("ground_id", targetGround))
		log.Debug("checking campground availability")
		newAvailability, deltas, err := CheckForAvailabilityChange(ctx, log, proxy, fs, targetTime, now, targetGround)
		if err != nil {
			log.Error("failed to check availability change", zap.Error(err))
			return nil, nil, err
		}

		// if there are no deltas, continue
		if len(deltas) == 0 {
			log.Debug("found no deltas")
			return nil, nil, err
		}

		// if there are deltas, add the new availability and deltas to the returned lists
		log.Debug("deltas found", zap.Int("delta_count", len(deltas)))
		allDeltas = append(allDeltas, deltas...)
		newAvailabilities[GetAvailabilityRef(targetGround, targetTime)] = newAvailability

	}

	return newAvailabilities, allDeltas, nil
}
