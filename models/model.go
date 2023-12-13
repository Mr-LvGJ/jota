package models

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	DefaultDBMaxIdleConns    = 10
	DefaultDBMaxOpenConns    = 100
	DefaultDBConnMaxLifetime = 120
	DefaultDBConnMaxIdleTime = 120
)

type BaseModel struct {
	Id        string         `gorm:"type:varchar(255);not null;primaryKey"`
	CreatedAt time.Time      `gorm:"type:datetime(6);not null;index:created_at;default:current_timestamp(6)"`
	UpdatedAt time.Time      `gorm:"type:datetime(6);not null;index:updated_at;default:current_timestamp(6) on update current_timestamp(6)"`
	DeletedAt gorm.DeletedAt `gorm:"type:datetime(6);index:deleted_at"`
}

type DatabaseConfig struct {
	User                  string
	Password              string
	Host                  string
	Schema                string
	PrepareStmt           bool
	SetDefaultTransaction bool
	SystemVars            map[string]any

	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime int
	ConnMaxIdleTime int

	EnableLog bool
	Logger    logger.Interface
}

func New(config *DatabaseConfig) (*gorm.DB, error) {
	url := formatDBUrl(config.User, config.Password, config.Host, config.Schema, config.SystemVars)

	if !config.EnableLog {
		config.Logger = logger.Discard
	}

	conn, err := gorm.Open(mysql.Open(url), &gorm.Config{
		Logger:                                   config.Logger,
		PrepareStmt:                              config.PrepareStmt,
		SkipDefaultTransaction:                   !config.SetDefaultTransaction,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, err
	}

	if config.MaxIdleConns == 0 {
		config.MaxIdleConns = DefaultDBMaxIdleConns
	}
	if config.MaxOpenConns == 0 {
		config.MaxOpenConns = DefaultDBMaxIdleConns
	}
	if config.ConnMaxLifetime == 0 {
		config.ConnMaxLifetime = DefaultDBConnMaxLifetime
	}
	if config.ConnMaxIdleTime == 0 {
		config.ConnMaxIdleTime = DefaultDBConnMaxIdleTime
	}

	db, err := conn.DB()
	if err != nil {
		return nil, err
	}

	// connection pool settings
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetConnMaxLifetime(time.Duration(config.ConnMaxLifetime) * time.Second)
	db.SetConnMaxIdleTime(time.Duration(config.ConnMaxIdleTime) * time.Second)
	return conn, nil
}

func formatDBUrl(user, password, host, schema string, systemVars map[string]interface{}) string {
	url := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		user, password, host, schema)
	for key, value := range systemVars {
		url = fmt.Sprintf("%s&%s=%v", url, key, value)
	}
	return url
}
