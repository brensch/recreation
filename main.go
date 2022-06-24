package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {

	// init logger
	logConfig := zap.NewProductionConfig()
	logConfig.Level.SetLevel(zap.DebugLevel)
	// this ensures google logs pick things up properly
	logConfig.EncoderConfig.MessageKey = "message"
	logConfig.EncoderConfig.LevelKey = "severity"
	logConfig.EncoderConfig.TimeKey = "time"
	logConfig.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	logConfig.EncoderConfig.EncodeDuration = zapcore.MillisDurationEncoder
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
	// 	// "availabilities_test",
	// 	// "availabilities_test2",
	// 	// "availability",
	// 	"availability_deltas2",
	// 	// "availability_deltas_test",
	// 	// "availability_detailed2",
	// 	"availability_detailed_test2",
	// 	"availability_detailed",
	// 	// "monitor_ground_requests",
	// 	// "monitor_site_requests",
	// }
	// batch := fs.Batch()
	// docCount := 0

	// for _, coll := range colls {
	// 	iter := fs.Collection(coll).Documents(context.Background())
	// 	for {
	// 		doc, err := iter.Next()
	// 		if err == iterator.Done {
	// 			// don't return an error, an empty object is what we want instead
	// 			break
	// 		}

	// 		batch.Delete(doc.Ref)
	// 		docCount++
	// 		if docCount == 500 {
	// 			log.Info("flushed delete batch")
	// 			batch.Commit(context.Background())
	// 			docCount = 0
	// 			batch = fs.Batch()
	// 		}

	// 		fmt.Println("deleted ", doc.Ref.ID)
	// 	}
	// }
	// batch.Commit(context.Background())

	// err = MonitorGround(context.Background(), log, fs, "10005253", time.Now().Add(-24*time.Hour))
	// if err != nil {
	// 	panic(err)
	// }
	// grounds := []string{
	// 	"232464", // kalaloch
	// 	"259084", // fairholme
	// 	"247592", // hoh rainforest
	// 	"251365", // falls creek
	// }

	// for _, ground := range grounds {
	// 	err = MonitorGroundOnDates(context.Background(), log, fs, ground, []time.Time{
	// 		time.Date(2022, 7, 25, 0, 0, 0, 0, time.UTC),
	// 		time.Date(2022, 7, 26, 0, 0, 0, 0, time.UTC),
	// 		time.Date(2022, 7, 27, 0, 0, 0, 0, time.UTC),
	// 		time.Date(2022, 7, 28, 0, 0, 0, 0, time.UTC),
	// 	})
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }

	// init influx
	// token, err := ioutil.ReadFile(".influx")
	// if err != nil {
	// 	panic(err)
	// }
	url := "https://us-central1-1.gcp.cloud2.influxdata.com"
	// url := "http://localhost:8086"
	ifdb := influxdb2.NewClient(url, "AuGH2KZIYj84lb6T9-yGnqyVAQJdp8V7Gh3tv_Jd9zqtrQEGIMStseu0KgzrQv4HsLLAkkQqOZgR8qUFty2VmA==")
	// ifdb := influxdb2.NewClient(url, "my-super-secret-auth-token")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	// start webserver
	http.Handle("/", HandleAvailabilitySync(log, fs, ifdb))
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)

}
