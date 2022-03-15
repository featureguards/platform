// Package main implements a server for Greeter service.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"stackv2/go/cmd"
	pb_global "stackv2/go/proto/global"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement global.GlobalServer.
type globalServer struct {
	pb_global.UnimplementedGlobalServer
}

// SayHello implements global.GlobalServer
func (s *globalServer) SayHello(ctx context.Context, in *pb_global.GreetRequest) (*pb_global.GreetReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb_global.GreetReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	flag.Parse()
	cmd.Init()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	recovery := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandlerContext(cmd.Recovery),
	}

	// No auth
	global := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(grpc_recovery.StreamServerInterceptor(recovery...))),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(grpc_recovery.UnaryServerInterceptor(recovery...))),
	)

	pb_global.RegisterGlobalServer(global, &globalServer{})

	log.Printf("server listening at %v", lis.Addr())
	if err := global.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
