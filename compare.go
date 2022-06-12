package recreation

import (
	"time"

	"github.com/brensch/recreation/api"
)

type CampsiteDelta struct {
	SiteID       string
	OldState     string
	NewState     string
	DateAffected time.Time
}

type CheckDelta struct {
	// GroundID  string
	Deltas    []CampsiteDelta
	CheckTime time.Time
}

func FindAvailabilityDeltas(oldGround, newGround api.Availability) ([]CampsiteDelta, error) {

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
				return nil, err
			}

			deltas = append(deltas, CampsiteDelta{
				SiteID:       siteID,
				OldState:     oldSite.Availabilities[dateString],
				NewState:     availability,
				DateAffected: date,
			})
		}

	}

	return deltas, nil

}
