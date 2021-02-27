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
	Come        = "come"
	Go          = "go"
	Show        = "show"
	ShowAllInfo = "show_all_info"
	Edit        = "edit"
	Delete      = "delete"
)

var numericKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("/"+Come),
		tgbotapi.NewKeyboardButton("/"+Go),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("/"+Show),
		tgbotapi.NewKeyboardButton("/"+ShowAllInfo),
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
	//bot.Debug = true
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
				database.Last(&month)
				database.Last(&day)

				dayTime := time.Time(day.Come)

				if dayTime.Day() == date.Day() && date.Month() == dayTime.Month() {
					day.Hours = float64(day.Go.Hour())-float64(day.Come.Hour())
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

				msg.Text = fmt.Sprint(day)

			case Go:
				date := time.Now()
				database.Last(&day)
				day.Go = date
				day.Hours = float64(day.Go.Hour()) - float64(day.Come.Hour())
				database.Save(&day)
				msg.Text = fmt.Sprint(day.Hours)

			case ShowAllInfo:

				database.Last(&month)
				var days []model.Day

				database.Find(&days, "month_id = ?", month.ID)

				var outputStr string
				outputStr = outputStr + "| Day \t  \t  \t| Hours|\n"
				for _, item := range days {
					dateFormate := item.Come.Format("2006-01-02")
					outputStr = outputStr + fmt.Sprintf(`| %s		| %d|`, dateFormate, int(item.Hours)) + "\n"
				}

				msg.Text = fmt.Sprint(outputStr)

			case Show:
				database.Last(&day)
				msg.Text = fmt.Sprint(day.Hours)

			default:
				msg.Text = "I don't know that command"
			}
			bot.Send(msg)
		}
	}

}
