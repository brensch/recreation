package main

import (
	"context"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestCampsiteMonitorRequest(t *testing.T) {

	ctx := context.Background()
	fs, err := InitFirestore(ctx)
	if err != nil {
		t.FailNow()
	}

	log, _ := zap.NewDevelopment()

	sentReq := CampsiteMonitorRequest{
		SiteID:            "6969",
		GroundID:          "9696",
		NotificationToken: "toke",
		Username:          "brend",
		Date:              time.Date(2022, 6, 28, 0, 0, 0, 0, time.UTC),
	}
	ref, err := CreateCampsiteMonitorRequest(ctx, log, fs, sentReq)
	if err != nil {
		t.Error(err)
	}

	foundReqs, err := GetCampsiteMonitorRequests(ctx, log, fs, sentReq.SiteID, sentReq.Date)
	if err != nil {
		t.Error(err)
	}

	t.Log(foundReqs)

	if len(foundReqs) != 1 {
		t.Log("didn't get exactly one doc returned")
		t.FailNow()
	}

	if foundReqs[0].SiteID != sentReq.SiteID {
		t.Errorf("got campsite id: %s, expected %s", foundReqs[0].SiteID, sentReq.SiteID)
	}

	if !foundReqs[0].Date.Equal(sentReq.Date) {
		t.Errorf("got date: %s, expected %s", foundReqs[0].Date, sentReq.Date)
	}

	_, err = fs.Collection(CollectionMonitorRequests).Doc(ref).Delete(ctx)
	if err != nil {
		t.Error(err)
	}

}
