package data

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sync/singleflight"

	"github.com/BradleyZhang/ecommerce-review-service/common"
	"github.com/BradleyZhang/ecommerce-review-service/gerr"
	"github.com/BradleyZhang/ecommerce-review-service/internal/biz"
	"github.com/BradleyZhang/ecommerce-review-service/internal/data/model"
	"github.com/BradleyZhang/ecommerce-review-service/internal/data/query"
	"github.com/elastic/go-elasticsearch/v9/typedapi/esdsl"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ReviewRepo struct {
	data   *Data
	logger *log.Logger
}

func NewReviewRepo(data *Data, logger *log.Logger) biz.ReviewRepo {
	return &ReviewRepo{
		data:   data,
		logger: logger,
	}
}

// AppealReview implements [biz.ReviewRepo].
func (r *ReviewRepo) AppealReview(ctx context.Context, param *biz.AppealParam) (*model.ReviewAppealInfo, error) {
	// 申诉过且申诉已审核的，不许再申诉
	appeal := r.data.query.ReviewAppealInfo
	existAppeal, err := appeal.WithContext(ctx).Where(
		appeal.ReviewID.Eq(param.ReviewID),
		appeal.StoreID.Eq(param.StoreID),
	).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, gerr.DBFailedErr(err.Error())
	}
	if err == nil && existAppeal.Status > 10 {
		return nil, gerr.DuplicateAppealErr("该评论已有审核过的申诉记录")
	}

	// 申诉过而没有审核过，更新申诉信息
	// if
	panic("unimplemented")
}

// AuditAppeal implements [biz.ReviewRepo].
func (r *ReviewRepo) AuditAppeal(ctx context.Context, param *biz.AuditAppealParam) error {
	panic("unimpl")
}

// AuditReview implements [biz.ReviewRepo].
func (r *ReviewRepo) AuditReview(ctx context.Context, param *biz.AuditParam) error {
	review := r.data.query.ReviewInfo
	_, err := review.WithContext(ctx).Where(review.ReviewID.Eq(param.ReviewID)).Updates(
		map[string]any{
			"status":     param.Status,
			"op_user":    param.OpUser,
			"op_reason":  param.OpReason,
			"op_remarks": param.OpRemarks,
		})
	return err
}

