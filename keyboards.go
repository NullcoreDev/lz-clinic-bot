package main

import (
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func mainKeyboard() tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📅 Записаться на приём"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📋 Моя запись"),
			tgbotapi.NewKeyboardButton("ℹ️ Информация"),
		),
	)

	keyboard.ResizeKeyboard = true

	return keyboard
}

func serviceKeyboard() tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🧠 Неврология"),
			tgbotapi.NewKeyboardButton("❤️ Кардиология"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🩺 Терапия"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("⬅️ Назад"),
		),
	)

	keyboard.ResizeKeyboard = true

	return keyboard
}

func doctorKeyboard(serviceID int64) tgbotapi.ReplyKeyboardMarkup {
	switch serviceID {
	case 1:
		return tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Раков Александр Михайлович"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("⬅️ Назад"),
			),
		)

	case 2:
		return tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Гуревич Оксана Васильевна"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("⬅️ Назад"),
			),
		)

	case 3:
		return tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Игнатенкова Эльвира Ильгизовна"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("⬅️ Назад"),
			),
		)
	}

	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("⬅️ Назад"),
		),
	)
}

func dateKeyboard() tgbotapi.ReplyKeyboardMarkup {
	now := time.Now()

	rows := make([][]tgbotapi.KeyboardButton, 0, 5)

	for i := 1; i <= 4; i++ {
		date := now.AddDate(0, 0, i)

		button := tgbotapi.NewKeyboardButton(
			fmt.Sprintf("📅 %s", date.Format("02.01.2006")),
		)

		rows = append(rows, []tgbotapi.KeyboardButton{
			button,
		})
	}

	rows = append(rows, []tgbotapi.KeyboardButton{
		tgbotapi.NewKeyboardButton("⬅️ Назад"),
	})

	keyboard := tgbotapi.ReplyKeyboardMarkup{
		Keyboard:       rows,
		ResizeKeyboard: true,
	}

	return keyboard
}

func timeKeyboard() tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("10:00"),
			tgbotapi.NewKeyboardButton("11:00"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("12:00"),
			tgbotapi.NewKeyboardButton("14:00"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("15:00"),
			tgbotapi.NewKeyboardButton("16:00"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("⬅️ Назад"),
		),
	)

	keyboard.ResizeKeyboard = true

	return keyboard
}

func confirmKeyboard() tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("✅ Подтвердить запись"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("❌ Отменить"),
		),
	)

	keyboard.ResizeKeyboard = true

	return keyboard
}
