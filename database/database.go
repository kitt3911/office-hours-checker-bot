package database

import (
	"fmt"

	"github.com/kitt3911/office-hours-checker-bot/config"
	"github.com/kitt3911/office-hours-checker-bot/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDatabase(config *config.Config) (*gorm.DB, error) {
	var database *gorm.DB
	connect := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s",
		"localhost", 5400, config.DB.User, config.DB.Password, config.DB.Name)

	database, err := gorm.Open(postgres.Open(connect), &gorm.Config{})

	database.Table("days").AutoMigrate(&model.Day{})
	database.Table("months").AutoMigrate(&model.Month{})

	database.Model(&model.Month{}).Association("days")
	if err != nil {
		return database, err
	}

	return database, nil

}