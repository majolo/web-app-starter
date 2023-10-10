package database

import (
	"context"
	"fmt"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

type PostgresConnConfig struct {
	Address    Address
	Auth       Auth
	TLSConfig  TLSConfig
	PoolConfig PoolConfig
}

type Address struct {
	Host   string
	Port   string
	DBName string
}

type Auth struct {
	User     string
	Password string
}

type TLSConfig struct {
	TLSMode    string
	CACertPath string
	KeyPath    string
	CertPath   string
}

type PoolConfig struct {
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

func PostgresGormDB(config PostgresConnConfig) (*gorm.DB, error) {
	// Create the connection string.
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Address.Host,
		config.Auth.User,
		config.Auth.Password,
		config.Address.DBName,
		config.Address.Port,
		config.TLSConfig.TLSMode,
	)

	if config.TLSConfig.TLSMode == "verify-ca" || config.TLSConfig.TLSMode == "verify-full" {
		_, err := os.Stat(config.TLSConfig.CACertPath)
		if err != nil {
			return nil, fmt.Errorf("cannot use TLS without CA cert path")
		}
		dsn = fmt.Sprintf("%s sslrootcert=%s", dsn, config.TLSConfig.CACertPath)

		_, err = os.Stat(config.TLSConfig.CertPath)
		certPathExists := err == nil
		_, err = os.Stat(config.TLSConfig.KeyPath)
		keyPathExists := err == nil

		if certPathExists != keyPathExists {
			return nil, fmt.Errorf("cannot use mTLS without both key and cert for client")
		} else if certPathExists {
			dsn = fmt.Sprintf("%s sslcert=%s", dsn, config.TLSConfig.CertPath)
			dsn = fmt.Sprintf("%s sslkey=%s", dsn, config.TLSConfig.KeyPath)
		}
	}

	gormConfig := &gorm.Config{}

	// Open the connection.
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, err
	}

	// Set connection pool config.
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if config.PoolConfig.MaxIdleConns != 0 {
		sqlDB.SetMaxIdleConns(config.PoolConfig.MaxIdleConns)
	}
	if config.PoolConfig.MaxOpenConns != 0 {
		sqlDB.SetMaxOpenConns(config.PoolConfig.MaxOpenConns)
	}
	if config.PoolConfig.ConnMaxLifetime != 0 {
		sqlDB.SetConnMaxLifetime(config.PoolConfig.ConnMaxLifetime)
	}

	// Check database connection is active.
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	err = sqlDB.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Migrate(db *gorm.DB, migrations []*gormigrate.Migration) error {
	err := migrate(db, migrations)
	if err != nil {
		log.Println("migrate error: ", err)
	}
	return nil
}

func migrate(db *gorm.DB, migrations []*gormigrate.Migration) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, migrations)
	return m.Migrate()
}
