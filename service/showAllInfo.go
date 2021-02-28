package service

import (
	"fmt"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/kitt3911/office-hours-checker-bot/model"
	"gorm.io/gorm"
)

func ShowAllInfo(bot *tgbotapi.BotAPI, database *gorm.DB) string {

	var month model.Month
	database.Last(&month)
	var days []model.Day

	database.Find(&days, "month_id = ?", month.ID)
	var outputStr string
	outputStr = outputStr + "| Day \t  \t  \t| Hours|\n"
	for _, item := range days {
		dateFormate := item.Come.Format("2006-01-02")
		outputStr = outputStr + fmt.Sprintf(`| %s		| %d|`, dateFormate, int(item.Hours)) + "\n"
	}

	return outputStr
}
