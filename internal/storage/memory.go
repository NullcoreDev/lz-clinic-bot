package storage

import (
	"sync"
	"time"
)

type Appointment struct {
	ID          int64
	TelegramID  int64
	PatientName string
	Phone       string
	ServiceID   int64
	DoctorID    int64
	Date        string
	Time        string
	CreatedAt   time.Time
}

type MemoryStorage struct {
	mu           sync.RWMutex
	appointments []Appointment
	nextID       int64
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		appointments: make([]Appointment, 0),
		nextID:       1,
	}
}

func (s *MemoryStorage) CreateAppointment(
	telegramID int64,
	patientName string,
	phone string,
	serviceID int64,
	doctorID int64,
	date string,
	timeValue string,
) Appointment {

	s.mu.Lock()
	defer s.mu.Unlock()

	appointment := Appointment{
		ID:          s.nextID,
		TelegramID:  telegramID,
		PatientName: patientName,
		Phone:       phone,
		ServiceID:   serviceID,
		DoctorID:    doctorID,
		Date:        date,
		Time:        timeValue,
		CreatedAt:   time.Now(),
	}

	s.appointments = append(s.appointments, appointment)
	s.nextID++

	return appointment
}

func (s *MemoryStorage) GetByTelegramID(telegramID int64) []Appointment {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []Appointment

	for _, appointment := range s.appointments {
		if appointment.TelegramID == telegramID {
			result = append(result, appointment)
		}
	}

	return result
}

func (s *MemoryStorage) IsSlotBusy(
	doctorID int64,
	date string,
	timeValue string,
) bool {

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, appointment := range s.appointments {
		if appointment.DoctorID == doctorID &&
			appointment.Date == date &&
			appointment.Time == timeValue {
			return true
		}
	}

	return false
}
