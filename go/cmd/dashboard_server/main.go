// Package main implements a server for Greeter service.
package main

import (
	"flag"
	"fmt"
	"net"

	"stackv2/go/cmd"
	"stackv2/go/grpc/dashboard"
	"stackv2/go/grpc/middleware/web_auth"
	"stackv2/go/grpc/middleware/web_log"
	pb_dashboard "stackv2/go/proto/dashboard"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

func main() {
	flag.Parse()
	app, err := cmd.Init()
	if err != nil {
		log.Fatal(err)
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	auth, err := web_auth.New(web_auth.AuthOpts{
		AllowedUnverifiedEmailMethods: []string{"/dashboard.Dashboard/GetUser"},
	})
	if err != nil {
		log.Fatal(err)
	}
	logger, err := web_log.New(web_log.Opts{})
	if err != nil {
		log.Fatal(err)
	}

	recovery := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandlerContext(cmd.Recovery),
	}

	server := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_recovery.StreamServerInterceptor(recovery...), logger.StreamServerInterceptor(), auth.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(recovery...), logger.UnaryServerInterceptor(), auth.UnaryServerInterceptor(),
		)),
	)

	dashboardServer, err := dashboard.New(dashboard.DashboardOpts{App: app})
	if err != nil {
		log.Fatal(err)
	}

	pb_dashboard.RegisterDashboardServer(server, dashboardServer)

	log.Printf("server listening at %v", lis.Addr())
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
