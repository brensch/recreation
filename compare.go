package main

import (
	"time"

	"github.com/brensch/recreation/api"
	"go.uber.org/zap"
)

type CampsiteDelta struct {
	SiteID       string
	GroundID     string
	OldState     string
	NewState     string
	DateAffected time.Time
}

type CampsiteDelta2 struct {
	SiteID       string
	GroundID     string
	OldState     AvailabilityState
	NewState     AvailabilityState
	DateAffected time.Time
}

type CheckDelta struct {
	// GroundID  string
	Deltas    []CampsiteDelta
	CheckTime time.Time
}

// FindAvailabilityDeltas compares old and new availability and returns all deltas between the two
func FindAvailabilityDeltas(log *zap.Logger, oldGround, newGround api.Availability, groundID string, now time.Time) ([]CampsiteDelta, error) {

	var deltas []CampsiteDelta

	// iterate through each field in new and check what the previous value was
	for siteID, newSite := range newGround.Campsites {
		oldSite := oldGround.Campsites[siteID]
		for dateString, availability := range newSite.Availabilities {

			// ignore things that haven't changed.
			// using a map here is nice, i think it's efficient. May try other approaches if i get frisky
			if oldSite.Availabilities[dateString] == availability {
				continue
			}

			date, err := time.Parse(time.RFC3339, dateString)
			if err != nil {
				log.Error("failed to parse date in availability map", zap.Error(err))
				return nil, err
			}

			// ignore dates in the past. the api reports them inconsistently, plus we're not interested in them.
			if date.Before(now) {
				continue
			}

			deltas = append(deltas, CampsiteDelta{
				SiteID:       siteID,
				GroundID:     groundID,
				OldState:     oldSite.Availabilities[dateString],
				NewState:     availability,
				DateAffected: date,
			})
		}

	}

	return deltas, nil

}
