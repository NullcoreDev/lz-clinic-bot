package booking

import (
	"fmt"

	"tg-bot/internal/catalog"
	"tg-bot/internal/storage"
)

type BookingStorage interface {
	CreateAppointment(
		telegramID int64,
		patientName string,
		phone string,
		serviceID int64,
		doctorID int64,
		date string,
		timeValue string,
	) (storage.Appointment, error)

	GetByTelegramID(
		telegramID int64,
	) ([]storage.Appointment, error)

	IsSlotBusy(
		doctorID int64,
		date string,
		timeValue string,
	) (bool, error)

	GetBusyTimes(
		doctorID int64,
		date string,
	) ([]string, error)
}

type BookingService struct {
	storage BookingStorage
}

func NewBookingService(storage BookingStorage) *BookingService {
	return &BookingService{
		storage: storage,
	}
}

func (b *BookingService) GetServices() []catalog.Service {
	return catalog.GetServices()
}

func (b *BookingService) GetDoctors(
	serviceID int64,
) []catalog.Doctor {
	return catalog.GetDoctorsByService(serviceID)
}

func (b *BookingService) GetAvailableTimes(
	doctorID int64,
	date string,
) ([]string, error) {
	// Это временное расписание прототипа.
	// Позже здесь можно подключить реальное расписание клиники.
	allTimes := []string{
		"10:00",
		"11:00",
		"12:00",
		"13:00",
		"14:00",
		"15:00",
		"16:00",
		"17:00",
	}

	busyTimes, err := b.storage.GetBusyTimes(
		doctorID,
		date,
	)

	if err != nil {
		return nil, err
	}

	busy := make(map[string]bool)

	for _, timeValue := range busyTimes {
		busy[timeValue] = true
	}

	var available []string

	for _, timeValue := range allTimes {
		if !busy[timeValue] {
			available = append(available, timeValue)
		}
	}

	return available, nil
}

func (b *BookingService) CreateAppointment(
	telegramID int64,
	patientName string,
	phone string,
	serviceID int64,
	doctorID int64,
	date string,
	timeValue string,
) (storage.Appointment, error) {
	service, serviceExists := catalog.GetServiceByID(serviceID)

	if !serviceExists {
		return storage.Appointment{}, fmt.Errorf(
			"услуга не найдена",
		)
	}

	doctor, doctorExists := catalog.GetDoctorByID(doctorID)

	if !doctorExists {
		return storage.Appointment{}, fmt.Errorf(
			"врач не найден",
		)
	}

	doctorCanProvideService := false

	for _, id := range doctor.ServiceIDs {
		if id == service.ID {
			doctorCanProvideService = true
			break
		}
	}

	if !doctorCanProvideService {
		return storage.Appointment{}, fmt.Errorf(
			"выбранный врач не оказывает эту услугу",
		)
	}

	availableTimes, err := b.GetAvailableTimes(
		doctorID,
		date,
	)

	if err != nil {
		return storage.Appointment{}, err
	}

	timeIsAvailable := false

	for _, availableTime := range availableTimes {
		if availableTime == timeValue {
			timeIsAvailable = true
			break
		}
	}

	if !timeIsAvailable {
		return storage.Appointment{}, fmt.Errorf(
			"это время уже занято или недоступно",
		)
	}

	busy, err := b.storage.IsSlotBusy(
		doctorID,
		date,
		timeValue,
	)

	if err != nil {
		return storage.Appointment{}, err
	}

	if busy {
		return storage.Appointment{}, fmt.Errorf(
			"это время уже занято",
		)
	}

	return b.storage.CreateAppointment(
		telegramID,
		patientName,
		phone,
		serviceID,
		doctorID,
		date,
		timeValue,
	)
}

func (b *BookingService) GetUserAppointments(
	telegramID int64,
) ([]storage.Appointment, error) {
	return b.storage.GetByTelegramID(telegramID)
}
