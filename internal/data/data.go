package data

import (
	"log"
	"review-service/common"
	"review-service/internal/data/query"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Data struct {
	query  *query.Query
	logger *log.Logger // log *log.Helper
}

func NewData(db *gorm.DB, logger *log.Logger) (*Data, func(), error) {
	cleanup := func() {
		logger.Println("closing the data resources")
	}
	query.SetDefault(db)
	return &Data{query: query.Q, logger: logger}, cleanup, nil

}

func NewDB(dsn string) (*gorm.DB, error) {
	common.SysLog("using mysql(only support mysql now)")
	return gorm.Open(mysql.Open(dsn))
}
