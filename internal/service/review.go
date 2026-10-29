package service

import (
	"context"

	pb "github.com/BradleyZhang/ecommerce-review-service/api/gen/go/review/v1"
	"github.com/BradleyZhang/ecommerce-review-service/common"
	"github.com/BradleyZhang/ecommerce-review-service/internal/biz"
	"github.com/BradleyZhang/ecommerce-review-service/internal/data/model"
)

type ReviewService struct {
	pb.UnimplementedReviewServiceServer
	uc *biz.ReviewUsecase
}

func NewReviewService(uc *biz.ReviewUsecase) *ReviewService {
	return &ReviewService{uc: uc}
}

// C端创建评价
func (r *ReviewService) CreateReview(ctx context.Context, request *pb.CreateReviewRequest) (*pb.CreateReviewResponse, error) {
	var anonymous int32
	if request.Anonymous {
		anonymous = 1
	}
	review, err := r.uc.CreateReview(ctx, &model.ReviewInfo{
		UserID:       request.UserId,
		OrderID:      request.OrderId,
		Score:        request.Score,
		ServiceScore: request.ServiceScore,
		ExpressScore: request.ExpressScore,
		Content:      request.Content,
		PicInfo:      request.PicInfo,
		VideoInfo:    request.VideoInfo,
		Anonymous:    anonymous,
		StoreID:      request.StoreId,
		Status:       0,
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateReviewResponse{ReviewId: review.ReviewID}, err
}

// // C端获取评价详情
// func (r *ReviewService) GetReview(context.Context, *pb.GetReviewRequest) (*pb.GetReviewResponse, error) {
// 	panic("unimpl")
// }

// // O端审核评价
// func (r *ReviewService) AuditReview(context.Context, *pb.AuditReviewRequest) (*pb.AuditReviewResponse, error) {
// 	panic("unimpl")
// }

// B端回复评价
func (r *ReviewService) ReplyReview(ctx context.Context, request *pb.ReplyReviewRequest) (*pb.ReplyReviewResponse, error) {
	common.DebugLogf("[service] ResponseReview req:%v", request)
	param := &biz.ReplyParam{
		ReviewID:  request.ReviewId,
		StoreID:   request.StoreId,
		Content:   request.Content,
		PicInfo:   request.PicInfo,
		VideoInfo: request.VideoInfo,
	}
	reply, err := r.uc.CreateReply(ctx, param)
	if err != nil {
		return nil, err
	}
	return &pb.ReplyReviewResponse{ReplyId: reply.ReplyID}, nil
}

// // B端申诉评价
// func (r *ReviewService) AppealReview(context.Context, *pb.AppealReviewRequest) (*pb.AppealReviewResponse, error) {
// 	panic("unimpl")
// }

// // O端评价申诉审核
// func (r *ReviewService) AuditAppeal(context.Context, *pb.AuditAppealRequest) (*pb.AuditAppealResponse, error) {
// 	panic("unimpl")
// }

// // C端查看userID下所有评价
// func (r *ReviewService) ListReviewByUserID(context.Context, *pb.ListReviewByUserIDRequest) (*pb.ListReviewByUserIDResponse, error) {
// 	panic("unimpl")
// }
