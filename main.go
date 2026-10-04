package main

import (
	"context"
	"log"
	"os"

	"tg-bot/internal/booking"
	"tg-bot/internal/database"
	"tg-bot/internal/storage"

	"github.com/joho/godotenv"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Файл .env не найден, используются переменные окружения")
	}

	if err := database.CreateDatabaseIfNotExists(); err != nil {
		log.Fatal(err)
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close(context.Background())

	log.Println("Подключение к PostgreSQL успешно")

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN не найден")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Бот запущен: %s", bot.Self.UserName)

	postgresStorage, err := storage.NewPostgresStorage(db)
	if err != nil {
		log.Fatal(err)
	}

	bookingService := booking.NewBookingService(postgresStorage)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		handleMessage(bot, update.Message, bookingService)
	}
}