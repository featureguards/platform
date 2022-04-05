package web_log

import (
	"context"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Opts struct {
}

type Logger struct {
	logger *log.Logger
}

func New(opts Opts) (*Logger, error) {
	logger := log.New()
	logger.SetReportCaller(false)
	return &Logger{logger: logger}, nil
}

func (l *Logger) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		l.logger.Info("-> " + info.FullMethod)
		defer l.logger.Info("<- " + info.FullMethod)
		return handler(ctx, req)
	}
}

func (l *Logger) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		l.logger.Info("-> " + info.FullMethod)
		defer l.logger.Info("<- " + info.FullMethod)
		return handler(srv, stream)
	}
}
