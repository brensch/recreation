package api

import (
	"context"
	"fmt"
	"testing"

	"go.uber.org/zap"
)

func TestDoSearchGeo(t *testing.T) {
	ctx := context.Background()
	log, _ := zap.NewDevelopment()
	// o := InitAgentRandomiser(ctx)

	res, err := SearchGeo(ctx, log, "http://recreation.gov", 37.3859, -122.0882)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	// should be around 262
	t.Log(len(res.Results))

	for _, result := range res.Results {
		fmt.Print(result.EntityID, ",")
	}
}
