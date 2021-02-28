package service

import (
	"time"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/google/uuid"
	"github.com/kitt3911/office-hours-checker-bot/model"
	"gorm.io/gorm"
)

func Come(userId int, bot *tgbotapi.BotAPI, database *gorm.DB) model.Day{
	date := time.Now()

	var month model.Month
	var day model.Day
	
	if date.Day() == 1 {
		database.Create(&model.Month{
			ID:       uuid.Must(uuid.NewRandom()),
			UserID:   userId,
			SumHours: 0,
			Name:     date.Month().String(),
		})
	}
	database.Last(&month)
	database.Last(&day)

	dayTime := time.Time(day.Come)

	if dayTime.Day() == date.Day() && date.Month() == dayTime.Month() {
		day.Hours = float64(day.Go.Hour()) - float64(day.Come.Hour())
		day.Come = date
		database.Save(&day)
	} else {
		database.Create(&model.Day{
			ID:        uuid.Must(uuid.NewRandom()),
			UserID:    userId,
			MonthID:   month.ID,
			DayOfWeek: date.Weekday().String(),
			Hours:     0,
			Come:      date,
		})
	}

	database.Last(&day)
	return day
}
