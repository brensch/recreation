package recreation

import (
	"context"
	"net/http/httputil"

	"cloud.google.com/go/firestore"
	"go.uber.org/zap"
)

// using globals since i'm using cloud functions and can't pass things like i would regular handler funcs
var (
	// o     api.Obfuscator
	// p     ProxyClient
	log   *zap.Logger
	proxy httputil.ReverseProxy
	fs    *firestore.Client
)

func init() {
	// o = api.InitAgentRandomiser(context.Background())
	logConfig := zap.NewProductionConfig()
	logConfig.Level.SetLevel(zap.DebugLevel)
	// this ensures google logs pick things up properly
	logConfig.EncoderConfig.MessageKey = "message"
	logConfig.EncoderConfig.LevelKey = "severity"
	logConfig.EncoderConfig.TimeKey = "time"
	// logConfig.Encoding = "console"

	proxy = MakeProxy()

	var err error
	fs, err = InitFirestore(context.Background())
	if err != nil {
		panic(err)
	}

	// init logger
	log, err = logConfig.Build()
	if err != nil {
		// this indicates a bug or some way that zap can fail i'm not aware of
		panic(err)
	}
}
