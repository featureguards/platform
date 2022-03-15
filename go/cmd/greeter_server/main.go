// Package main implements a server for Greeter service.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"stackv2/go/cmd"
	"stackv2/go/grpc/middleware/web_auth"
	pb_greeter "stackv2/go/proto/greeter"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement greeter.GreeterServer.
type greeterServer struct {
	pb_greeter.UnimplementedGreeterServer
}

// SayHello implements greeter.GreeterServer
func (s *greeterServer) SayHello(ctx context.Context, in *pb_greeter.HelloRequest) (*pb_greeter.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb_greeter.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	flag.Parse()
	cmd.Init()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	auth, err := web_auth.New(web_auth.AuthOpts{})
	if err != nil {
		log.Fatal(err)
	}

	recovery := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandlerContext(cmd.Recovery),
	}

	greeter := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_auth.StreamServerInterceptor(auth.Authenticate), grpc_recovery.StreamServerInterceptor(recovery...),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_auth.UnaryServerInterceptor(auth.Authenticate), grpc_recovery.UnaryServerInterceptor(recovery...))),
	)

	pb_greeter.RegisterGreeterServer(greeter, &greeterServer{})

	log.Printf("server listening at %v", lis.Addr())
	if err := greeter.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
