package server

import (
	pb "review-service/api/gen/go/review/v1"
	"review-service/common"

	"buf.build/go/protovalidate"
	protovalidate_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"google.golang.org/grpc"
)

func NewGRPCServer(srv pb.ReviewServiceServer) *grpc.Server {
	var opts []grpc.ServerOption
	validator, err := protovalidate.New()
	if err != nil {
		common.FatalLog("validator 初始化失败", err)
	}
	interceptor := protovalidate_middleware.UnaryServerInterceptor(validator)
	opts = append(opts, grpc.UnaryInterceptor(interceptor))
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterReviewServiceServer(grpcServer, srv)
	return grpcServer
}
