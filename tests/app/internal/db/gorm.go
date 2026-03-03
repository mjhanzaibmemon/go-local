package db

import (
	"fmt"
	"time"

	"gorm.io/gorm/logger"

	mysqlbase "github.com/go-sql-driver/mysql"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	gormopentracing "gorm.io/plugin/opentracing"
)

// Config defines connection and pooling parameters for creating a GORM DB instance.
type Config struct {
	DSN             string
	ConnMaxLifetime time.Duration
	MaxOpenConns    int
	MaxIdleConns    int
	Plugins         []gorm.Plugin
}

// NewConnection opens a new GORM connection using the provided Config and registers metrics.
func NewConnection(cfg *Config) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: logger.Discard,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open gorm DB connection: %w", err)
	}

	for _, p := range cfg.Plugins {
		if errUse := db.Use(p); errUse != nil {
			return nil, fmt.Errorf("failed to use gorm plugin %s: %w", p.Name(), errUse)
		}
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get SQL DB: %w", err)
	}

	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)

	dsnConfig, err := mysqlbase.ParseDSN(cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database DSN: %w", err)
	}

	if dsnConfig.DBName != "" {
		prometheus.Unregister(collectors.NewDBStatsCollector(sqlDB, dsnConfig.DBName))

		err = prometheus.Register(collectors.NewDBStatsCollector(sqlDB, dsnConfig.DBName))
		if err != nil {
			return nil, fmt.Errorf("failed to register prometheus DB stats collector: %w", err)
		}
	}

	return db, nil
}

// NewGORMOpentracingPlugin creates a GORM plugin that instruments operations with OpenTracing.
func NewGORMOpentracingPlugin(tracer opentracing.Tracer) gorm.Plugin {
	return gormopentracing.New(
		gormopentracing.WithTracer(tracer),
		gormopentracing.WithSqlParameters(false),
		gormopentracing.WithCreateOpName("gorm_create"),
		gormopentracing.WithUpdateOpName("gorm_update"),
		gormopentracing.WithQueryOpName("gorm_query"),
		gormopentracing.WithDeleteOpName("gorm_delete"),
		gormopentracing.WithRowOpName("gorm_row"),
		gormopentracing.WithRawOpName("gorm_raw"),
	)
}
