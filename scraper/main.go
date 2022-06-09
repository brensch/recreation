package main

import (
	"net/http"
	"os"

	"github.com/brensch/recreation"
	"go.uber.org/zap"
)

var (
	log *zap.Logger
)

func main() {
	logConfig := zap.NewProductionConfig()
	logConfig.Level.SetLevel(zap.DebugLevel)
	logConfig.Encoding = "console"

	// init logger
	var err error
	log, err = logConfig.Build()
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	log.Info("starting server")

	http.HandleFunc("/", HandleLogIP)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Info("using default port", zap.String("port", port))
	}

	// Start HTTP server.
	log.Info("listening", zap.String("port", port))
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("server had error", zap.Error(err))
	}
}

func HandleLogIP(w http.ResponseWriter, r *http.Request) {
	recreation.PingAll()

}
