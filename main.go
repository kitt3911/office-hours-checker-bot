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
	UserID    int     `json:"user_id"`
	WeekID    uuid.UUID `json:"week_id"`
	DayOfWeek string  `json:"day_od_week"`
	Date      string  `json:"date"`
	Hours     int     `json:"hours"`
	Minutes   int     `json:"minutes"`
}

type Week struct {
	ID       uuid.UUID `json:"id" gorm:"primaryKey"`
	MonthID  uuid.UUID `json:"month_id`
	UserID   int     `json:"user_id"`
	Days     []Day   `json:"days" gorm:"foreignKey:WeekID"`
	SumHours int     `json:"sum_hours"`
	Number   int     `json:"number"`
}

type Month struct {
	gorm.Model
	ID       uuid.UUID `json:"id" gorm:"primaryKey"`
	UserID   int     `json:"user_id"`
	Weeks    []Week  `json:"weeks" gorm:"foreignKey:MonthID"`
	SumHours int     `json:"sum_hours"`
	Name     string  `json:"name_of_month"`
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
	Come      = "/Пришел"
	Go        = "/Ушел"
	Show      = "/Показать за неделю"
	ShowMonth = "/Показать за месяц"
	Edit      = "/Редактировать"
	Delete    = "/Удалить"
)

var numericKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(Come),
		tgbotapi.NewKeyboardButton(Go),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(Show),
		tgbotapi.NewKeyboardButton(ShowMonth),
		tgbotapi.NewKeyboardButton(Delete),
		tgbotapi.NewKeyboardButton(Edit),
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

	database.Table("Days").AutoMigrate(&Day{})
	database.Table("Weeks").AutoMigrate(&Week{})
	database.Table("Months").AutoMigrate(&Month{})


	database.Model(&Week{}).Association("Days")
	database.Model(&Month{}).Association("Weeks")
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
			switch update.Message.Command() {
			case "start":
				msg.ReplyMarkup = numericKeyboard
				msg.Text = fmt.Sprint(userId)
				date := time.Now()
				uuidMonth := uuid.New()
				uuidWeek := uuid.New()
				database.Create(&Month{
					ID:	uuidMonth,
					UserID:   userId,
					SumHours: 0,
					Name:     date.Month().String(),
				})
				database.Create(&Week{
					ID:	uuidWeek,
					MonthID: uuidMonth,
					UserID:   userId,
					Number:   1,
					SumHours: 0,
				})
			
				database.Create(&Day{
					ID: uuid.New(),
					WeekID: uuidWeek,
					UserID:    userId,
					DayOfWeek: date.Weekday().String(),
					Date:      date.String(),
					Hours:     0,
					Minutes:   0,
				})
				
			case Come:
				date := time.Now()
				uuidWeek := uuid.New()
				uuidMonth := uuid.New()
				if date.Day() == 1 {
					database.Create(&Month{
						ID:	uuidMonth,
						UserID:   userId,
						SumHours: 0,
						Name:     date.Month().String(),
					})
					database.Create(&Week{
						MonthID: uuidMonth,
						UserID:   userId,
						ID: uuidWeek,
						Number:   1,
						SumHours: 0,
					})
				}
				

				msg.Text = fmt.Sprint(userId)

			case Go:

				msg.Text = fmt.Sprint(userId)
			case Show:

				msg.Text = fmt.Sprint(userId)
			default:
				msg.Text = "I don't know that command"
			}
			bot.Send(msg)
		}
	}

}
