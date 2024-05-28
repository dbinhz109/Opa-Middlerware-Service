package database

import (
	"crypto/tls"
	"go-app/src/logger"
	"net/http"
	"sync"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-redis/redis/v9"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" //import postgres
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	globalDB    *sqlx.DB      = nil
	redisClient *redis.Client = nil
	doOnce      sync.Once
	doOnceRedis sync.Once
)

func GetDbInstance() *sqlx.DB {
	if globalDB != nil {
		return globalDB
	}
	doOnce.Do(connectDb)
	return globalDB
}

func GetDbTransaction() (*sqlx.Tx, error) {
	db := GetDbInstance()
	tx, err := db.Beginx()
	if err != nil {
		logger.Warn("GetDbTransaction failed", zap.Error(err))
	}
	return tx, err
}

func connectDb() {
	var err error
	connString := viper.GetString("database.postgresql.connection-string")
	logger.Info("dbconn", zap.String("connString", connString))
	globalDB, err = sqlx.Open("postgres", connString)
	if err != nil {
		logger.Error("DBConnectionError", zap.Error(err))
	}
	globalDB.SetMaxOpenConns(100)
	globalDB.SetMaxIdleConns(5)
}

func GetRedisClient() *redis.Client {
	doOnceRedis.Do(connectRedis)
	return redisClient
}

func connectRedis() {
	viper.SetDefault("database.redis.host", "localhost")
	viper.SetDefault("database.redis.port", "6379")
	viper.SetDefault("database.redis.password", "")
	viper.SetDefault("", 0)

	redisClient = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("database.redis.host") + ":" + viper.GetString("database.redis.port"),
		Password: viper.GetString("database.redis.password"),
		DB:       viper.GetInt("database.redis.db"),
	})
}

func GetESLogClient() (*elasticsearch.Client, error) {
	viper.SetDefault("services.eslogdata.user", "")
	viper.SetDefault("services.eslogdata.password", "")
	viper.SetDefault("services.eslogdata.addresses", "http://localhost:9200")

	cfg := elasticsearch.Config{
		Username:  viper.GetString("services.eslogdata.user"),
		Password:  viper.GetString("services.eslogdata.password"),
		Addresses: viper.GetStringSlice("services.eslogdata.addresses"),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		logger.Warn("CreateESClientFailed", zap.Error(err))
	}
	return client, err
}
