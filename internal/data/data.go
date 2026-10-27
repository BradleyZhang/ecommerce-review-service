package data

import (
	"io"
	"log"
	"os"
	"review-service/common"
	"review-service/internal/data/query"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	return gorm.Open(mysql.Open(dsn), newGormConfig(true))
}
func newGormConfig(prepareStmt bool) *gorm.Config {
	return &gorm.Config{
		PrepareStmt: prepareStmt,
		Logger:      newGormLogger(os.Stdout),
	}
}
func newGormLogger(w io.Writer) logger.Interface {
	return logger.New(log.New(w, "\n[GORM]", log.LstdFlags), logger.Config{
		// SlowThreshold:             time.Duration(slowThresholdMs) * time.Millisecond,
		LogLevel:                  logger.Warn,
		IgnoreRecordNotFoundError: true,
		ParameterizedQueries:      !common.DebugMode,
		Colorful:                  true,
	})
}
