package main

import (
	"bytes"
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/xconstruct/go-pushbullet"
	"go.uber.org/zap"
)

type NotificationType string

const (
	NotificationTypePushBullet NotificationType = "pushbullet"
)

type Notification struct {
	delta    CampsiteDelta2
	notifier NotificationTarget
}

func NotifyOfDeltas(ctx context.Context, log *zap.Logger, fs *firestore.Client, deltas []CampsiteDelta2) error {

	if len(deltas) == 0 {
		return fmt.Errorf("not no deltas pass to notify")
	}

	var notifications []Notification
	for _, delta := range deltas {
		log := log.With(
			zap.String("site_id", delta.SiteID),
			zap.Time("delta_date", delta.DateAffected),
		)

		groundMonitors, err := GetMonitorRequests(ctx, log, fs, delta.GroundID, delta.DateAffected)
		if err != nil {
			continue
		}

		for _, monitor := range groundMonitors {
			notifications = append(notifications, Notification{
				delta:    delta,
				notifier: monitor.Notifier,
			})
		}

	}

	// // var
	// // do this to stop having to call a million empty fs reads
	// if len(deltas) > 50 {
	// 	allGroundMonitors, err := GetAllCampgroundMonitorRequests(ctx, log, fs, deltas[0].GroundID)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	allSiteMonitors, err := GetAllCampsiteMonitorRequests(ctx, log, fs, deltas[0].SiteID)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	for _, delta := range deltas {

	// 		for _, monitor := range allGroundMonitors {
	// 			notifications = append(notifications, Notification{
	// 				delta:    delta,
	// 				notifier: monitor.Notifier,
	// 			})
	// 		}
	// 		for _, monitor := range allSiteMonitors {
	// 			notifications = append(notifications, Notification{
	// 				delta:    delta,
	// 				notifier: monitor.Notifier,
	// 			})
	// 		}
	// 	}
	// } else {

	// 	// all deltas will belong to the one ground
	// 	for _, delta := range deltas {
	// 		log := log.With(
	// 			zap.String("site_id", delta.SiteID),
	// 			zap.Time("delta_date", delta.DateAffected),
	// 		)

	// 		groundMonitors, err := GetMonitorRequests(ctx, log, fs, deltas[0].GroundID, delta.DateAffected)
	// 		if err != nil {
	// 			continue
	// 		}

	// 		for _, monitor := range groundMonitors {
	// 			notifications = append(notifications, Notification{
	// 				delta:    delta,
	// 				notifier: monitor.Notifier,
	// 			})
	// 		}

	// 		siteMonitors, err := GetCampgroundMonitorRequests(ctx, log, fs, delta.SiteID, delta.DateAffected)
	// 		if err != nil {
	// 			continue
	// 		}

	// 		for _, monitor := range siteMonitors {
	// 			notifications = append(notifications, Notification{
	// 				delta:    delta,
	// 				notifier: monitor.Notifier,
	// 			})
	// 		}

	// 	}
	// }

	// add all notifications to a map and then ship them all at once
	userNotifications := make(map[NotificationTarget][]CampsiteDelta2)

	for _, notification := range notifications {

		userNotifications[notification.notifier] = append(userNotifications[notification.notifier], notification.delta)
		// TODO: add switch for different methods here. ie this but good.

	}

	for target, deltas := range userNotifications {

		pb := pushbullet.New(target.NotificationToken)
		devs, err := pb.Devices()
		if err != nil {
			log.Error("couldn't get user devices", zap.Error(err))
			continue
		}

		if len(devs) == 0 {
			log.Error("user has no devices")
			continue
		}

		var buf bytes.Buffer

		if len(deltas) < 100 {
			for _, delta := range deltas {
				buf.WriteString(fmt.Sprintf("Campsite %s changed from %s to %s for %s: https://www.recreation.gov/camping/campsites/%s\n",
					delta.SiteID,
					StateStrings[delta.OldState],
					StateStrings[delta.NewState],
					delta.DateAffected.UTC().Format("2006/01/02"),
					delta.SiteID,
				))

			}
		} else {
			buf.WriteString(fmt.Sprintf("found %d changes.", len(deltas)))
		}

		log.Info("sending notification to user", zap.String("user", target.Username))

		err = pb.PushNote(
			devs[0].Iden,
			fmt.Sprintf("Campsite %s availability changed. ", deltas[0].GroundID),
			buf.String(),
		)
		if err != nil {
			log.Error("send notification to user", zap.Error(err))
			continue
		}
	}

	return nil

}

// func NotifierUUID(notifier NotificationTarget) string {
// 	return fmt.Sprintf("%s-%s", notifier.Username, notifier.NotificationType)
// }

// o.gMQajO2iSWtnlJpUBdXDl55CRO9UNhLw
