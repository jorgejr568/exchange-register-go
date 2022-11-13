package server

import "context"

type Server interface {
	GracefulListenAndShutdown(ctx context.Context) error
}
