package models

import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
	"go_short/conf"
	"fmt"
)

var DB *gorm.DB

func InitDatabase() (*gorm.DB, error) {
	config := conf.Conf()
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.DB.Host, config.DB.Port, config.DB.User, config.DB.Password, config.DB.DB_name)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        return nil, err
    }
    
    // 將數據庫連接賦值給全局變數 DB
    DB = db
	db.AutoMigrate(&UrlMapping{})
	return db, nil
}