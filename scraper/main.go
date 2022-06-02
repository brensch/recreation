package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/brensch/recreation"
	"google.golang.org/api/option"
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

	ctx := context.Background()
	client := recreation.InitObfuscator(ctx)
	_ = client
	opt := option.WithCredentialsFile("creds.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		fmt.Printf("error initializing app: %v", err)
		return
	}

	fs, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer fs.Close()
	fmt.Println("starting")

	log.Print("starting server...")
	http.HandleFunc("/", HandleDoSite(client, fs))

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}

	// GetAllAvailabilities(ctx, client)

}

func HandleDoSite(client *recreation.Obfuscator, fs *firestore.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("got request")
		DoSites(r.Context(), client, fs)
	}
}

func DoSites(ctx context.Context, client *recreation.Obfuscator, fs *firestore.Client) {
	// get sites
	res, err := recreation.DoSearchGeo(ctx, client, 37.3859, -122.0882)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, campGround := range res.Results {
		var avails []recreation.Availability
		for i := 0; i < 5; i++ {
			targetTime := time.Date(time.Now().Year(), time.Now().Month()+time.Month(i), 1, 0, 0, 0, 0, time.UTC)
			availability, err := recreation.GetAvailability(ctx, client, Park, targetTime)
			if err != nil {
				fmt.Println(err)
				return
			}
			avails = append(avails, availability)
		}

		sacredGround := AvailabilitiesToSacredGround(avails)

		currentSacredGroundSnap, err := fs.Collection("grounds").Doc(campGround.EntityID).Get(context.Background())
		if err != nil {
			fmt.Println(err)
			return
		}

		var currentSacredGround SacredGround
		err = currentSacredGroundSnap.DataTo(&currentSacredGround)
		if err != nil {
			fmt.Println(err)
			return
		}

		_, err = fs.Collection("grounds").Doc(campGround.EntityID).Set(context.Background(), sacredGround)
		if err != nil {
			fmt.Println(err)
			return
		}
		return
	}
}

// SacredGround is a silly name meaning a campground in the format that i want it.
// the idea is that i will group every site into the different states present to help doing deltas
type SacredGround struct {
	Sites map[string]Site
}

type siteState int

const (
	stateAvailable siteState = iota
	stateReserved
	stateNotReservableManagement
)

// TODO: map states to enums
// var (
// 	stateMappings =
// )

type Site struct {
	// ID string
	// this is actually a date as the key
	// todo make values sitestates
	Availabilities map[string]string
}

func AvailabilitiesToSacredGround(avails []recreation.Availability) SacredGround {
	s := SacredGround{
		Sites: make(map[string]Site),
	}

	for _, avail := range avails {

		for _, site := range avail.Campsites {

			// TODO: trim times for space savings

			_, ok := s.Sites[site.CampsiteID]

			if !ok {
				s.Sites[site.CampsiteID] = Site{
					Availabilities: site.Availabilities,
				}
				continue
			}

			// if it exists, have to iterate through map and add. bit ugly.
			for date, state := range site.Availabilities {
				s.Sites[site.CampsiteID].Availabilities[date] = state
			}
		}
	}

	return s
}

// func GetAllAvailabilities(ctx context.Context, client recreation.HTTPClient) {

// 	res, err := recreation.DoSearchGeo(ctx, client, 37.3859, -122.0882)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	fmt.Println(len(res.Results))

// 	for _, campground := range res.Results[0:50] {
// 		fmt.Println(campground.EntityID, campground.Name, campground.City)

// 		daysFree := 0
// 		for i := 0; i < 6; i++ {

// 			targetTime := time.Now()
// 			targetTime = time.Date(targetTime.Year(), targetTime.Month()+time.Month(i), 1, 0, 0, 0, 0, time.UTC)

// 			availability, err := recreation.GetAvailability(ctx, client, Park, targetTime)
// 			if err != nil {
// 				fmt.Println(err)
// 				return
// 			}
// 			for _, camp := range availability.Campsites {
// 				for _, avail := range camp.Availabilities {
// 					if avail == recreation.StateAvailable {
// 						daysFree++
// 					}
// 				}
// 			}
// 		}
// 		fmt.Println(daysFree)
// 	}
// }

// func GetDailyAvailabilities(ctx context.Context, client recreation.HTTPClient) {
// 	allAvailabilities := make(map[string][]time.Time)

// 	for i := 0; i < 6; i++ {

// 		targetTime := time.Now()
// 		targetTime = time.Date(targetTime.Year(), targetTime.Month()+time.Month(i), 1, 0, 0, 0, 0, time.UTC)

// 		availability, err := recreation.GetAvailability(ctx, client, Park, targetTime)
// 		if err != nil {
// 			fmt.Println(err)
// 			return
// 		}

// 		for campsiteID, camp := range availability.Campsites {
// 			for date, avail := range camp.Availabilities {
// 				if avail == recreation.StateAvailable {
// 					allAvailabilities[campsiteID] = append(allAvailabilities[campsiteID], date)
// 				}
// 			}

// 		}
// 	}
// 	campsiteAvailability, ok := allAvailabilities[Campsite]
// 	if !ok {
// 		fmt.Println("didn't find the campsite you specified, getrekt")
// 	}

// 	fmt.Printf("availabilities for the next five months at %s:%s\n", Park, Campsite)
// 	sort.Slice(campsiteAvailability, func(i, j int) bool {
// 		return campsiteAvailability[i].Before(campsiteAvailability[j])
// 	})
// 	for _, availableDate := range campsiteAvailability {
// 		fmt.Println(availableDate.Format("2006-01-02 Monday"))
// 	}
// }
