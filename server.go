package recreation

import (
	"context"
	"time"
)

type Server struct {
	client HTTPClient

	ctx context.Context
}

func InitServer(ctx context.Context, apiPause time.Duration) *Server {

	s := &Server{
		client: initObfuscator(ctx, apiPause),
		ctx:    ctx,
	}
	return s
}
