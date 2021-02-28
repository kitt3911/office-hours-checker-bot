package main

import (
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/kitt3911/office-hours-checker-bot/config"
	"github.com/kitt3911/office-hours-checker-bot/database"
	"github.com/kitt3911/office-hours-checker-bot/model"
	"github.com/kitt3911/office-hours-checker-bot/service"
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

			var day model.Day

			switch update.Message.Command() {
			case "start":

				msg.ReplyMarkup = numericKeyboard
				day = service.Start(userId, bot, database)
				msg.Text = fmt.Sprint(day)

			case Come:

				day = service.Come(userId, bot, database)
				msg.Text = fmt.Sprint(day)

			case Go:
				date := time.Now()
				database.Last(&day)
				day.Go = date
				day.Hours = float64(day.Go.Hour()) - float64(day.Come.Hour())
				database.Save(&day)
				msg.Text = fmt.Sprint(day.Hours)

			case ShowAllInfo:

				outputStr := service.ShowAllInfo(bot, database)

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
