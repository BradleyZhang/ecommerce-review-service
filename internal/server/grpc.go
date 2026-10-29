package server

import (
	"fmt"

	pb "github.com/BradleyZhang/ecommerce-review-service/api/gen/go/review/v1"

	"buf.build/go/protovalidate"
	protovalidate_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func NewGRPCServer(srv pb.ReviewServiceServer) (*grpc.Server, error) {
	var opts []grpc.ServerOption
	validator, err := protovalidate.New()
	if err != nil {
		return nil, fmt.Errorf("validator 初始化失败: %w", err)
	}
	interceptor := protovalidate_middleware.UnaryServerInterceptor(validator)
	opts = append(opts, grpc.UnaryInterceptor(interceptor))
	grpcServer := grpc.NewServer(opts...)

	healthcheck := health.NewServer() // healthcheck 服务器 用于 consul
	healthpb.RegisterHealthServer(grpcServer, healthcheck)
	// healthcheck.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	pb.RegisterReviewServiceServer(grpcServer, srv)
	return grpcServer, nil
}
