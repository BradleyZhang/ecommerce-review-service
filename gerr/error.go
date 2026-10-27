package gerr

import (
	"review-service/common"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewGRPCError(code codes.Code, msg string, reason Reason) error {
	st := status.New(code, msg)
	detail := &errdetails.ErrorInfo{
		Reason: string(reason),
		Domain: common.Domain,
	}
	st, err := st.WithDetails(detail)
	if err != nil {
		common.DebugLog("[gRPC Err]fail to create details code:", code, " msg:", msg, " reason:", reason, "err:", err)
		return status.Error(codes.Internal, "fail to create error details")
	}
	return st.Err()
}
