package main

import (
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/google/uuid"
	"github.com/kitt3911/office-hours-checker-bot/config"
	"github.com/kitt3911/office-hours-checker-bot/database"
	"github.com/kitt3911/office-hours-checker-bot/model"
)

var (
	Come      = "come"
	Go        = "go"
	Show      = "show"
	ShowMonth = "show_month"
	Edit      = "edit"
	Delete    = "delete"
)

var numericKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("/"+Come),
		tgbotapi.NewKeyboardButton("/"+Go),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("/"+Show),
		tgbotapi.NewKeyboardButton("/"+ShowMonth),
		tgbotapi.NewKeyboardButton("/"+Delete),
		tgbotapi.NewKeyboardButton("/"+Edit),
	),
)

func main() {

	config := config.ConfigNew()
	database, err := database.InitDatabase(config)

	if err != nil {
		log.Fatalln(err)
	}

	bot, err := tgbotapi.NewBotAPI(config.BotToken)
	if err != nil {
		log.Fatalln(err)
	}
	bot.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)

	fmt.Println("Service has been started!")
	for update := range updates {

		if update.CallbackQuery != nil {

			bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))
			bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data))
		}
		if update.Message != nil {

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			userId := update.Message.From.ID

			var month model.Month
			var day model.Day

			switch update.Message.Command() {
			case "start":
				msg.ReplyMarkup = numericKeyboard
				date := time.Now()

				database.First(&month)
				database.First(&day)

				if month.Name == date.Month().String() {
					dayTime := time.Time(day.Come)
					if dayTime.Day() == date.Day() && dayTime.Month() == date.Month() {
						day.Come = date
						day.Hours = 0
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

				msg.Text = fmt.Sprint(month)

			case Come:
				date := time.Now()
				uuidMonth := uuid.Must(uuid.NewRandom())

				if date.Day() == 1 {
					database.Create(&model.Month{
						ID:       uuidMonth,
						UserID:   userId,
						SumHours: 0,
						Name:     date.Month().String(),
					})
				}

				database.First(&month)
				database.First(&day)

				dayTime := time.Time(day.Come)

				if dayTime.Day() == date.Day() && dayTime.Month() == date.Month() {
					day.Come = date
					day.Hours = 0
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

				msg.Text = fmt.Sprint(userId)

			case Go:
				date := time.Now()
				database.First(&day)
				workHour := date.Sub(day.Come)
				day.Go = date
				day.Hours = float64(workHour.Minutes()) / 60
				database.Save(&day)
				msg.Text = fmt.Sprint(day.Hours)
			case Show:

				var day []model.Day
				var count int64
				database.Find(&day).Count(&count)
				msg.Text = fmt.Sprint(count)
			default:
				msg.Text = "I don't know that command"
			}
			bot.Send(msg)
		}
	}

}
