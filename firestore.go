package main

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func InitFirestore(ctx context.Context) (*firestore.Client, error) {

	var app *firebase.App

	// check if creds exist
	_, err := os.Stat("creds.json")
	if err != nil {
		// use default creds if not (means we're in the cloud)
		app, err = firebase.NewApp(context.Background(), nil)
	} else {
		// use file if creds.json does exist
		opt := option.WithCredentialsFile("creds.json")
		conf := &firebase.Config{ProjectID: "campr-app"}
		app, err = firebase.NewApp(context.Background(), conf, opt)
	}
	if err != nil {
		return nil, err
	}

	return app.Firestore(ctx)

}
