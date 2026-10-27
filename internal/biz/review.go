package biz

import (
	"context"
	"errors"
	"fmt"
	"review-service/common"
	"review-service/gerr"
	"review-service/internal/data/model"
	"review-service/pkg/snowflake"
)

type ReviewRepo interface {
	SaveReview(context.Context, *model.ReviewInfo) (*model.ReviewInfo, error)
	GetReviewByOrderID(context.Context, int64) ([]*model.ReviewInfo, error)
}
type ReviewUsecase struct {
	repo ReviewRepo
	// log  *log.Helper
}

func NewReviewUsecase(repo ReviewRepo) *ReviewUsecase {
	return &ReviewUsecase{
		repo: repo,
	}
}
func (r *ReviewUsecase) CreateReview(ctx context.Context, review *model.ReviewInfo) (*model.ReviewInfo, error) {
	common.DebugLog(fmt.Sprintf("[biz] CreateReview,req:%+v", review))
	// 数据校验
	// 1.1 参数基础校验
	// 1.2 参数业务校验
	reviews, err := r.repo.GetReviewByOrderID(ctx, review.OrderID)
	if err != nil {
		return nil, errors.New("查询数据库失败")
	}
	if len(reviews) > 0 {
		return nil, gerr.OrderReviewedErr(fmt.Sprintf("订单 %d 已评价", review.OrderID))
	}
	// 生成review ID
	review.ReviewID = snowflake.GenID()
	// 查询订单和商品快照信息

	// 入库
	_, err = r.repo.SaveReview(ctx, review)
	if err != nil {
		return nil, err
	}
	return review, nil
}
