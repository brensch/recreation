package recreation

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"go.uber.org/zap"
)

func SubmitData() error {

	token, err := ioutil.ReadFile(".influx")
	if err != nil {
		return err
	}

	url := "https://us-central1-1.gcp.cloud2.influxdata.com"

	client := influxdb2.NewClient(url, string(token))

	defer client.Close()

	org := "brensch@tuta.io"
	bucket := "test"
	writeAPI := client.WriteAPIBlocking(org, bucket)
	for value := 0; value < 5; value++ {
		tags := map[string]string{
			"tagname1": "tagvalue1",
		}
		fields := map[string]interface{}{
			"field1": value,
		}
		point := write.NewPoint("measurement1", tags, fields, time.Now())
		time.Sleep(1 * time.Second) // separate points by 1 second

		if err := writeAPI.WritePoint(context.Background(), point); err != nil {
			log.Fatal("failed to write", zap.Error(err))
		}
	}

	queryAPI := client.QueryAPI(org)
	// query := `from(bucket: "test")
	//         |> range(start: -10m)
	//         |> filter(fn: (r) => r._measurement == "measurement1")`

	query := `from(bucket: "test")
              |> range(start: -10m)
              |> filter(fn: (r) => r._measurement == "measurement1")
              |> mean()`
	results, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		log.Fatal("failed to query", zap.Error(err))
	}
	for results.Next() {
		fmt.Println(results.Record())
	}
	if err := results.Err(); err != nil {
		log.Fatal("failed to get results", zap.Error(err))
	}

	return nil
}
