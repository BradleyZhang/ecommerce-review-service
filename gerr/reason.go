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

var ( // 通用
	InvalidArgument = ReasonCode{"ERR_INVALIDE_ARGUMENT", C}
	DBFailed        = ReasonCode{"ERR_DB_FAILED", S}
	AccessDenied    = ReasonCode{"ACCESS_DENIED", PD}
)

var ( // 业务
	OrderReviewed  = ReasonCode{"ERR_ORDER_REVIEWED", C}
	DuplicateReply = ReasonCode{"DUPLICATE_REPLY", C}
)
