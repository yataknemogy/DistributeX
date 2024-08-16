package main

import (
	"DistributeX/DistributeX/api/executionnodepb"
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
	executionnodepb.UnimplementedExecutionNodeServer
}

func loggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	log.Printf("gRPC method called: %s", info.FullMethod)
	resp, err := handler(ctx, req)
	if err != nil {
		log.Printf("Error in gRPC method %s: %v", info.FullMethod, err)
	}
	return resp, err
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	serverOpts := []grpc.ServerOption{
		grpc.UnaryInterceptor(loggingInterceptor),
	}

	s := grpc.NewServer(serverOpts...)
	executionnodepb.RegisterExecutionNodeServer(s, &server{})

	log.Println("ExecutionNode is running on port 50051")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
