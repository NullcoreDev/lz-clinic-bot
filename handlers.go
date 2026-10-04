package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	chatID := message.Chat.ID
	text := message.Text

	if _, exists := users[chatID]; !exists {
		users[chatID] = &UserState{}
	}

	user := users[chatID]

	switch text {

	case "/start":
		handleStart(bot, chatID, user)

	case "📅 Записаться на приём":
		handleStartBooking(bot, chatID, user)

	case "🦷 Лечение зубов",
		"🧼 Профессиональная чистка",
		"🦷 Удаление зуба",
		"🔍 Консультация":

		handleService(bot, chatID, user, text)

	case "👨‍⚕️ Иванов И.И.",
		"👩‍⚕️ Петрова А.А.",
		"👨‍⚕️ Сидоров Д.С.":

		handleDoctor(bot, chatID, user, text)

	case "📅 6 октября",
		"📅 7 октября",
		"📅 8 октября",
		"📅 9 октября":

		handleDate(bot, chatID, user, text)

	case "10:00",
		"11:00",
		"12:00",
		"14:00",
		"15:00",
		"16:00":

		handleTime(bot, chatID, user, text)

	case "✅ Подтвердить запись":
		handleConfirm(bot, chatID, user)

	case "📋 Моя запись":
		handleMyBooking(bot, chatID, user)

	case "ℹ️ Информация":
		handleInformation(bot, chatID)

	case "⬅️ Назад":
		handleBack(bot, chatID, user)

	case "❌ Отменить":
		handleCancel(bot, chatID, user)

	default:
		msg := tgbotapi.NewMessage(
			chatID,
			"Пожалуйста, используйте кнопки меню.",
		)

		bot.Send(msg)
	}
}

func handleStart(
	bot *tgbotapi.BotAPI,
	chatID int64,
	user *UserState,
) {
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
}

func handleStartBooking(
	bot *tgbotapi.BotAPI,
	chatID int64,
	user *UserState,
) {
	user.Step = "service"

	msg := tgbotapi.NewMessage(
		chatID,
		"📅 Запись на приём\n\nВыберите услугу:",
	)

	msg.ReplyMarkup = serviceKeyboard()
	bot.Send(msg)
}

func handleService(
	bot *tgbotapi.BotAPI,
	chatID int64,
	user *UserState,
	service string,
) {
	user.Booking.Service = service
	user.Step = "doctor"

	msg := tgbotapi.NewMessage(
		chatID,
		"✅ Услуга: "+service+"\n\n"+
			"Теперь выберите специалиста:",
	)

	msg.ReplyMarkup = doctorKeyboard()
	bot.Send(msg)
}

func handleDoctor(
	bot *tgbotapi.BotAPI,
	chatID int64,
	user *UserState,
	doctor string,
) {
	user.Booking.Doctor = doctor
	user.Step = "date"

	msg := tgbotapi.NewMessage(
		chatID,
		"👨‍⚕️ Специалист: "+doctor+"\n\n"+
			"Выберите удобную дату:",
	)

	msg.ReplyMarkup = dateKeyboard()
	bot.Send(msg)
}

func handleDate(
	bot *tgbotapi.BotAPI,
	chatID int64,
	user *UserState,
	date string,
) {
	user.Booking.Date = date
	user.Step = "time"

	msg := tgbotapi.NewMessage(
		chatID,
		"📅 Дата: "+date+"\n\n"+
			"Выберите свободное время:",
	)

	msg.ReplyMarkup = timeKeyboard()
	bot.Send(msg)
}

func handleTime(
	bot *tgbotapi.BotAPI,
	chatID int64,
	user *UserState,
	time string,
) {
	user.Booking.Time = time
	user.Step = "confirm"

	msg := tgbotapi.NewMessage(
		chatID,
		"📋 Проверьте данные записи:\n\n"+
			"🦷 Услуга: "+user.Booking.Service+"\n"+
			"👨‍⚕️ Врач: "+user.Booking.Doctor+"\n"+
			"📅 Дата: "+user.Booking.Date+"\n"+
			"🕐 Время: "+user.Booking.Time+"\n\n"+
			"Всё верно?",
	)

	msg.ReplyMarkup = confirmKeyboard()
	bot.Send(msg)
}

func handleConfirm(
	bot *tgbotapi.BotAPI,
	chatID int64,
	user *UserState,
) {
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
}

func handleMyBooking(
	bot *tgbotapi.BotAPI,
	chatID int64,
	user *UserState,
) {
	if user.Booking.Service == "" || user.Step != "booked" {
		msg := tgbotapi.NewMessage(
			chatID,
			"📋 У вас пока нет активных записей.",
		)

		bot.Send(msg)
		return
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
}

func handleInformation(
	bot *tgbotapi.BotAPI,
	chatID int64,
) {
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
}

func handleBack(
	bot *tgbotapi.BotAPI,
	chatID int64,
	user *UserState,
) {
	user.Step = ""

	msg := tgbotapi.NewMessage(
		chatID,
		"Вы вернулись в главное меню.",
	)

	msg.ReplyMarkup = mainKeyboard()
	bot.Send(msg)
}

func handleCancel(
	bot *tgbotapi.BotAPI,
	chatID int64,
	user *UserState,
) {
	user.Step = ""
	user.Booking = Booking{}

	msg := tgbotapi.NewMessage(
		chatID,
		"❌ Запись отменена.\n\nВыберите действие:",
	)

	msg.ReplyMarkup = mainKeyboard()
	bot.Send(msg)
}