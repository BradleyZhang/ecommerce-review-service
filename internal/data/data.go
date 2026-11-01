package data

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/BradleyZhang/ecommerce-review-service/common"
	"github.com/BradleyZhang/ecommerce-review-service/internal/data/query"

	es "github.com/elastic/go-elasticsearch/v9"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Data struct {
	query  *query.Query
	logger *log.Logger // log *log.Helper
	es     *es.TypedClient
	rdb    *redis.Client
}

func NewData(db *gorm.DB, es *es.TypedClient, rdb *redis.Client, logger *log.Logger) (*Data, func(context.Context), error) {
	data := &Data{query: query.Q, logger: logger, es: es, rdb: rdb}
	cleanup := func(ctx context.Context) {
		logger.Println("closing the data resources")
		err := data.es.Close(ctx)
		if err != nil {
			common.ErrorLogf("failed to close es :%v", err)
		} else {
			common.SysLog("close es...")
		}
		err = rdb.Close()
		if err != nil {
			common.ErrorLogf("failed to close redis :%v", err)
		} else {
			common.SysLog("close redis...")
		}
		sqlDB, err := db.DB()
		if err != nil {
			common.ErrorLogf("get sql.DB failed: %v", err)
			return
		} else {
			err = sqlDB.Close()
			if err != nil {
				common.ErrorLogf("failed to close db :%v", err)
			} else {
				common.SysLog("close db...")
			}

		}
	}
	query.SetDefault(db)
	return data, cleanup, nil
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

type RedisConfig struct {
	Address  string
	Password string
	DB       int
}

func NewRedisCli(cfg RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return rdb, nil
}
func NewRedisConfigDefault() RedisConfig {
	return RedisConfig{
		Address:  common.RedisAddress,
		Password: common.RedisPassword,
		DB:       common.RedisDB,
	}
}
