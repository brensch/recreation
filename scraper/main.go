package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/brensch/recreation"
)

var (
	Park     = "232770"
	Campsite = "86766"
)

func init() {
	flag.StringVar(&Park, "park", Park, "Which park to check")
	flag.StringVar(&Campsite, "campsite", Campsite, "Which campsite to check")

}

func main() {
	flag.Parse()

	client := http.DefaultClient
	ctx := context.Background()
	allAvailabilities := make(map[string][]time.Time)

	for i := 0; i < 6; i++ {

		targetTime := time.Now()
		targetTime = time.Date(targetTime.Year(), targetTime.Month()+time.Month(i), 1, 0, 0, 0, 0, time.UTC)

		availability, err := recreation.GetAvailability(ctx, client, Park, targetTime)
		if err != nil {
			fmt.Println(err)
			return
		}

		for campsiteID, camp := range availability.Campsites {
			for date, avail := range camp.Availabilities {
				if avail == recreation.StateAvailable {
					allAvailabilities[campsiteID] = append(allAvailabilities[campsiteID], date)
				}
			}

		}
	}
	campsiteAvailability, ok := allAvailabilities[Campsite]
	if !ok {
		fmt.Println("didn't find the campsite you specified, getrekt")
	}

	fmt.Printf("availabilities for the next five months at %s:%s\n", Park, Campsite)
	sort.Slice(campsiteAvailability, func(i, j int) bool {
		return campsiteAvailability[i].Before(campsiteAvailability[j])
	})
	for _, availableDate := range campsiteAvailability {
		fmt.Println(availableDate.Format("2006-01-02 Monday"))
	}
}
