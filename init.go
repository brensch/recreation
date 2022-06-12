package recreation

import (
	"net/http"
	"net/http/httputil"
	"strings"

	"go.uber.org/zap"
)

// using globals since i'm using cloud functions and can't pass things like i would regular handler funcs
var (
	// o     api.Obfuscator
	// p     ProxyClient
	log   *zap.Logger
	proxy httputil.ReverseProxy
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

	proxy = httputil.ReverseProxy{
		// Transport: roundTripper(rt),
		Director: func(req *http.Request) {
			target := "www.recreation.gov"
			req.URL.Scheme = "https"
			req.URL.Host = target
			req.Host = target
			req.Header["X-Forwarded-For"] = nil
			req.Header.Set("User-Agent", "PostmanRuntime/7.29.0")
			req.URL.Path = strings.TrimPrefix(req.URL.Path, "/HandleProxyRequest")
		},
	}

	// init logger
	var err error
	log, err = logConfig.Build()
	if err != nil {
		// this indicates a bug or some way that zap can fail i'm not aware of
		panic(err)
	}
}
