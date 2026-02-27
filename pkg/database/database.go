package database

import (
	"fmt"
	"hyperscal-go/config"
	"hyperscal-go/internal/domain"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ConnectDatabase establishes database connection based on configuration
func ConnectDatabase(cfg *config.Config) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	switch cfg.Database.Driver {
	case "postgres":
		db, err = connectPostgres(cfg.Database.Postgres, gormConfig)
	case "oracle":
		db, err = connectOracle(cfg.Database.Oracle, gormConfig)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Database.Driver)
	}

	if err != nil {
		return nil, err
	}

	log.Printf("Successfully connected to %s database", cfg.Database.Driver)

	// Auto migrate schemas
	if err := autoMigrate(db); err != nil {
		return nil, err
	}

	return db, nil
}

// connectPostgres establishes connection to PostgreSQL database
func connectPostgres(cfg config.PostgresConfig, gormConfig *gorm.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	return db, nil
}

// connectOracle establishes connection to Oracle database
func connectOracle(cfg config.OracleConfig, gormConfig *gorm.Config) (*gorm.DB, error) {
	// Note: For Oracle, you need to install Oracle Instant Client
	// DSN format: user/password@host:port/service_name
	_ = fmt.Sprintf(
		"user=\"%s\" password=\"%s\" connectString=\"%s:%s/%s\"",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.ServiceName,
	)

	// Using godror driver for Oracle
	// You'll need to import: "github.com/godror/godror"
	// and create a custom dialector or use existing gorm oracle driver

	// For now, returning an error with instructions
	// In production, you would use: gorm.Open(oracle.Open(dsn), gormConfig)
	return nil, fmt.Errorf("Oracle support requires additional setup. Please install Oracle Instant Client and configure godror driver")
}

// autoMigrate runs database migrations
func autoMigrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	if err := db.AutoMigrate(
		&domain.Country{},
		&domain.City{},
		&domain.User{},
		// Add other models here
	); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}
