package recreation

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type State string

var (
	StateAvailable               State = "Available"
	StateReserved                State = "Reserved"
	StateNotReservableManagement State = "Not Reservable Management"
)

type Availability struct {
	Campsites map[string]Campsite
	Count     int
}

type Campsite struct {
	Availabilities map[time.Time]State `json:"availabilities"`

	CampsiteID          string      `json:"campsite_id"`
	CampsiteReserveType string      `json:"campsite_reserve_type"`
	CampsiteRules       interface{} `json:"campsite_rules"` // TODO: find example of this
	CampsiteType        string      `json:"campsite_type"`
	CapacityRating      string      `json:"capacity_rating"`
	Loop                string      `json:"loop"`
	MaxNumPeople        int         `json:"max_num_people"`
	MinNumPeople        int         `json:"min_num_people"`
	Site                string      `json:"site"`
	TypeOfUse           string      `json:"type_of_use"`

	// not sure what quantities means
	// TODO: figure out if we need it
	Quantities struct{} `json:"quantities"`
}

func GetAvailability(ctx context.Context, client *http.Client, site string, targetTime time.Time) (Availability, error) {

	endpoint := fmt.Sprintf("%s/api/camps/availability/campground/%s/month", RecreationGovURI, site)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return Availability{}, err
	}

	// round the time to the start of the target month and put in param "start_date"
	monthStart := GetStartOfMonth(targetTime)
	v := req.URL.Query()
	v.Add("start_date", monthStart.Format("2006-01-02T15:04:05.000Z"))
	req.URL.RawQuery = v.Encode()

	// set user agent since it seems cloudfront is blocking go's default one
	req.Header.Set("User-Agent", UserAgent)

	res, err := client.Do(req)
	if err != nil {
		return Availability{}, err
	}
	defer res.Body.Close()

	var availability Availability
	err = json.NewDecoder(res.Body).Decode(&availability)
	if err != nil {
		return Availability{}, err
	}

	return availability, nil

}
