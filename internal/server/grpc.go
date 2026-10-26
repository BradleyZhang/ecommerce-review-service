package server

import (
	v1 "review-service/api/review/v1"

	"google.golang.org/grpc"
)

func NewGRPCServer(srv v1.ReviewServer) *grpc.Server {
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	v1.RegisterReviewServer(grpcServer, srv)
	return grpcServer
}
