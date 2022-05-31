// Package main implements a server for Greeter service.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"platform/go/cmd"
	"platform/go/grpc/auth"
	"platform/go/grpc/server"

	log "github.com/sirupsen/logrus"
)

var (
	port = flag.Int("port", 50052, "The server port")
)

func main() {
	flag.Parse()
	app, err := cmd.Init()
	if err != nil {
		log.Fatal(fmt.Sprintf("%+v", err))
	}

	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		sig := <-sigs
		log.Infof("Ratelimit server received %v, shutting down gracefully", sig)
		cancel()
	}()

	_, srv, lis, err := auth.Listen(ctx, app, server.WithPort(*port))
	if err != nil {
		log.Fatalf("failed to listen: %s", err)
	}

	log.Printf("server listening at %s", lis.Addr())

	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
