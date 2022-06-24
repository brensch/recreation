package main

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"go.uber.org/zap"
)

type NotificationTarget struct {
	Username          string
	NotificationToken string
	NotificationType  NotificationType
}

type CampsiteMonitorSiteRequest struct {
	Notifier NotificationTarget
	SiteID   string
	Date     time.Time
}

type CampsiteMonitorGroundRequest struct {
	Notifier NotificationTarget
	GroundID string
	Date     time.Time
}

type MonitorRequest struct {
	Notifier NotificationTarget
	GroundID string
	// if siteids is nil, do any
	SiteIDs []string
	Dates   []time.Time
}

const (
	CollectionMonitorSiteRequests   = "monitor_site_requests"
	CollectionMonitorGroundRequests = "monitor_ground_requests"
	CollectionMonitorRequests       = "monitor_requests"
)

func CreateCampsiteMonitorRequest(ctx context.Context, log *zap.Logger, fs *firestore.Client, req CampsiteMonitorSiteRequest) (string, error) {

	ref, _, err := fs.Collection(CollectionMonitorSiteRequests).Add(ctx, req)
	if err != nil {
		log.Error("failed to record campsite monitor request")
		return "", err
	}

	log.Info("recorded campsite monitor request", zap.String("firestore_doc", ref.ID))

	return ref.ID, nil
}

func GetMonitorRequests(ctx context.Context, log *zap.Logger, fs *firestore.Client, groundID string, date time.Time) ([]CampsiteMonitorSiteRequest, error) {

	docs, err := fs.Collection(CollectionMonitorRequests).
		Where("GroundID", "==", groundID).
		Where("Dates", "array-contains", date).
		Documents(ctx).GetAll()
	if err != nil {
		log.Error("failed to get campsite monitor request")
		return nil, err
	}

	var reqs []CampsiteMonitorSiteRequest
	for _, doc := range docs {
		var req CampsiteMonitorSiteRequest
		err = doc.DataTo(&req)
		if err != nil {
			return nil, err
		}
		reqs = append(reqs, req)
	}

	log.Info("got all monitor requests for delta", zap.Int("count", len(reqs)))

	return reqs, nil
}

func GetCampsiteMonitorRequests(ctx context.Context, log *zap.Logger, fs *firestore.Client, campsiteID string, date time.Time) ([]CampsiteMonitorSiteRequest, error) {

	docs, err := fs.Collection(CollectionMonitorSiteRequests).
		Where("SiteID", "==", campsiteID).
		Where("Date", "==", date).
		Documents(ctx).GetAll()
	if err != nil {
		log.Error("failed to get campsite monitor request")
		return nil, err
	}

	var reqs []CampsiteMonitorSiteRequest
	for _, doc := range docs {
		var req CampsiteMonitorSiteRequest
		err = doc.DataTo(&req)
		if err != nil {
			return nil, err
		}
		reqs = append(reqs, req)
	}

	log.Info("got all site requests for delta", zap.Int("count", len(reqs)))

	return reqs, nil
}

func GetAllCampsiteMonitorRequests(ctx context.Context, log *zap.Logger, fs *firestore.Client, campsiteID string) ([]CampsiteMonitorSiteRequest, error) {

	docs, err := fs.Collection(CollectionMonitorSiteRequests).
		Where("SiteID", "==", campsiteID).
		Documents(ctx).GetAll()
	if err != nil {
		log.Error("failed to get campsite monitor request")
		return nil, err
	}

	var reqs []CampsiteMonitorSiteRequest
	for _, doc := range docs {
		var req CampsiteMonitorSiteRequest
		err = doc.DataTo(&req)
		if err != nil {
			return nil, err
		}
		reqs = append(reqs, req)
	}

	log.Info("got all site requests for delta", zap.Int("count", len(reqs)))

	return reqs, nil
}

func GetCampgroundMonitorRequests(ctx context.Context, log *zap.Logger, fs *firestore.Client, groundID string, date time.Time) ([]CampsiteMonitorGroundRequest, error) {

	docs, err := fs.Collection(CollectionMonitorGroundRequests).
		Where("GroundID", "==", groundID).
		Where("Date", "==", date).
		Documents(ctx).GetAll()
	if err != nil {
		log.Error("failed to get campsite monitor request", zap.Error(err))
		return nil, err
	}

	var reqs []CampsiteMonitorGroundRequest
	for _, doc := range docs {
		var req CampsiteMonitorGroundRequest
		err = doc.DataTo(&req)
		if err != nil {
			return nil, err
		}
		reqs = append(reqs, req)
	}

	log.Info("got all ground requests for delta", zap.Int("count", len(reqs)))

	return reqs, nil
}

