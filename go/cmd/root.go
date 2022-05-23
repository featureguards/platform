package cmd

import (
	"context"
	"net/url"
	"os"
	"runtime/debug"
	"strconv"

	"platform/go/core/app"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	EnvPhysicalBits = "APP_PHYSICAL_BITS"
	EnvDSN          = "APP_DSN"

	SmtpDSN = "SMTP_DSN"
)

func Recovery(ctx context.Context, p interface{}) (err error) {
	log.Errorf("panic triggered: %v\n%v", p, string(debug.Stack()))
	return status.Error(codes.Internal, "Unexpected")
}

func Init() (*app.App, error) {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
	formatter := log.TextFormatter{}
	formatter.DisableQuote = true
	log.SetFormatter(&formatter)

	physicalBitsStr := os.Getenv(EnvPhysicalBits)
	if physicalBitsStr == "" {
		physicalBitsStr = "2"
	}

	physicalBits, err := strconv.Atoi(physicalBitsStr)
	if err != nil {
		return nil, err
	}
	dsn := os.Getenv(EnvDSN)
	smtp, err := url.Parse(os.Getenv(SmtpDSN))
	if err != nil {
		log.Fatal(err)
	}
	config := app.Config{
		PhysicalBits: physicalBits,
		DSN:          dsn,
		SmtpURL:      smtp,
	}
	return app.Initialize(config)
}
