package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type State string

var (
	StateAvailable               State = "Available"
	StateReserved                State = "Reserved"
	StateNotReservableManagement State = "Not Reservable Management"
)

type Availability struct {
	Campsites map[string]Campsite `json:"campsites,omitempty"`
	Count     int                 `json:"count,omitempty"`
}

type Campsite struct {
	// keeping this as a string even though it's a time object for less processing
	Availabilities map[string]string `json:"availabilities"`

	CampsiteID          string `json:"campsite_id"`
	CampsiteReserveType string `json:"campsite_reserve_type"`
	CampsiteType        string `json:"campsite_type"`
	CapacityRating      string `json:"capacity_rating"`
	Loop                string `json:"loop"`
	MaxNumPeople        int    `json:"max_num_people"`
	MinNumPeople        int    `json:"min_num_people"`
	Site                string `json:"site"` // not sure what this represents
	TypeOfUse           string `json:"type_of_use"`

	// TODO: find example of this. haven't seen what form it takes yet.
	CampsiteRules interface{} `json:"campsite_rules"`

	// not sure what quantities means
	// TODO: figure out if we need it
	Quantities struct{} `json:"quantities"`
}

// GetAvailability ensures that the targettime is snapped to the start of the month, then queries the API for all availabilities at that ground
func GetAvailability(ctx context.Context, log *zap.Logger, baseURI string, campgroundID string, targetTime time.Time) (Availability, error) {

	start := time.Now()
	log.Debug("getting availability from api")
	endpoint := fmt.Sprintf("%s/api/camps/availability/campground/%s/month", baseURI, campgroundID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		log.Error("couldn't create request", zap.Error(err))
		return Availability{}, err
	}

	// round the time to the start of the target month and put in param "start_date"
	monthStart := GetStartOfMonth(targetTime)

	// params need to be url encoded. ie base64
	v := req.URL.Query()
	v.Add("start_date", monthStart.Format("2006-01-02T15:04:05.000Z"))
	req.URL.RawQuery = v.Encode()

	res, err := client.Do(req)
	if err != nil {
		log.Error("couldn't do request", zap.Error(err))
		return Availability{}, err
	}
	defer res.Body.Close()

	resContents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error("couldn't read response", zap.Error(err))
		return Availability{}, err
	}

	// fmt.Println(string(resContents))

	if res.StatusCode != http.StatusOK {
		log.Error("got bad statuscode getting availability", zap.Int("status_code", res.StatusCode))
		log.Debug("body of bad request", zap.String("body", string(resContents)))
		return Availability{}, fmt.Errorf(string(resContents))
	}

	var availability Availability
	err = json.Unmarshal(resContents, &availability)
	if err != nil {
		log.Error("couldn't unmarshal", zap.Error(err))
		return Availability{}, err
	}

	log.Debug("completed getting availability from api", zap.Duration("duration", time.Since(start)))

	return availability, nil

}