func GetAllCampgroundMonitorRequests(ctx context.Context, log *zap.Logger, fs *firestore.Client, groundID string) ([]CampsiteMonitorGroundRequest, error) {

	docs, err := fs.Collection(CollectionMonitorGroundRequests).
		Where("GroundID", "==", groundID).
		Documents(ctx).GetAll()
	if err != nil {
		log.Error("failed to get campsite monitor request", zap.Error(err))
		return nil, err
	}

	var reqs []CampsiteMonitorGroundRequest
	for _, doc := range docs {
		var req CampsiteMonitorGroundRequest
		err = doc.DataTo(&req)
		if err != nil {
			return nil, err
		}
		reqs = append(reqs, req)
	}

	log.Info("got all ground requests for delta", zap.Int("count", len(reqs)))

	return reqs, nil
}

// func MonitorWholeSite(ctx context.Context, log *zap.Logger, fs *firestore.Client, groundID string) error {

// 	// get avail to get sites
// 	// TODO:figure out endpoint that returns sites
// 	avail, err := api.GetAvailability(ctx, log, "https://us-west1-campr-app.cloudfunctions.net/HandleProxyRequest", groundID, time.Now())
// 	if err != nil {
// 		return err
// 	}

// 	batch := fs.Batch()
// 	_ = batch
// 	docCount := 0

// 	for i := 0; i < 1; i++ {
// 		for _, site := range avail.Campsites {

// 			monitorRequest := CampsiteMonitorRequest{
// 				SiteID:            site.CampsiteID,
// 				Date:              time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()+i, 0, 0, 0, 0, time.UTC),
// 				Username:          "brend",
// 				NotificationToken: "o.gMQajO2iSWtnlJpUBdXDl55CRO9UNhLw",
// 			}
// 			_ = monitorRequest

// 			docCount++
// 			// ref := fs.Collection(CollectionMonitorRequests).NewDoc()
// 			// batch.Create(ref, monitorRequest)
// 			// if docCount == 500 {
// 			// 	log.Info("flushed monitor requests")
// 			// 	batch.Commit(ctx)
// 			// 	docCount = 0
// 			// 	batch = fs.Batch()
// 			// }

// 		}
// 	}

// 	fmt.Println("will use docs", docCount)

// 	return nil
// }

func MonitorWholeGround(ctx context.Context, log *zap.Logger, fs *firestore.Client, groundID string, date time.Time) error {

	for i := 0; i < 90; i++ {

		monitorRequest := CampsiteMonitorGroundRequest{
			Notifier: NotificationTarget{
				Username:          "brend",
				NotificationToken: "o.gMQajO2iSWtnlJpUBdXDl55CRO9UNhLw",
			},
			GroundID: groundID,
			Date:     time.Date(date.Year(), date.Month(), date.Day()+i, 0, 0, 0, 0, time.UTC),
		}

		_, _, err := fs.Collection(CollectionMonitorGroundRequests).Add(ctx, monitorRequest)
		if err != nil {
			return err
		}
	}

	log.Info("added whole ground monitor", zap.String("ground_id", groundID))

	return nil
}

func MonitorGround(ctx context.Context, log *zap.Logger, fs *firestore.Client, groundID string, date time.Time) error {

	monitorRequest := MonitorRequest{
		Notifier: NotificationTarget{
			Username:          "brend",
			NotificationToken: "o.gMQajO2iSWtnlJpUBdXDl55CRO9UNhLw",
		},
		GroundID: groundID,
	}

	for i := 0; i < 90; i++ {

		monitorRequest.Dates = append(monitorRequest.Dates, time.Date(date.Year(), date.Month(), date.Day()+i, 0, 0, 0, 0, time.UTC))

	}

	_, _, err := fs.Collection(CollectionMonitorRequests).Add(ctx, monitorRequest)
	if err != nil {
		return err
	}

	log.Info("added whole ground monitor", zap.String("ground_id", groundID))

	return nil
}

func MonitorGroundOnDates(ctx context.Context, log *zap.Logger, fs *firestore.Client, groundID string, dates []time.Time) error {
	monitorRequest := MonitorRequest{
		Notifier: NotificationTarget{
			Username:          "brend",
			NotificationToken: "o.gMQajO2iSWtnlJpUBdXDl55CRO9UNhLw",
		},
		GroundID: groundID,
		Dates:    dates,
	}

	_, _, err := fs.Collection(CollectionMonitorRequests).Add(ctx, monitorRequest)
	if err != nil {
		return err
	}

	log.Info("added whole ground monitor", zap.String("ground_id", groundID))

	return nil
}

// 259084