// GetReview implements [biz.ReviewRepo].
func (r *ReviewRepo) GetReview(ctx context.Context, reviewID int64) (*model.ReviewInfo, error) {
	info := r.data.query.ReviewInfo
	review, err := info.WithContext(ctx).Where(info.ReviewID.Eq(reviewID)).First()
	if err != nil {
		return nil, gerr.DBFailedErr(err.Error())
	}
	return review, nil
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

// GetReviewReply implements [biz.ReviewRepo].
func (r *ReviewRepo) GetReviewReply(ctx context.Context, reviewID int64) (*model.ReviewReplyInfo, error) {
	reply := r.data.query.ReviewReplyInfo
	ans, err := reply.WithContext(ctx).Where(reply.ReviewID.Eq(reviewID)).First()
	if err != nil {
		return nil, gerr.DBFailedErr(err.Error())
	}
	return ans, nil
}

// ListReviewByUserID implements [biz.ReviewRepo].
func (r *ReviewRepo) ListReviewByUserID(ctx context.Context, userID int64, offset int, limit int) ([]*model.ReviewInfo, error) {
	common.DebugLogf("[data ListReviewByUserID] userID:%v offset:%v limit:%v", userID, offset, limit)
	data, err := r.getReviewByFieldSingleFlignt(ctx, offset, limit, "user_id", strconv.FormatInt(userID, 10))
	// data, err := r.esGetReviewByField(ctx, offset, limit, "user_id", strconv.FormatInt(userID, 10))
	if err != nil {
		return nil, fmt.Errorf("es get review err:%w", err)
	}
	hm := &types.HitsMetadata{}
	if err := json.Unmarshal(data, hm); err != nil {
		return nil, fmt.Errorf("failed to unmarshal hitsmetadata err:%w", err)
	}

	list := make([]*model.ReviewInfo, 0, hm.Total.Value)
	for _, hit := range hm.Hits {
		tmp := &esReviewInfo{}
		if err := json.Unmarshal(hit.Source_, tmp); err != nil {
			common.ErrorLogf("[data] failed to unmarshal review info err:%v source:%+v", err, string(data))
			continue
		}
		list = append(list, tmp.ToReviewInfo())
	}
	return list, nil
	// panic("unimplemented")
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

// SaveReview implements [biz.ReviewRepo].
func (r *ReviewRepo) SaveReview(ctx context.Context, review *model.ReviewInfo) (*model.ReviewInfo, error) {
	err := r.data.query.ReviewInfo.WithContext(ctx).Save(review)
	if err != nil {
		return nil, gerr.DBFailedErr(err.Error())
	}
	return review, nil
}

var g = new(singleflight.Group)

func (r *ReviewRepo) getReviewByFieldSingleFlignt(ctx context.Context, offset, limit int, field, value string) ([]byte, error) {
	key := fmt.Sprintf("review:%s:%s:%d:%d", field, value, offset, limit)
	v, err, _ := g.Do(key, func() (any, error) {
		// 查缓存
		data, err := r.getBytesFromCache(ctx, key)
		if err == nil {
			return data, nil
		}
		if errors.Is(err, redis.Nil) { // 未命中, 查 es 更新缓存
			data, err := r.esGetReviewByField(ctx, offset, limit, field, value)
			if err == nil {
				return data, r.setCache(ctx, key, data)
			} else { // 查es失败
				return nil, err
			}
		}
		// 查缓存失败
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return v.([]byte), nil
}

func (r *ReviewRepo) getBytesFromCache(ctx context.Context, key string) ([]byte, error) {
	return r.data.rdb.Get(ctx, key).Bytes()
}
func (r *ReviewRepo) setCache(ctx context.Context, key string, data []byte) error {
	return r.data.rdb.Set(ctx, key, data, time.Second*10).Err()
}
func (r *ReviewRepo) esGetReviewByField(ctx context.Context, offset, limit int, field, value string) ([]byte, error) {
	common.DebugLogf("[data esGetReviewByField] field:%v value: %v offset:%v limit:%v", field, value, offset, limit)
	resp, err := r.data.es.Search().
		Index(common.ESIndex).
		From(offset).
		Size(limit).
		Query(
			esdsl.NewBoolQuery().
				Filter(
					esdsl.NewTermQuery(
						field,
						esdsl.NewFieldValue().String(value),
					),
				),
		).
		Do(ctx)

	common.DebugLogf("[data esGetReviewByField] es resp:%+v", resp)
	if err != nil {
		return nil, err
	}
	return json.Marshal(resp.Hits)
}

type esTime time.Time

var _ json.Unmarshaler = (*esTime)(nil) // 检查实现 json.Unmarshaler

// implements [encoding/json.Unmarshaler].
func (t *esTime) UnmarshalJSON(data []byte) error {
	// 2026-10-31 05:58:46
	s := strings.Trim(string(data), `"`)
	temp, err := time.Parse(time.DateTime, s)
	if err != nil {
		return err
	}
	*t = esTime(temp)
	return nil
}

type esReviewInfo struct {
	*model.ReviewInfo
	CreateAt     esTime `json:"create_at"`
	UpdateAt     esTime `json:"update_at"`
	ID           int64  `json:"id,string"`
	Version      int32  `json:"version,string"`
	ReviewID     int64  `json:"review_id,string"`
	ServiceScore int32  `json:"service_score,string"`
	ExpressScore int32  `json:"express_score,string"`
	HasMedia     int32  `json:"has_media,string"`
	OrderID      int64  `json:"order_id,string"`
	SkuID        int64  `json:"sku_id,string"`
	SpuID        int64  `json:"spu_id,string"`
	StoreID      int64  `json:"store_id,string"`
	UserID       int64  `json:"user_id,string"`
	Anonymous    int32  `json:"anonymous,string"`
	Status       int32  `json:"status,string"`
	Score        int32  `json:"score,string"`
	IsDefault    int32  `json:"is_default,string"`
	HasReply     int32  `json:"has_reply,string"`
}

func (r *esReviewInfo) ToReviewInfo() *model.ReviewInfo {
	info := r.ReviewInfo

	info.CreateAt = time.Time(r.CreateAt)
	info.UpdateAt = time.Time(r.UpdateAt)
	info.ID = r.ID
	info.Version = r.Version
	info.ReviewID = r.ReviewID
	info.ServiceScore = r.ServiceScore
	info.ExpressScore = r.ExpressScore
	info.HasMedia = r.HasMedia
	info.OrderID = r.OrderID
	info.SkuID = r.SkuID
	info.SpuID = r.SpuID
	info.StoreID = r.StoreID
	info.UserID = r.UserID
	info.Anonymous = r.Anonymous
	info.Status = r.Status
	info.Score = r.Score
	info.IsDefault = r.IsDefault
	info.HasReply = r.HasReply

	return info
}
