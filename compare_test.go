package main

import (
	"testing"
	"time"

	"github.com/brensch/recreation/api"
	"go.uber.org/zap"
)

func TestCompareCampgroundStates(t *testing.T) {

	date, err := time.Parse(time.RFC3339, "2022-06-01T00:00:00Z")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	old := api.Availability{
		Campsites: map[string]api.Campsite{
			"site_1": {
				Availabilities: map[string]string{
					"2022-06-01T00:00:00Z": "Not Reservable",
					"2022-06-02T00:00:00Z": "Not Reservable",
					"2022-06-03T00:00:00Z": "Available",
					"2022-06-04T00:00:00Z": "Available",
					"2022-06-05T00:00:00Z": "Available",
					"2022-06-06T00:00:00Z": "Not Reservable",
					"2022-06-07T00:00:00Z": "Available",
					"2022-06-08T00:00:00Z": "Available",
					"2022-06-09T00:00:00Z": "Available",
					"2022-06-10T00:00:00Z": "Available",
				},
			},
		},
	}

	new := api.Availability{
		Campsites: map[string]api.Campsite{
			"site_1": {
				Availabilities: map[string]string{
					"2022-06-01T00:00:00Z": "Not Reservable",
					"2022-06-02T00:00:00Z": "Not Reservable",
					"2022-06-03T00:00:00Z": "Available",
					"2022-06-04T00:00:00Z": "Available",
					"2022-06-05T00:00:00Z": "Available",
					"2022-06-06T00:00:00Z": "Available",
					"2022-06-07T00:00:00Z": "Available",
					"2022-06-08T00:00:00Z": "Available",
					"2022-06-09T00:00:00Z": "Available",
					"2022-06-10T00:00:00Z": "Not Reservable",
				},
			},
		},
	}

	log, _ := zap.NewDevelopment()
	deltas, err := FindAvailabilityDeltas(log, old, new, "testGround", date)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if len(deltas) != 2 {
		t.Error("got wrong number of deltas")
	}
	t.Log(deltas)

}

func BenchmarkCompareCampgroundStates(b *testing.B) {

	date, err := time.Parse(time.RFC3339, "2022-06-01T00:00:00Z")
	if err != nil {
		b.Error(err)
		b.FailNow()
	}

	old := api.Availability{
		Campsites: map[string]api.Campsite{
			"site_1": {
				Availabilities: map[string]string{
					"2022-06-01T00:00:00Z": "Not Reservable",
					"2022-06-02T00:00:00Z": "Not Reservable",
					"2022-06-03T00:00:00Z": "Available",
					"2022-06-04T00:00:00Z": "Available",
					"2022-06-05T00:00:00Z": "Available",
					"2022-06-06T00:00:00Z": "Not Reservable",
					"2022-06-07T00:00:00Z": "Available",
					"2022-06-08T00:00:00Z": "Available",
					"2022-06-09T00:00:00Z": "Available",
					"2022-06-10T00:00:00Z": "Available",
				},
			},
		},
	}

	new := api.Availability{
		Campsites: map[string]api.Campsite{
			"site_1": {
				Availabilities: map[string]string{
					"2022-06-01T00:00:00Z": "Not Reservable",
					"2022-06-02T00:00:00Z": "Not Reservable",
					"2022-06-03T00:00:00Z": "Available",
					"2022-06-04T00:00:00Z": "Available",
					"2022-06-05T00:00:00Z": "Available",
					"2022-06-06T00:00:00Z": "Available",
					"2022-06-07T00:00:00Z": "Available",
					"2022-06-08T00:00:00Z": "Available",
					"2022-06-09T00:00:00Z": "Available",
					"2022-06-10T00:00:00Z": "Not Reservable",
				},
			},
		},
	}

	log, _ := zap.NewDevelopment()

	for i := 0; i < b.N; i++ {
		FindAvailabilityDeltas(log, old, new, "test_ground", date)
	}
}
