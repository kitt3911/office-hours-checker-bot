package main

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Day struct {
	UserID    string `json:"user_id" gorm:"primary_key"`
	DayOfWeek string `json:"day_od_week"`
	Date      string `json:"date"`
	Hours     int    `json:"hours"`
	Minutes   int    `json:"minutes"`
}

type Week struct {
	UserID   string `json:"user_id" gorm:"primary_key"`
	Days     []Day  `json:"days" gorm:"foreignKey:UserID"`
	SumHours int    `json:"sum_hours"`
}

type Month struct {
	gorm.Model
	UserID   string `json:"user_id" gorm:"primary_key"`
	Weeks    []Week `json:"weeks" gorm:"foreignKey:UserID"`
	SumHours int    `json:"sum_hours"`
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

var (
	Database *gorm.DB
)

var (
	Come      = "/Пришел"
	Go        = "/Ушел"
	Show      = "/Показать за неделю"
	ShowMonth = "/Показать за месяц"
	Add       = "/Добавить долг"
)

func InitDatabase(config *Config) (*gorm.DB, error) {
	connect := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s",
		"localhost", 5432, config.DB.User, config.DB.Password, config.DB.Name)

	Database, err := gorm.Open(postgres.Open(connect), &gorm.Config{})

	Database.Table("Days").AutoMigrate(&Day{})

	Database.Table("Weeks").AutoMigrate(&Week{})
	Database.Table("Months").AutoMigrate(&Month{})

	Database.Model(&Month{}).Association("Weeks")
	Database.Model(&Week{}).Association("Days")
	if err != nil {
		return Database, err
	}

	return Database, nil

}

func main() {

	config := ConfigNew()
	InitDatabase(config)


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
			bot.Send(msg)
		}
	}
	
}
