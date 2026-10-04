package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Booking struct {
	Service string
	Doctor  string
	Date    string
	Time    string
}

type UserState struct {
	Step    string
	Booking Booking
}

var users = make(map[int64]*UserState)

func main() {
	// ВСТАВЬ СЮДА СВОЙ ТОКЕН
	bot, err := tgbotapi.NewBotAPI("REMOVED_SECRET")

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Бот запущен:", bot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID
		text := update.Message.Text

		if _, exists := users[chatID]; !exists {
			users[chatID] = &UserState{}
		}

		user := users[chatID]

		switch text {

		case "/start":
			user.Step = ""
			user.Booking = Booking{}

			msg := tgbotapi.NewMessage(
				chatID,
				"👋 Здравствуйте!\n\n"+
					"Я бот клиники LZ.\n\n"+
					"Чем могу помочь?",
			)

			msg.ReplyMarkup = mainKeyboard()
			bot.Send(msg)

		case "📅 Записаться на приём":
			user.Step = "service"

			keyboard := tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("🦷 Лечение зубов"),
					tgbotapi.NewKeyboardButton("🧼 Профессиональная чистка"),
				),
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("🦷 Удаление зуба"),
					tgbotapi.NewKeyboardButton("🔍 Консультация"),
				),
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("⬅️ Назад"),
				),
			)

			msg := tgbotapi.NewMessage(
				chatID,
				"📅 Запись на приём\n\nВыберите услугу:",
			)

			msg.ReplyMarkup = keyboard
			bot.Send(msg)

		case "🦷 Лечение зубов",
			"🧼 Профессиональная чистка",
			"🦷 Удаление зуба",
			"🔍 Консультация":

			user.Booking.Service = text
			user.Step = "doctor"

			keyboard := tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("👨‍⚕️ Иванов И.И."),
					tgbotapi.NewKeyboardButton("👩‍⚕️ Петрова А.А."),
				),
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("👨‍⚕️ Сидоров Д.С."),
				),
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("⬅️ Назад"),
				),
			)

			msg := tgbotapi.NewMessage(
				chatID,
				"✅ Услуга: "+text+"\n\n"+
					"Теперь выберите специалиста:",
			)

			msg.ReplyMarkup = keyboard
			bot.Send(msg)

		case "👨‍⚕️ Иванов И.И.",
			"👩‍⚕️ Петрова А.А.",
			"👨‍⚕️ Сидоров Д.С.":

			user.Booking.Doctor = text
			user.Step = "date"

			keyboard := tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("📅 6 октября"),
					tgbotapi.NewKeyboardButton("📅 7 октября"),
				),
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("📅 8 октября"),
					tgbotapi.NewKeyboardButton("📅 9 октября"),
				),
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("⬅️ Назад"),
				),
			)

			msg := tgbotapi.NewMessage(
				chatID,
				"👨‍⚕️ Специалист: "+text+"\n\n"+
					"Выберите удобную дату:",
			)

			msg.ReplyMarkup = keyboard
			bot.Send(msg)

		case "📅 6 октября",
			"📅 7 октября",
			"📅 8 октября",
			"📅 9 октября":

			user.Booking.Date = text
			user.Step = "time"

			keyboard := tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("10:00"),
					tgbotapi.NewKeyboardButton("11:00"),
					tgbotapi.NewKeyboardButton("12:00"),
				),
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("14:00"),
					tgbotapi.NewKeyboardButton("15:00"),
					tgbotapi.NewKeyboardButton("16:00"),
				),
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("⬅️ Назад"),
				),
			)

			msg := tgbotapi.NewMessage(
				chatID,
				"📅 Дата: "+text+"\n\n"+
					"Выберите свободное время:",
			)

			msg.ReplyMarkup = keyboard
			bot.Send(msg)

		case "10:00",
			"11:00",
			"12:00",
			"14:00",
			"15:00",
			"16:00":

			user.Booking.Time = text
			user.Step = "confirm"

			keyboard := tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("✅ Подтвердить запись"),
				),
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("❌ Отменить"),
				),
			)

			msg := tgbotapi.NewMessage(
				chatID,
				"📋 Проверьте данные записи:\n\n"+
					"🦷 Услуга: "+user.Booking.Service+"\n"+
					"👨‍⚕️ Врач: "+user.Booking.Doctor+"\n"+
					"📅 Дата: "+user.Booking.Date+"\n"+
					"🕐 Время: "+user.Booking.Time+"\n\n"+
					"Всё верно?",
			)

			msg.ReplyMarkup = keyboard
			bot.Send(msg)

		case "✅ Подтвердить запись":

			user.Step = "booked"

			msg := tgbotapi.NewMessage(
				chatID,
				"🎉 Запись успешно создана!\n\n"+
					"🦷 Услуга: "+user.Booking.Service+"\n"+
					"👨‍⚕️ Врач: "+user.Booking.Doctor+"\n"+
					"📅 Дата: "+user.Booking.Date+"\n"+
					"🕐 Время: "+user.Booking.Time+"\n\n"+
					"Ждём вас в клинике LZ!",
			)

			msg.ReplyMarkup = mainKeyboard()
			bot.Send(msg)

		case "📋 Моя запись":

			if user.Booking.Service == "" || user.Step != "booked" {
				msg := tgbotapi.NewMessage(
					chatID,
					"📋 У вас пока нет активных записей.",
				)

				bot.Send(msg)
				continue
			}

			msg := tgbotapi.NewMessage(
				chatID,
				"📋 Ваша запись:\n\n"+
					"🦷 Услуга: "+user.Booking.Service+"\n"+
					"👨‍⚕️ Врач: "+user.Booking.Doctor+"\n"+
					"📅 Дата: "+user.Booking.Date+"\n"+
					"🕐 Время: "+user.Booking.Time,
			)

			bot.Send(msg)

		case "ℹ️ Информация":

			msg := tgbotapi.NewMessage(
				chatID,
				"ℹ️ Клиника LZ\n\n"+
					"Стоматологическая клиника в Смоленске.\n\n"+
					"Здесь позже разместим:\n"+
					"• адрес\n"+
					"• телефон\n"+
					"• режим работы\n"+
					"• список услуг\n"+
					"• ссылку на сайт",
			)

			bot.Send(msg)

		case "⬅️ Назад":

			user.Step = ""

			msg := tgbotapi.NewMessage(
				chatID,
				"Вы вернулись в главное меню.",
			)

			msg.ReplyMarkup = mainKeyboard()
			bot.Send(msg)

		case "❌ Отменить":

			user.Step = ""
			user.Booking = Booking{}

			msg := tgbotapi.NewMessage(
				chatID,
				"❌ Запись отменена.\n\nВыберите действие:",
			)

			msg.ReplyMarkup = mainKeyboard()
			bot.Send(msg)

		default:

			msg := tgbotapi.NewMessage(
				chatID,
				"Пожалуйста, используйте кнопки меню.",
			)

			bot.Send(msg)
		}
	}
}

func mainKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📅 Записаться на приём"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📋 Моя запись"),
			tgbotapi.NewKeyboardButton("ℹ️ Информация"),
		),
	)
}