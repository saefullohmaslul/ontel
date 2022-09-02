package postgres

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

// Module postgres
var Module = fx.Options(
	fx.Provide(NewClient),
	fx.Provide(NewCustomerPostgres),
)

// Database postgres database
type Database struct {
	*gorm.DB
}

// NewClient postgres
func NewClient() *Database {
	var (
		connection *gorm.DB
		err        error
	)

	master := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Shanghai",
		os.Getenv("DB_HOST_MASTER"),
		os.Getenv("DB_USER_MASTER"),
		os.Getenv("DB_PASS_MASTER"),
		os.Getenv("DB_NAME_MASTER"),
		os.Getenv("DB_PORT_MASTER"),
		os.Getenv("DB_SSL_MASTER"),
	)

	if connection, err = gorm.Open(postgres.Open(master), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	}); err != nil {
		return &Database{}
	}

	if os.Getenv("DB_HOST_SLAVE") != "" {
		slave := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Shanghai",
			os.Getenv("DB_HOST_SLAVE"),
			os.Getenv("DB_USER_SLAVE"),
			os.Getenv("DB_PASS_SLAVE"),
			os.Getenv("DB_NAME_SLAVE"),
			os.Getenv("DB_PORT_SLAVE"),
			os.Getenv("DB_SSL_SLAVE"),
		)

		if err = connection.Use(dbresolver.Register(dbresolver.Config{
			Replicas: []gorm.Dialector{postgres.Open(slave)},
			Policy:   dbresolver.RandomPolicy{},
		})); err != nil {
			return &Database{}
		}
	}

	sqlDB, err := connection.DB()

	if _, found := os.LookupEnv("DB_MAX_OPEN_CONNECTION"); found {
		if maxOpenConnection, err := strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONNECTION")); err == nil {
			sqlDB.SetMaxOpenConns(maxOpenConnection)
		}
	}

	if _, found := os.LookupEnv("DB_MAX_IDLE_CONNECTION"); found {
		if maxIdleConnection, err := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNECTION")); err == nil {
			sqlDB.SetMaxIdleConns(maxIdleConnection)
		}
	}

	if _, found := os.LookupEnv("DB_MAX_LIFETIME_IN_MINUTE"); found {
		if maxLifetime, err := strconv.Atoi(os.Getenv("DB_MAX_LIFETIME_IN_MINUTE")); err == nil {
			sqlDB.SetConnMaxLifetime(time.Duration(maxLifetime) * time.Minute)
		}
	}

	return &Database{DB: connection}
}
