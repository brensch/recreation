package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"go.uber.org/zap"
)

func main() {

	// init logger
	logConfig := zap.NewProductionConfig()
	logConfig.Level.SetLevel(zap.DebugLevel)
	// this ensures google logs pick things up properly
	logConfig.EncoderConfig.MessageKey = "message"
	logConfig.EncoderConfig.LevelKey = "severity"
	logConfig.EncoderConfig.TimeKey = "time"
	log, err := logConfig.Build()
	if err != nil {
		// this indicates a bug or some way that zap can fail i'm not aware of
		panic(err)
	}

	// init firestore
	fs, err := InitFirestore(context.Background())
	if err != nil {
		panic(err)
	}

	// // this deletes everything from firestore
	// colls := []string{
	// 	"availabilities_test",
	// 	"availabilities_test2",
	// 	"availability",
	// 	"availability_deltas",
	// 	"availability_deltas_test",
	// 	"availability_detailed",
	// }
	// for _, coll := range colls {
	// 	iter := fs.Collection(coll).Documents(context.Background())
	// 	for {
	// 		doc, err := iter.Next()
	// 		if err == iterator.Done {
	// 			// don't return an error, an empty object is what we want instead
	// 			break
	// 		}

	// 		_, err = doc.Ref.Delete(context.Background())
	// 		if err != nil {
	// 			fmt.Println(err)
	// 			return
	// 		}
	// 		fmt.Println("deleted ", doc.Ref.ID)
	// 	}
	// }

	// init influx
	// token, err := ioutil.ReadFile(".influx")
	// if err != nil {
	// 	panic(err)
	// }
	url := "https://us-central1-1.gcp.cloud2.influxdata.com"
	ifdb := influxdb2.NewClient(url, "AuGH2KZIYj84lb6T9-yGnqyVAQJdp8V7Gh3tv_Jd9zqtrQEGIMStseu0KgzrQv4HsLLAkkQqOZgR8qUFty2VmA==")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	// start webserver
	http.Handle("/", HandleAvailabilitySync(log, fs, ifdb))
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)

}
