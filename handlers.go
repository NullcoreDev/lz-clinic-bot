package main

import (
	"fmt"
	"strings"
	"time"

	"tg-bot/internal/booking"
	"tg-bot/internal/catalog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

	case "✅ Подтвердить запись":
		handleConfirm(bot, chatID, user, bookingService)

	default:
		switch user.Step {
		case "date":
			handleDate(bot, chatID, user, text, bookingService)
			return

		case "time":
			handleTime(bot, chatID, user, text, bookingService)
			return

		case "patient_name":
			handlePatientName(bot, chatID, user, text)
			return

		case "phone":
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
	resetUser(user)

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
	resetUser(user)

	user.Step = "service"

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
	user.Booking.Date = ""
	user.Booking.Time = ""
	user.Step = "doctor"

	doctors := bookingService.GetDoctors(serviceID)

	if len(doctors) == 0 {
		sendError(bot, chatID, "Для этой услуги пока нет доступных врачей.")
		return
	}

	var text strings.Builder

	text.WriteString("Вы выбрали:\n")
	text.WriteString("🩺 " + service.Name + "\n\n")
	text.WriteString("Теперь выберите врача:")

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

	if user.Booking.ServiceID == 0 {
		sendError(bot, chatID, "Сначала выберите направление.")
		return
	}

	doctors := bookingService.GetDoctors(user.Booking.ServiceID)

	doctorAvailable := false

	for _, item := range doctors {
		if item.ID == doctorID {
			doctorAvailable = true
			break
		}
	}

	if !doctorAvailable {
		sendError(bot, chatID, "Этот врач недоступен для выбранного направления.")
		return
	}

	user.Booking.DoctorID = doctorID
	user.Booking.Date = ""
	user.Booking.Time = ""
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
	bookingService *booking.BookingService,
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

	now := time.Now()
	today := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0,
		0,
		0,
		0,
		now.Location(),
	)

	selectedDate := time.Date(
		parsedDate.Year(),
		parsedDate.Month(),
		parsedDate.Day(),
		0,
		0,
		0,
		0,
		now.Location(),
	)

	if selectedDate.Before(today) {
		msg := tgbotapi.NewMessage(
			chatID,
			"Эта дата уже прошла. Выберите другую.",
		)

		msg.ReplyMarkup = dateKeyboard()
		bot.Send(msg)
		return
	}

	user.Booking.Date = selectedDate.Format("02.01.2006")

	availableTimes, err := bookingService.GetAvailableTimes(
		user.Booking.DoctorID,
		user.Booking.Date,
	)

	if err != nil {
		sendError(bot, chatID, "Не удалось получить доступное время.")
		return
	}

	if len(availableTimes) == 0 {
		msg := tgbotapi.NewMessage(
			chatID,
			"На выбранную дату свободных записей у этого врача нет.\n\n"+
				"Выберите другую дату:",
		)

		msg.ReplyMarkup = dateKeyboard()
		bot.Send(msg)

		user.Booking.Time = ""
		user.Step = "date"

		return
	}

	user.Step = "time"

	msg := tgbotapi.NewMessage(
		chatID,
		"Дата: "+user.Booking.Date+
			"\n\nВыберите свободное время:",
	)

	msg.ReplyMarkup = timeKeyboard(availableTimes)
	bot.Send(msg)
}

func handleTime(
	bot *tgbotapi.BotAPI,
	chatID int64,
	user *UserState,
	timeValue string,
	bookingService *booking.BookingService,
) {
	timeValue = strings.TrimSpace(timeValue)

	_, err := time.Parse("15:04", timeValue)

	if err != nil {
		sendError(bot, chatID, "Пожалуйста, выберите время кнопкой.")
		return
	}

	availableTimes, err := bookingService.GetAvailableTimes(
		user.Booking.DoctorID,
		user.Booking.Date,
	)

	if err != nil {
		sendError(bot, chatID, "Не удалось проверить доступность времени.")
		return
	}

	isAvailable := false

	for _, availableTime := range availableTimes {
		if availableTime == timeValue {
			isAvailable = true
			break
		}
	}

	if !isAvailable {
		msg := tgbotapi.NewMessage(
			chatID,
			"Это время уже занято или недоступно.\n\n"+
				"Выберите другое время:",
		)

		msg.ReplyMarkup = timeKeyboard(availableTimes)
		bot.Send(msg)
		return
	}

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
	name = strings.TrimSpace(name)

	if len([]rune(name)) < 2 {
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
	phone = strings.TrimSpace(phone)

	if len([]rune(phone)) < 5 {
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

	resetUser(user)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = mainKeyboard()

	bot.Send(msg)
}

func handleMyBookings(
	bot *tgbotapi.BotAPI,
	chatID int64,
	bookingService *booking.BookingService,
) {
	appointments, err := bookingService.GetUserAppointments(chatID)

	if err != nil {
		sendError(bot, chatID, "Не удалось получить ваши записи.")
		return
	}

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
		user.Booking.Time = ""

		msg := tgbotapi.NewMessage(
			chatID,
			"Выберите дату:",
		)

		msg.ReplyMarkup = dateKeyboard()
		bot.Send(msg)

	case "patient_name":
		availableTimes, err := bookingService.GetAvailableTimes(
			user.Booking.DoctorID,
			user.Booking.Date,
		)

		if err != nil {
			user.Step = "date"

			msg := tgbotapi.NewMessage(
				chatID,
				"Не удалось получить время.\n\nВыберите дату:",
			)

			msg.ReplyMarkup = dateKeyboard()
			bot.Send(msg)
			return
		}

		user.Step = "time"
		user.Booking.Time = ""

		msg := tgbotapi.NewMessage(
			chatID,
			"Выберите свободное время:",
		)

		msg.ReplyMarkup = timeKeyboard(availableTimes)
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
	resetUser(user)

	msg := tgbotapi.NewMessage(
		chatID,
		"❌ Запись отменена.\n\nВы вернулись в главное меню.",
	)

	msg.ReplyMarkup = mainKeyboard()
	bot.Send(msg)
}

func resetUser(user *UserState) {
	user.Step = ""
	user.Booking = Booking{}
	user.PatientName = ""
	user.Phone = ""
}

func sendError(
	bot *tgbotapi.BotAPI,
	chatID int64,
	text string,
) {
	msg := tgbotapi.NewMessage(chatID, "❌ "+text)
	bot.Send(msg)
}
