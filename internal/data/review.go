package data

import (
	"context"
	"log"
	"review-service/internal/biz"
	"review-service/internal/data/model"
)

type ReviewRepo struct {
	data   *Data
	logger *log.Logger
}

// GetReviewByOrderID implements [biz.ReviewRepo].
func (r *ReviewRepo) GetReviewByOrderID(ctx context.Context, orderID int64) ([]*model.ReviewInfo, error) {
	info := r.data.query.ReviewInfo
	ans, err := info.WithContext(ctx).Where(info.OrderID.Eq(orderID)).Find()
	if err != nil {
		return nil, err
	}
	return ans, nil
}

// SaveReview implements [biz.ReviewRepo].
func (r *ReviewRepo) SaveReview(ctx context.Context, review *model.ReviewInfo) (*model.ReviewInfo, error) {
	err := r.data.query.ReviewInfo.WithContext(ctx).Save(review)
	return review, err
}

func NewReviewRepo(data *Data, logger *log.Logger) biz.ReviewRepo {
	return &ReviewRepo{
		data:   data,
		logger: logger,
	}
}
