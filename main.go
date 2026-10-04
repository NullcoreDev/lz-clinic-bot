package main

import (
	"log"
	"os"

	"tg-bot/internal/booking"
	"tg-bot/internal/storage"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Println("Файл .env не найден, используются переменные окружения")
	}

	token := os.Getenv("TELEGRAM_BOT_TOKEN")

	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN не найден")
	}

	bot, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Бот запущен: %s", bot.Self.UserName)

	// Хранилище записей.
	store := storage.NewMemoryStorage()

	// Сервис записи.
	bookingService := booking.NewBookingService(store)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		handleMessage(
			bot,
			update.Message,
			bookingService,
		)
	}
}
