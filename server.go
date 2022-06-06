package recreation

import (
	"context"
	"time"

	"go.uber.org/zap"
)

type Server struct {
	client HTTPClient
	log    *zap.Logger

	ctx context.Context
}

func InitServer(ctx context.Context, log *zap.Logger, apiPause time.Duration) *Server {

	s := &Server{
		client: initObfuscator(ctx, apiPause),
		log:    log,
		ctx:    ctx,
	}
	return s
}
