package service

import (
	"context"
	pb "review-service/api/review/v1"
	"review-service/internal/biz"
	"review-service/internal/data/model"
)

type ReviewService struct {
	pb.UnimplementedReviewServiceServer
	uc *biz.ReviewUsecase
}

func NewReviewService(uc *biz.ReviewUsecase) *ReviewService {
	return &ReviewService{uc: uc}
}
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
		Status:       0,
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateReviewResponse{ReviewId: review.ReviewID}, err
}
