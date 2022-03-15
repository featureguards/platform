package cmd

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Recovery(ctx context.Context, p interface{}) (err error) {
	log.Errorf("panic triggered: %v", p)
	return status.Error(codes.Internal, "Unexpected")
}

func Init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}
