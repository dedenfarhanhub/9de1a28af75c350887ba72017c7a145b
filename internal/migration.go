package internal

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file" // Needed for file-based migrations
	"gorm.io/gorm"
	"log"
	"strconv"
)

// RunMigrations run migration script
func RunMigrations(db *gorm.DB) error {
	// Get *sql.DB from GORM
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get *sql.DB from GORM: %v", err)
	}

	// Initialize the MySQL driver
	driver, err := mysql.WithInstance(sqlDB, &mysql.Config{})
	if err != nil {
		return fmt.Errorf("could not create MySQL driver: %v", err)
	}

	// Initialize the migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations", // Path to the migrations directory
		"mysql", driver,
	)
	if err != nil {
		return fmt.Errorf("could not initialize migrate: %v", err)
	}

	// Run the migrations
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("could not run up migrations: %v", err)
	}

	log.Println("Migrations ran successfully")
	return nil
}

// RollbackMigrations rollback migration script
func RollbackMigrations(db *gorm.DB, stepsStr string) error {
	// Get *sql.DB from GORM
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get *sql.DB from GORM: %v", err)
	}

	// Initialize the MySQL driver
	driver, err := mysql.WithInstance(sqlDB, &mysql.Config{})
	if err != nil {
		return fmt.Errorf("could not create MySQL driver: %v", err)
	}

	// Initialize the migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations", // Path to the migrations directory
		"mysql", driver,
	)
	if err != nil {
		return fmt.Errorf("could not initialize migrate: %v", err)
	}

	// Parse the steps argument
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return fmt.Errorf("invalid steps: %v", err)
	}

	// Rollback the migrations
	if err := m.Steps(-steps); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("could not roll back migrations: %v", err)
	}

	log.Println("Rollback ran successfully")
	return nil
}
