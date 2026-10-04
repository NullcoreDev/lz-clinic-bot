package main

import (
	"fmt"
	"time"

	"tg-bot/internal/catalog"

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
			tgbotapi.NewKeyboardButton("❌ Отменить"),
		),
	)

	keyboard.ResizeKeyboard = true

	return keyboard
}

func doctorKeyboard(serviceID int64) tgbotapi.ReplyKeyboardMarkup {
	doctors := catalog.GetDoctorsByService(serviceID)

	var rows [][]tgbotapi.KeyboardButton

	for _, doctor := range doctors {
		rows = append(
			rows,
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(doctor.Name),
			),
		)
	}

	rows = append(
		rows,
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("⬅️ Назад"),
			tgbotapi.NewKeyboardButton("❌ Отменить"),
		),
	)

	keyboard := tgbotapi.NewReplyKeyboard(rows...)
	keyboard.ResizeKeyboard = true

	return keyboard
}

func dateKeyboard() tgbotapi.ReplyKeyboardMarkup {
	var rows [][]tgbotapi.KeyboardButton

	now := time.Now()

	for i := 1; i <= 7; i++ {
		date := now.AddDate(0, 0, i)

		button := fmt.Sprintf(
			"📅 %s",
			date.Format("02.01.2006"),
		)

		rows = append(
			rows,
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(button),
			),
		)
	}

	rows = append(
		rows,
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("⬅️ Назад"),
			tgbotapi.NewKeyboardButton("❌ Отменить"),
		),
	)

	keyboard := tgbotapi.NewReplyKeyboard(rows...)
	keyboard.ResizeKeyboard = true

	return keyboard
}

func timeKeyboard(times []string) tgbotapi.ReplyKeyboardMarkup {
	var rows [][]tgbotapi.KeyboardButton

	var currentRow []tgbotapi.KeyboardButton

	for _, timeValue := range times {
		currentRow = append(
			currentRow,
			tgbotapi.NewKeyboardButton(timeValue),
		)

		if len(currentRow) == 2 {
			rows = append(rows, currentRow)
			currentRow = nil
		}
	}

	if len(currentRow) > 0 {
		rows = append(rows, currentRow)
	}

	rows = append(
		rows,
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("⬅️ Назад"),
			tgbotapi.NewKeyboardButton("❌ Отменить"),
		),
	)

	keyboard := tgbotapi.NewReplyKeyboard(rows...)
	keyboard.ResizeKeyboard = true

	return keyboard
}

func confirmKeyboard() tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("✅ Подтвердить запись"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("⬅️ Назад"),
			tgbotapi.NewKeyboardButton("❌ Отменить"),
		),
	)

	keyboard.ResizeKeyboard = true

	return keyboard
}
