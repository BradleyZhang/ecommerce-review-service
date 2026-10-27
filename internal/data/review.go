package data

import (
	"context"
	"log"
	"review-service/common"
	"review-service/gerr"
	"review-service/internal/biz"
	"review-service/internal/data/model"
	"review-service/internal/data/query"
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
		return nil, gerr.DBFailedErr(err.Error())
	}
	return ans, nil
}

// SaveReview implements [biz.ReviewRepo].
func (r *ReviewRepo) SaveReview(ctx context.Context, review *model.ReviewInfo) (*model.ReviewInfo, error) {
	err := r.data.query.ReviewInfo.WithContext(ctx).Save(review)
	if err != nil {
		return nil, gerr.DBFailedErr(err.Error())
	}
	return review, nil
}
func (r *ReviewRepo) GetReview(ctx context.Context, id int64) (*model.ReviewInfo, error) {
	info := r.data.query.ReviewInfo

	review, err := info.WithContext(ctx).Where(info.ReviewID.Eq(id)).First()
	if err != nil {
		return nil, gerr.DBFailedErr(err.Error())
	}
	return review, nil
}

// SaveReply implements [biz.ReviewRepo].
func (r *ReviewRepo) SaveReply(ctx context.Context, reply *model.ReviewReplyInfo) (*model.ReviewReplyInfo, error) {
	err := r.data.query.Transaction(func(tx *query.Query) error {
		if err := tx.ReviewReplyInfo.WithContext(ctx).
			Save(reply); err != nil {
			common.ErrorLogf("[data] SaveReply create reply fail,err:%v", err)
			return err
		}
		// 更新review的hasReply
		if _, err := tx.ReviewInfo.WithContext(ctx).
			Where(tx.ReviewInfo.ReviewID.Eq(reply.ReviewID)).
			Update(tx.ReviewInfo.HasReply, 1); err != nil {
			common.ErrorLogf("[data] SaveReply update review fail,err:%v", err)
			return err
		}
		return nil
	})
	if err != nil {
		return nil, gerr.DBFailedErr(err.Error())
	}
	return reply, nil
}

func NewReviewRepo(data *Data, logger *log.Logger) biz.ReviewRepo {
	return &ReviewRepo{
		data:   data,
		logger: logger,
	}
}
