package grpc

import (
	"context"
	"google.golang.org/grpc"
	"log"
)

func CreateGRPCServer() *grpc.Server {
	loggingInterceptor := func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		log.Printf("gRPC method called: %s", info.FullMethod)
		return handler(ctx, req)
	}

	return grpc.NewServer(grpc.UnaryInterceptor(loggingInterceptor))
}
