package main

import (
	"fmt"
	"log"
	"os"
	"time"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Day struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey"`
	UserID    int       `json:"user_id"`
	MonthID   uuid.UUID `json:"month_id"`
	DayOfWeek string    `json:"day_od_week"`
	Hours     float64   `json:"hours"`
	Go        time.Time `json:"go_time"`
	Come      time.Time `json:"come_time"`
}

type Month struct {
	gorm.Model
	ID       uuid.UUID `json:"id" gorm:"primaryKey"`
	UserID   int       `json:"user_id"`
	Days     []Day     `json:"weeks" gorm:"foreignKey:MonthID"`
	SumHours int       `json:"sum_hours"`
	Name     string    `json:"name_of_month"`
}

type DatabaseConfig struct {
	User     string
	Password string
	Name     string
}

type Config struct {
	DB       DatabaseConfig
	BotToken string
}

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

func ConfigNew() *Config {
	_ = godotenv.Load()
	return &Config{
		DB: DatabaseConfig{
			User:     getEnv("DB_USER", ""),
			Password: getEnv("DB_PASS", ""),
			Name:     getEnv("DB_NAME", ""),
		},
		BotToken: getEnv("BOT_TOKEN", ""),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func InitDatabase(config *Config) (*gorm.DB, error) {
	var database *gorm.DB
	connect := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s",
		"localhost", 5400, config.DB.User, config.DB.Password, config.DB.Name)

	database, err := gorm.Open(postgres.Open(connect), &gorm.Config{})

	database.Table("days").AutoMigrate(&Day{})
	database.Table("months").AutoMigrate(&Month{})

	database.Model(&Month{}).Association("days")
	if err != nil {
		return database, err
	}

	return database, nil

}

func main() {

	config := ConfigNew()
	database, err := InitDatabase(config)

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

			var month Month
			var day Day

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
						database.Create(&Day{
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
					database.Create(&Month{
						ID:       uuidMonth,
						UserID:   userId,
						SumHours: 0,
						Name:     date.Month().String(),
					})

					database.Create(&Day{
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
					database.Create(&Month{
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
					database.Create(&Day{
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

				var day []Day
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
