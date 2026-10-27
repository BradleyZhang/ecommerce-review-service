package gerr

import "google.golang.org/grpc/codes"

//go:generate go run ./tool/gen_func.go

type Reason string

var (
	C = codes.InvalidArgument // 客户端出问题
	S = codes.Internal        // 服务端出问题

	Unauth = codes.Unauthenticated  // 认证失败
	PD     = codes.PermissionDenied // 授权失败
	NF     = codes.NotFound         // 资源找不到
)

type ReasonCode struct {
	Reason Reason
	Code   codes.Code
}

var (
	OrderReviewed = ReasonCode{"ORDER_REVIEWED", C}
	DBFailed      = ReasonCode{"DB_FAILED", S}
)
