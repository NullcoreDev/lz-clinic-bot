package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"tg-bot/internal/booking"
	"tg-bot/internal/catalog"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleMessage(
	bot *tgbotapi.BotAPI,
	message *tgbotapi.Message,
	bookingService *booking.BookingService,
) {
	chatID := message.Chat.ID
	text := strings.TrimSpace(message.Text)

	user, exists := users[chatID]

	if !exists {
		user = &UserState{}
		users[chatID] = user
	}

	switch text {

	case "/start":
		handleStart(bot, chatID, user)

	case "📅 Записаться на приём":
		handleStartBooking(bot, chatID, user)

	case "📋 Моя запись":
		handleMyBookings(bot, chatID, bookingService)

	case "ℹ️ Информация":
		handleInformation(bot, chatID)

	case "⬅️ Назад":
		handleBack(bot, chatID, user, bookingService)

	case "❌ Отменить":
		handleCancel(bot, chatID, user)

	case "🧠 Неврология":
		handleService(bot, chatID, user, 1, bookingService)

	case "❤️ Кардиология":
		handleService(bot, chatID, user, 2, bookingService)

	case "🩺 Терапия":
		handleService(bot, chatID, user, 3, bookingService)

	case "Раков Александр Михайлович":
		handleDoctor(bot, chatID, user, 1, bookingService)

	case "Гуревич Оксана Васильевна":
		handleDoctor(bot, chatID, user, 2, bookingService)

	case "Игнатенкова Эльвира Ильгизовна":
		handleDoctor(bot, chatID, user, 3, bookingService)

	case "10:00", "11:00", "12:00", "14:00", "15:00", "16:00":
		handleTime(bot, chatID, user, text)

	case "✅ Подтвердить запись":
		handleConfirm(bot, chatID, user, bookingService)

	default:
		if user.Step == "date" {
			handleDate(bot, chatID, user, text)
			return
		}

		if user.Step == "patient_name" {
			handlePatientName(bot, chatID, user, text)
			return
		}

		if user.Step == "phone" {
			handlePhone(bot, chatID, user, text)
			return
		}

		msg := tgbotapi.NewMessage(
			chatID,
			"Не понял команду. Используйте кнопки меню.",
		)

		msg.ReplyMarkup = mainKeyboard()
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
	user.PatientName = ""
	user.Phone = ""

	msg := tgbotapi.NewMessage(
		chatID,
		"Здравствуйте! 👋\n\n"+
			"Это бот записи в ЛЗ Клиник.\n\n"+
			"Выберите действие:",
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
	user.Booking = Booking{}
	user.PatientName = ""
	user.Phone = ""

	msg := tgbotapi.NewMessage(
		chatID,
		"Выберите направление:",
	)

	msg.ReplyMarkup = serviceKeyboard()

	bot.Send(msg)
}

func handleService(
	bot *tgbotapi.BotAPI,
	chatID int64,
	user *UserState,
	serviceID int64,
	bookingService *booking.BookingService,
) {
	service, exists := catalog.GetServiceByID(serviceID)

	if !exists {
		sendError(bot, chatID, "Услуга не найдена.")
		return
	}

	user.Booking.ServiceID = serviceID
	user.Booking.DoctorID = 0
	user.Step = "doctor"

	doctors := bookingService.GetDoctors(serviceID)

	var text strings.Builder

	text.WriteString("Вы выбрали:\n")
	text.WriteString("🩺 " + service.Name + "\n\n")
	text.WriteString("Теперь выберите врача:")

	if len(doctors) == 0 {
		sendError(bot, chatID, "Для этой услуги пока нет доступных врачей.")
		return
	}

	msg := tgbotapi.NewMessage(chatID, text.String())
	msg.ReplyMarkup = doctorKeyboard(serviceID)

	bot.Send(msg)
}

func handleDoctor(
	bot *tgbotapi.BotAPI,
	chatID int64,
	user *UserState,
	doctorID int64,
	bookingService *booking.BookingService,
) {
	doctor, exists := catalog.GetDoctorByID(doctorID)

	if !exists {
		sendError(bot, chatID, "Врач не найден.")
		return
	}

	user.Booking.DoctorID = doctorID
	user.Step = "date"

	msg := tgbotapi.NewMessage(
		chatID,
		"Вы выбрали врача:\n\n"+
			"👨‍⚕️ "+doctor.Name+"\n"+
			"Специализация: "+doctor.Specialization+
			"\n\nВыберите дату:",
	)

	msg.ReplyMarkup = dateKeyboard()

	bot.Send(msg)
}

func handleDate(
	bot *tgbotapi.BotAPI,
	chatID int64,
	user *UserState,
	text string,
) {
	text = strings.TrimPrefix(text, "📅 ")
	text = strings.TrimSpace(text)

	parsedDate, err := time.Parse("02.01.2006", text)

	if err != nil {
		msg := tgbotapi.NewMessage(
			chatID,
			"Пожалуйста, выберите дату кнопкой.",
		)

		msg.ReplyMarkup = dateKeyboard()
		bot.Send(msg)
		return
	}

	if parsedDate.Before(time.Now().Truncate(24 * time.Hour)) {
		msg := tgbotapi.NewMessage(
			chatID,
			"Эта дата уже прошла. Выберите другую.",
		)

		msg.ReplyMarkup = dateKeyboard()
		bot.Send(msg)
		return
	}

	user.Booking.Date = parsedDate.Format("02.01.2006")
	user.Step = "time"

	msg := tgbotapi.NewMessage(
		chatID,
		"Дата: "+user.Booking.Date+
			"\n\nВыберите время:",
	)

	msg.ReplyMarkup = timeKeyboard()

	bot.Send(msg)
}

func handleTime(
	bot *tgbotapi.BotAPI,
	chatID int64,
	user *UserState,
	timeValue string,
) {
	user.Booking.Time = timeValue
	user.Step = "patient_name"

	msg := tgbotapi.NewMessage(
		chatID,
		"Отлично 👍\n\n"+
			"Теперь введите ваше имя и фамилию:",
	)

	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

	bot.Send(msg)
}

func handlePatientName(
	bot *tgbotapi.BotAPI,
	chatID int64,
	user *UserState,
	name string,
) {
	if len(name) < 2 {
		bot.Send(tgbotapi.NewMessage(
			chatID,
			"Пожалуйста, введите имя и фамилию.",
		))
		return
	}

	user.PatientName = name
	user.Step = "phone"

	msg := tgbotapi.NewMessage(
		chatID,
		"Введите номер телефона для связи:",
	)

	bot.Send(msg)
}

func handlePhone(
	bot *tgbotapi.BotAPI,
	chatID int64,
	user *UserState,
	phone string,
) {
	if len(phone) < 5 {
		bot.Send(tgbotapi.NewMessage(
			chatID,
			"Похоже, номер слишком короткий. Введите номер телефона ещё раз.",
		))
		return
	}

	user.Phone = phone
	user.Step = "confirm"

	service, _ := catalog.GetServiceByID(user.Booking.ServiceID)
	doctor, _ := catalog.GetDoctorByID(user.Booking.DoctorID)

	text := fmt.Sprintf(
		"Проверьте данные записи:\n\n"+
			"🩺 Услуга: %s\n"+
			"👨‍⚕️ Врач: %s\n"+
			"📅 Дата: %s\n"+
			"🕐 Время: %s\n"+
			"👤 Пациент: %s\n"+
			"📞 Телефон: %s\n\n"+
			"Всё верно?",
		service.Name,
		doctor.Name,
		user.Booking.Date,
		user.Booking.Time,
		user.PatientName,
		user.Phone,
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = confirmKeyboard()

	bot.Send(msg)
}

func handleConfirm(
	bot *tgbotapi.BotAPI,
	chatID int64,
	user *UserState,
	bookingService *booking.BookingService,
) {
	appointment, err := bookingService.CreateAppointment(
		chatID,
		user.PatientName,
		user.Phone,
		user.Booking.ServiceID,
		user.Booking.DoctorID,
		user.Booking.Date,
		user.Booking.Time,
	)

	if err != nil {
		msg := tgbotapi.NewMessage(
			chatID,
			"❌ Не удалось создать запись.\n\n"+err.Error(),
		)

		msg.ReplyMarkup = mainKeyboard()
		bot.Send(msg)

		return
	}

	service, _ := catalog.GetServiceByID(appointment.ServiceID)
	doctor, _ := catalog.GetDoctorByID(appointment.DoctorID)

	text := fmt.Sprintf(
		"✅ Запись успешно создана!\n\n"+
			"Номер записи: #%d\n\n"+
			"🩺 %s\n"+
			"👨‍⚕️ %s\n"+
			"📅 %s\n"+
			"🕐 %s\n\n"+
			"Пациент: %s\n"+
			"Телефон: %s\n\n"+
			"Адрес клиники:\n"+
			"г. Смоленск, ул. Нарвская, д. 4",
		appointment.ID,
		service.Name,
		doctor.Name,
		appointment.Date,
		appointment.Time,
		appointment.PatientName,
		appointment.Phone,
	)

	user.Step = ""
	user.Booking = Booking{}
	user.PatientName = ""
	user.Phone = ""

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = mainKeyboard()

	bot.Send(msg)
}

func handleMyBookings(
	bot *tgbotapi.BotAPI,
	chatID int64,
	bookingService *booking.BookingService,
) {
	appointments := bookingService.GetUserAppointments(chatID)

	if len(appointments) == 0 {
		msg := tgbotapi.NewMessage(
			chatID,
			"📋 У вас пока нет записей.",
		)

		msg.ReplyMarkup = mainKeyboard()
		bot.Send(msg)
		return
	}

	var text strings.Builder

	text.WriteString("📋 Ваши записи:\n\n")

	for _, appointment := range appointments {
		service, _ := catalog.GetServiceByID(appointment.ServiceID)
		doctor, _ := catalog.GetDoctorByID(appointment.DoctorID)

		text.WriteString(
			fmt.Sprintf(
				"Запись #%d\n"+
					"🩺 %s\n"+
					"👨‍⚕️ %s\n"+
					"📅 %s\n"+
					"🕐 %s\n\n",
				appointment.ID,
				service.Name,
				doctor.Name,
				appointment.Date,
				appointment.Time,
			),
		)
	}

	msg := tgbotapi.NewMessage(chatID, text.String())
	msg.ReplyMarkup = mainKeyboard()

	bot.Send(msg)
}

func handleInformation(
	bot *tgbotapi.BotAPI,
	chatID int64,
) {
	text :=
		"ℹ️ ЛЗ Клиник\n\n" +
			"Многопрофильная клиника для взрослых и детей.\n\n" +
			"📍 г. Смоленск, ул. Нарвская, д. 4\n" +
			"📞 +7 (4812) 51-03-03\n\n" +
			"🕐 Пн-Пт: 08:00-20:00\n" +
			"🕐 Сб: 08:00-20:00\n" +
			"🕐 Вс: 09:00-19:00"

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = mainKeyboard()

	bot.Send(msg)
}

func handleBack(
	bot *tgbotapi.BotAPI,
	chatID int64,
	user *UserState,
	bookingService *booking.BookingService,
) {
	switch user.Step {

	case "service":
		handleStart(bot, chatID, user)

	case "doctor":
		user.Step = "service"

		msg := tgbotapi.NewMessage(
			chatID,
			"Выберите направление:",
		)

		msg.ReplyMarkup = serviceKeyboard()
		bot.Send(msg)

	case "date":
		user.Step = "doctor"

		msg := tgbotapi.NewMessage(
			chatID,
			"Выберите врача:",
		)

		msg.ReplyMarkup = doctorKeyboard(user.Booking.ServiceID)
		bot.Send(msg)

	case "time":
		user.Step = "date"

		msg := tgbotapi.NewMessage(
			chatID,
			"Выберите дату:",
		)

		msg.ReplyMarkup = dateKeyboard()
		bot.Send(msg)

	case "patient_name":
		user.Step = "time"

		msg := tgbotapi.NewMessage(
			chatID,
			"Выберите время:",
		)

		msg.ReplyMarkup = timeKeyboard()
		bot.Send(msg)

	case "phone":
		user.Step = "patient_name"

		msg := tgbotapi.NewMessage(
			chatID,
			"Введите ваше имя и фамилию:",
		)

		bot.Send(msg)

	case "confirm":
		user.Step = "phone"

		msg := tgbotapi.NewMessage(
			chatID,
			"Введите номер телефона:",
		)

		bot.Send(msg)

	default:
		handleStart(bot, chatID, user)
	}
}

func handleCancel(
	bot *tgbotapi.BotAPI,
	chatID int64,
	user *UserState,
) {
	user.Step = ""
	user.Booking = Booking{}
	user.PatientName = ""
	user.Phone = ""

	msg := tgbotapi.NewMessage(
		chatID,
		"❌ Запись отменена.\n\nВыберите действие:",
	)

	msg.ReplyMarkup = mainKeyboard()

	bot.Send(msg)
}

func sendError(
	bot *tgbotapi.BotAPI,
	chatID int64,
	text string,
) {
	msg := tgbotapi.NewMessage(chatID, "❌ "+text)
	msg.ReplyMarkup = mainKeyboard()
	bot.Send(msg)
}

func parseID(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}
