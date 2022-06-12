package recreation

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/brensch/recreation/api"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	CollectionAvailabilities     = "availabilities"
	CollectionAvailabilityDeltas = "availability_deltas_grouped"

	SitesToSyncEachIteration = 20
)

// GetGroundsToScrape returns everything we want to scrape. Currently it grabs all site ids,
// then selects SitesToSyncEachIteration site ids from that list at random
func GetGroundIDsToScrape(ctx context.Context, fs *firestore.Client) ([]string, error) {
	snap, err := fs.Collection(CollectionGroundsSummary).Doc(DocGroundsSummary).Get(ctx)
	if err != nil {
		return nil, err
	}

	var summary GroundSummary
	err = snap.DataTo(&summary)
	if err != nil {
		return nil, err
	}

	ids := summary.GroundIDs
	return SelectRandomIDs(ids), nil
}

// SelectRandomIDs picks 5 sites at random, ensuring they're all different
func SelectRandomIDs(input []string) []string {
	var randomIDs []string
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < SitesToSyncEachIteration; i++ {
		idToRemove := rand.Intn(len(input))

		// add the chosen id to the array
		randomIDs = append(randomIDs, input[idToRemove])

		// then remove it from the array by replacing and truncating
		input[idToRemove] = input[len(input)-1]
		input = input[:len(input)-1]
	}

	return randomIDs
}

// GetAvailabilityRef gives us a consisten ID to work with for our documents
func GetAvailabilityRef(groundID string, targetTime time.Time) string {
	targetTime = api.GetStartOfMonth(targetTime)
	return fmt.Sprintf("%s-%s", groundID, targetTime.Format(time.RFC3339))
}

// CheckForAvailabilityChange gets old and new states of availability and returns the new availability and deltas to old
func CheckForAvailabilityChange(ctx context.Context, log *zap.Logger, baseURI string, fs *firestore.Client, targetTime time.Time, targetGround string) (api.Availability, []CampsiteDelta, error) {

	// get new availability from API
	newAvailability, err := api.GetAvailability(ctx, log, baseURI, targetGround, targetTime)
	if err != nil {
		return api.Availability{}, nil, err
	}

	// get old abailability from firestore
	// NotFound errors ignored since the document not existing just results in an empty object, as intended
	oldAvailabilitySnap, err := fs.Collection(CollectionAvailabilities).Doc(GetAvailabilityRef(targetGround, targetTime)).Get(ctx)
	if err != nil && status.Code(err) != codes.NotFound {
		return api.Availability{}, nil, err
	}
	var oldAvailability api.Availability
	err = oldAvailabilitySnap.DataTo(&oldAvailability)
	if err != nil && status.Code(err) != codes.NotFound {
		return api.Availability{}, nil, err
	}

	// compare the old and new availabilities
	deltas, err := FindAvailabilityDeltas(oldAvailability, newAvailability)
	return newAvailability, deltas, err

}

func DoAvailabilitiesSync(ctx context.Context, log *zap.Logger, baseURI string, fs *firestore.Client, targetTime time.Time) error {

	// start is used to measure duration and as the change time in the delta object
	start := time.Now()
	syncID := uuid.New().String()

	log = log.With(
		zap.String("sync_id", syncID),
		zap.Time("target_time", targetTime),
	)
	log.Info("started availability sync routine")

	// select targets
	targetGrounds, err := GetGroundIDsToScrape(ctx, fs)
	if err != nil {
		log.Error("failed to get ground IDs to scrape", zap.Error(err))
		return err
	}

	// iterate over targets, looking for changes in availability for the given time
	var allDeltas []CampsiteDelta
	for _, targetGround := range targetGrounds {
		log := log.With(zap.String("ground_id", targetGround))
		log.Debug("checking campground availability")
		newAvailability, deltas, err := CheckForAvailabilityChange(ctx, log, baseURI, fs, targetTime, targetGround)
		if err != nil {
			log.Error("failed to check availability change", zap.Error(err))
			return err
		}

		// if there are no deltas, continue
		if len(deltas) == 0 {
			log.Debug("found no deltas, continuing")
			continue
		}

		log.Debug("deltas found", zap.Int("delta_count", len(deltas)))
		allDeltas = append(allDeltas, deltas...)

		// update availabilities
		_, err = fs.Collection(CollectionAvailabilities).Doc(GetAvailabilityRef(targetGround, targetTime)).Set(
			ctx,
			newAvailability,
		)
		if err != nil {
			log.Error("couldn't add availability to firestore", zap.Error(err))
			return err
		}
	}

	log.Info("completed availability sync routine", zap.Duration("duration", time.Since(start)))

	// if there are no deltas at all, end here
	if len(allDeltas) == 0 {
		return nil
	}

	// update deltas in firestore.
	// add all deltas to the same document since campsites are globally unique so cheaper for consumers to grab the whole
	// document and search through for campsites you want.
	log.Warn("syncing deltas to firestore", zap.Int("delta_count", len(allDeltas)))
	checkDelta := CheckDelta{
		Deltas:    allDeltas,
		CheckTime: start,
	}
	_, _, err = fs.Collection(CollectionAvailabilityDeltas).Add(
		ctx,
		checkDelta,
	)
	if err != nil {
		log.Error("couldn't add availability deltas to firestore", zap.Error(err))
		return err
	}

	return nil
}
