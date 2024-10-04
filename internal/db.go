package internal

import (
	"github.com/dedenfarhanhub/blog-service/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitDB initializes db
func InitDB() (*gorm.DB, error) {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database connection
	db, err := gorm.Open(mysql.Open(cfg.DBUser+":"+cfg.DBPassword+"@tcp("+cfg.DBHost+":"+cfg.DBPort+")/"+cfg.DBName+"?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
