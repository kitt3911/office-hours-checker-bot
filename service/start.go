package service

import (
	"time"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/google/uuid"
	"github.com/kitt3911/office-hours-checker-bot/model"
	"gorm.io/gorm"
)

func Start(userId int,bot *tgbotapi.BotAPI, database *gorm.DB) model.Day{
	date := time.Now()

	var month model.Month
	var day   model.Day

	database.Last(&month)
	database.Last(&day)

	if month.Name == date.Month().String() {
		dayTime := time.Time(day.Come)
		if dayTime.Day() == date.Day() && dayTime.Month() == date.Month() {
			day.Come = date
			day.Hours = float64(dayTime.Hour()) + day.Hours
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
	} else {

		uuidMonth := uuid.Must(uuid.NewRandom())
		database.Create(&model.Month{
			ID:       uuidMonth,
			UserID:   userId,
			SumHours: 0,
			Name:     date.Month().String(),
		})

		database.Create(&model.Day{
			ID:        uuid.Must(uuid.NewRandom()),
			MonthID:   uuidMonth,
			UserID:    userId,
			DayOfWeek: date.Weekday().String(),
			Hours:     0,
			Come:      date,
		})
	}

	database.Last(&day)

	return day
}
