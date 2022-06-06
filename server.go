package recreation

import (
	"context"

	"go.uber.org/zap"
)

type Server struct {
	client HTTPClient
	log    *zap.Logger

	ctx context.Context
}

func InitServer(ctx context.Context, log *zap.Logger) *Server {

	s := &Server{
		client: initObfuscator(ctx),
		log:    log,
		ctx:    ctx,
	}
	return s
}
