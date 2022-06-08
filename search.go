package recreation

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// I found two ways to search:
// api/search/geo - gives only geo based results but seems to return far more results
// api/search - gives everything

// Their search capability is actually pretty powerful. Seems like you just yeet any field name
// into a param with key 'fq' followed by %3A and the value you're searching for.
// eg: entity_type%3Acampground

type SearchResults struct {
	Latitude              string       `json:"latitude"`
	Location              string       `json:"location"`
	Longitude             string       `json:"longitude"`
	Radius                string       `json:"radius"`
	Results               []CampGround `json:"results"`
	Size                  int          `json:"size"`
	SpellingAutocorrected bool         `json:"spelling_autocorrected"`
	Start                 string       `json:"start"`
	Total                 int          `json:"total"`
}

type CampGround struct {
	AccessibleCampsitesCount int `json:"accessible_campsites_count,omitempty"`
	Activities               []struct {
		ActivityDescription    string `json:"activity_description"`
		ActivityFeeDescription string `json:"activity_fee_description"`
		ActivityID             int    `json:"activity_id"`
		ActivityName           string `json:"activity_name"`
	} `json:"activities"`
	Addresses []struct {
		AddressType    string `json:"address_type"`
		City           string `json:"city"`
		CountryCode    string `json:"country_code"`
		PostalCode     string `json:"postal_code"`
		StateCode      string `json:"state_code"`
		StreetAddress1 string `json:"street_address1"`
		StreetAddress2 string `json:"street_address2"`
		StreetAddress3 string `json:"street_address3"`
	} `json:"addresses"`
	AggregateCellCoverage float64   `json:"aggregate_cell_coverage,omitempty"`
	AverageRating         float64   `json:"average_rating,omitempty"`
	CampsiteAccessible    int       `json:"campsite_accessible,omitempty"`
	CampsiteEquipmentName []string  `json:"campsite_equipment_name,omitempty"`
	CampsiteReserveType   []string  `json:"campsite_reserve_type"`
	CampsiteTypeOfUse     []string  `json:"campsite_type_of_use"`
	CampsitesCount        string    `json:"campsites_count"`
	City                  string    `json:"city"`
	CountryCode           string    `json:"country_code"`
	Description           string    `json:"description"`
	Directions            string    `json:"directions"`
	Distance              string    `json:"distance"`
	EntityID              string    `json:"entity_id"`
	EntityType            string    `json:"entity_type"`
	GoLiveDate            time.Time `json:"go_live_date"`
	HTMLDescription       string    `json:"html_description"`
	ID                    string    `json:"id"`
	Latitude              string    `json:"latitude"`
	Links                 []struct {
		Description string `json:"description"`
		LinkType    string `json:"link_type"`
		Title       string `json:"title"`
		URL         string `json:"url"`
	} `json:"links"`
	Longitude string `json:"longitude"`
	Name      string `json:"name"`
	Notices   []struct {
		Text string `json:"text"`
		Type string `json:"type"`
	} `json:"notices"`
	NumberOfRatings int    `json:"number_of_ratings,omitempty"`
	OrgID           string `json:"org_id"`
	OrgName         string `json:"org_name"`
	ParentID        string `json:"parent_id"`
	ParentName      string `json:"parent_name"`
	ParentType      string `json:"parent_type"`
	PreviewImageURL string `json:"preview_image_url"`
	PriceRange      struct {
		AmountMax int    `json:"amount_max"`
		AmountMin int    `json:"amount_min"`
		PerUnit   string `json:"per_unit"`
	} `json:"price_range,omitempty"`
	Rate []struct {
		EndDate time.Time `json:"end_date"`
		Prices  []struct {
			Amount    int    `json:"amount"`
			Attribute string `json:"attribute"`
		} `json:"prices"`
		RateMap map[string]struct {
			GroupFees        interface{} `json:"group_fees"`
			SingleAmountFees Fees        `json:"single_amount_fees"`
		} `json:"rate_map"`
		SeasonDescription string    `json:"season_description"`
		SeasonType        string    `json:"season_type"`
		StartDate         time.Time `json:"start_date"`
	} `json:"rate"`
	Reservable bool   `json:"reservable"`
	StateCode  string `json:"state_code"`
	TimeZone   string `json:"time_zone,omitempty"`
	Type       string `json:"type"`
}

type Fees struct {
	Deposit   int `json:"deposit"`
	Holiday   int `json:"holiday"`
	PerNight  int `json:"per_night"`
	PerPerson int `json:"per_person"`
	Weekend   int `json:"weekend"`
}

func (s *Server) SearchGeo(ctx context.Context, log *zap.Logger, lat, lon float64) (SearchResults, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return searchGeo(ctx, log, s.client, lat, lon)
}

func searchGeo(ctx context.Context, log *zap.Logger, client HTTPClient, lat, lon float64) (SearchResults, error) {

	start := time.Now()
	log = log.With(
		zap.Float64("lat", lat),
		zap.Float64("lon", lon),
	)
	log.Debug("doing search using api")
	endpoint := fmt.Sprintf("%s/api/search/geo", RecreationGovURI)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		log.Error("couldn't create request", zap.Error(err))
		return SearchResults{}, err
	}

	v := req.URL.Query()

	v.Add("lat", fmt.Sprint(lat))
	v.Add("lng", fmt.Sprint(lon))

	// TODO: maybe make these customizable
	v.Add("exact", "false")
	v.Add("size", "1000")
	// v.Add("fq", `-entity_type%3A(tour%20OR%20timedentry_tour)`)
	// v.Add("fq", "campsite_type_of_use%3AOvernight")
	// v.Add("fq", "campsite_type_of_use%3Ana")
	v.Add("fq", "entity_type%3Acampground")
	// v.Add("fq", "campsite_type_of_use%3ADay")

	req.URL.RawQuery = v.Encode()

	res, err := client.Do(req)
	if err != nil {
		log.Error("couldn't do request", zap.Error(err))
		return SearchResults{}, err
	}
	defer res.Body.Close()

	// doing a readall since cloudflare dumps xml on you
	resContents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error("couldn't read response", zap.Error(err))
		return SearchResults{}, err
	}

	if res.StatusCode != http.StatusOK {
		log.Error("got bad statuscode searching geo",
			zap.Int("status_code", res.StatusCode),
			zap.String("body", string(resContents)),
		)
		return SearchResults{}, fmt.Errorf(string(resContents))
	}

	var results SearchResults
	err = json.Unmarshal(resContents, &results)
	if err != nil {
		log.Error("couldn't unmarshal", zap.Error(err))
		return SearchResults{}, err
	}

	log.Debug("completed search using api", zap.Duration("duration", time.Since(start)))

	return results, nil

}
