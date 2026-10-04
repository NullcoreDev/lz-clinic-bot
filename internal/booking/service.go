package booking

import (
	"fmt"

	"tg-bot/internal/catalog"
	"tg-bot/internal/storage"
)

type BookingService struct {
	storage *storage.MemoryStorage
}

func NewBookingService(storage *storage.MemoryStorage) *BookingService {
	return &BookingService{
		storage: storage,
	}
}

func (b *BookingService) GetServices() []catalog.Service {
	return catalog.GetServices()
}

func (b *BookingService) GetDoctors(serviceID int64) []catalog.Doctor {
	return catalog.GetDoctorsByService(serviceID)
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

	_, serviceExists := catalog.GetServiceByID(serviceID)
	if !serviceExists {
		return storage.Appointment{}, fmt.Errorf("услуга не найдена")
	}

	_, doctorExists := catalog.GetDoctorByID(doctorID)
	if !doctorExists {
		return storage.Appointment{}, fmt.Errorf("врач не найден")
	}

	if b.storage.IsSlotBusy(doctorID, date, timeValue) {
		return storage.Appointment{}, fmt.Errorf(
			"это время уже занято",
		)
	}

	appointment := b.storage.CreateAppointment(
		telegramID,
		patientName,
		phone,
		serviceID,
		doctorID,
		date,
		timeValue,
	)

	return appointment, nil
}

func (b *BookingService) GetUserAppointments(
	telegramID int64,
) []storage.Appointment {
	return b.storage.GetByTelegramID(telegramID)
}
