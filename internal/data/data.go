package data

import (
	"io"
	"log"
	"os"

	"github.com/BradleyZhang/ecommerce-review-service/common"
	"github.com/BradleyZhang/ecommerce-review-service/internal/data/query"

	es "github.com/elastic/go-elasticsearch/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Data struct {
	query  *query.Query
	logger *log.Logger // log *log.Helper
	es     *es.TypedClient
}

func NewData(db *gorm.DB, es *es.TypedClient, logger *log.Logger) (*Data, func(), error) {
	cleanup := func() {
		logger.Println("closing the data resources")
	}
	query.SetDefault(db)
	return &Data{query: query.Q, logger: logger, es: es}, cleanup, nil

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

type ESCfg struct {
	Addresses []string
	Apikey    string
}

func NewESClient(cfg ESCfg) (*es.TypedClient, error) {
	var opts []es.Option
	opts = append(opts,
		es.WithAddresses(cfg.Addresses...),
		es.WithAPIKey(cfg.Apikey),
	)

	client, err := es.NewTyped(opts...)
	if err != nil {
		return nil, err
	}
	common.SysLogf("connect to es with address: %v", cfg.Addresses)
	common.DebugLogf("[job] NewESClient cfg:%v", cfg)
	return client, nil
}
func NewESDefaultCfg() ESCfg {
	return ESCfg{
		Addresses: common.ESAddresses,
		Apikey:    common.ESAPIKey,
	}
}
