package recreation_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/brensch/recreation"
)

func TestDoSearchGeo(t *testing.T) {
	ctx := context.Background()
	client := recreation.InitObfuscator(ctx)

	res, err := recreation.DoSearchGeo(ctx, client, 37.3859, -122.0882)
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
