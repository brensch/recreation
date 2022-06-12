package recreation

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/brensch/recreation/api"
	"go.uber.org/zap"
)

func TestScrapeGrounds(t *testing.T) {

	ctx := context.Background()
	log, _ := zap.NewProduction()

	grounds, err := ScrapeGrounds(ctx, log, "http://www.recreation.gov")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	t.Log("got ground count", len(grounds.Results))

}

func TestUpdateGrounds(t *testing.T) {

	ctx := context.Background()
	// TODO: use a local firestore instance for this
	fs, err := InitFirestore(ctx)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer fs.Close()

	var grounds []api.CampGround
	err = json.Unmarshal(groundTests, &grounds)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	err = UpdateGrounds(ctx, fs, grounds)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
}

func TestInitGroundsSummary(t *testing.T) {
	ctx := context.Background()

	fs, err := InitFirestore(ctx)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer fs.Close()

	var grounds []api.CampGround
	err = json.Unmarshal(groundTests, &grounds)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	err = InitGroundsSummary(ctx, fs)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

}

func TestUpdateGroundsSummary(t *testing.T) {
	ctx := context.Background()

	fs, err := InitFirestore(ctx)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer fs.Close()

	var grounds []api.CampGround
	err = json.Unmarshal(groundTests, &grounds)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	err = UpdateGroundsSummary(ctx, fs, grounds)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	snap, err := fs.Collection(CollectionGroundsSummary).Doc(DocGroundsSummary).Get(ctx)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	// check total length is as expected
	// NB this will fail if using the same collection as production
	var summary GroundSummary
	err = snap.DataTo(&summary)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if len(summary.GroundIDs) != 3 {
		t.Errorf("expecting 3 ids, got %d", len(summary.GroundIDs))
	}
}

var (
	groundTests = []byte(`
	[
        {
            "accessible_campsites_count": 12,
            "activities": null,
            "addresses": null,
            "aggregate_cell_coverage": 3.545454588803378,
            "average_rating": 4.8,
            "campsite_accessible": 1,
            "campsite_reserve_type": [
                "Site-Specific"
            ],
            "campsite_type_of_use": [
                "Overnight"
            ],
            "campsites_count": "12",
            "city": "San Francisco",
            "country_code": "United States",
            "description": "Aquatic Park Cove is a vibrantly urban harbor that mimics a natural anchorage and is located on the Pacific shores of one of America's great West Coast cities.  San Francisco Maritime National Historical Park's Aquatic Park Cove is the small craft ga",
            "directions": "",
            "distance": "55.35",
            "entity_id": "273757",
            "entity_type": "campground",
            "go_live_date": "2018-09-30T00:00:00Z",
            "html_description": "",
            "id": "273757_asset",
            "latitude": "37.80666670000000",
            "links": null,
            "longitude": "-122.42388889999999",
            "name": "Aquatic Park Cove, Black Point Overnight Anchoring, Sailing Vessels only (auxiliary engine OK)",
            "notices": null,
            "number_of_ratings": 20,
            "org_id": "128",
            "org_name": "National Park Service",
            "parent_id": "2915",
            "parent_name": "San Francisco Maritime National Historical Park",
            "parent_type": "",
            "preview_image_url": "https://cdn.api.gov/public/2021/04/27/21/34/273757_999a6139-3cc0-45e4-a274-1e6a36242dfa_700.jpg",
            "price_range": {
                "amount_max": 0,
                "amount_min": 0,
                "per_unit": ""
            },
            "rate": null,
            "reservable": true,
            "state_code": "California",
            "type": "STANDARD"
        },
        {
            "accessible_campsites_count": 2,
            "activities": null,
            "addresses": null,
            "aggregate_cell_coverage": 4,
            "average_rating": 5,
            "campsite_accessible": 1,
            "campsite_reserve_type": [
                "Site-Specific"
            ],
            "campsite_type_of_use": [
                "Overnight"
            ],
            "campsites_count": "4",
            "city": "San Francisco",
            "country_code": "United States",
            "description": "As San Francisco's only group campground, Rob Hill offers a national park camping experience just minutes from the city. Here you will have access to all of the recreational opportunities the Presidio of San Francisco has to offer, including hiking a",
            "directions": "",
            "distance": "57.10",
            "entity_id": "10172170",
            "entity_type": "campground",
            "go_live_date": "2022-02-23T15:48:30.843Z",
            "html_description": "",
            "id": "10172170_asset",
            "latitude": "37.79753700000000",
            "links": null,
            "longitude": "-122.47562300000000",
            "name": "Rob Hill Group Campground",
            "notices": null,
            "number_of_ratings": 2,
            "org_id": "250",
            "org_name": "Presidio Trust",
            "parent_id": "RA1014",
            "parent_name": "Presidio of San Francisco",
            "parent_type": "",
            "preview_image_url": "https://cdn.api.gov/public/2022/02/24/21/33/10172170_ba4e1946-4c5a-44be-a92c-f68a9c7bc0e1_700.jpg",
            "price_range": {
                "amount_max": 0,
                "amount_min": 0,
                "per_unit": ""
            },
            "rate": null,
            "reservable": true,
            "state_code": "California",
            "type": "STANDARD"
        },
        {
            "accessible_campsites_count": 3,
            "activities": null,
            "addresses": null,
            "aggregate_cell_coverage": 3.3114754567380813,
            "average_rating": 4.6,
            "campsite_accessible": 1,
            "campsite_reserve_type": [
                "Site-Specific"
            ],
            "campsite_type_of_use": [
                "Overnight",
                "Day"
            ],
            "campsites_count": "6",
            "city": "Sausalito",
            "country_code": "United States",
            "description": "Kirby Cove is located just north of the Golden Gate Bridge at historic Battery Kirby. Visitors are awarded breathtaking views of San Francisco, its famous Golden Gate Bridge, and the rugged Pacific Coast of northern California. \nThe San Francisco Bay",
            "directions": "",
            "distance": "61.64",
            "entity_id": "232491",
            "entity_type": "campground",
            "go_live_date": "2018-09-30T00:00:00Z",
            "html_description": "",
            "id": "232491_asset",
            "latitude": "37.84035000000000",
            "links": null,
            "longitude": "-122.48888890000001",
            "name": "KIRBY COVE CAMPGROUND",
            "notices": null,
            "number_of_ratings": 75,
            "org_id": "128",
            "org_name": "National Park Service",
            "parent_id": "2730",
            "parent_name": "Golden Gate National Recreation Area",
            "parent_type": "",
            "preview_image_url": "https://cdn.api.gov/public/2020/04/08/05/22/232491_de74dcb9-2a90-427c-a267-a01c65c7e114_700.jpg",
            "price_range": {
                "amount_max": 0,
                "amount_min": 0,
                "per_unit": ""
            },
            "rate": null,
            "reservable": true,
            "state_code": "California",
            "type": "STANDARD"
        }
	]
`)
)
