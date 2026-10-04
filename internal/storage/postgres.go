package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type PostgresStorage struct {
	db *pgx.Conn
}

func NewPostgresStorage(db *pgx.Conn) (*PostgresStorage, error) {
	storage := &PostgresStorage{
		db: db,
	}

	if err := storage.createTables(); err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *PostgresStorage) createTables() error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	query := `
	CREATE TABLE IF NOT EXISTS appointments (
		id BIGSERIAL PRIMARY KEY,
		telegram_id BIGINT NOT NULL,
		patient_name TEXT NOT NULL,
		phone TEXT NOT NULL,
		service_id BIGINT NOT NULL,
		doctor_id BIGINT NOT NULL,
		date DATE NOT NULL,
		time TIME NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),

		UNIQUE (doctor_id, date, time)
	);
	`

	_, err := s.db.Exec(ctx, query)

	if err != nil {
		return fmt.Errorf(
			"ошибка создания таблицы appointments: %w",
			err,
		)
	}

	return nil
}

func (s *PostgresStorage) CreateAppointment(
	telegramID int64,
	patientName string,
	phone string,
	serviceID int64,
	doctorID int64,
	date string,
	timeValue string,
) (Appointment, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	parsedDate, err := time.Parse(
		"02.01.2006",
		date,
	)

	if err != nil {
		return Appointment{}, fmt.Errorf(
			"неверный формат даты: %w",
			err,
		)
	}

	parsedTime, err := time.Parse(
		"15:04",
		timeValue,
	)

	if err != nil {
		return Appointment{}, fmt.Errorf(
			"неверный формат времени: %w",
			err,
		)
	}

	query := `
		INSERT INTO appointments (
			telegram_id,
			patient_name,
			phone,
			service_id,
			doctor_id,
			date,
			time
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING
			id,
			telegram_id,
			patient_name,
			phone,
			service_id,
			doctor_id,
			TO_CHAR(date, 'DD.MM.YYYY'),
			TO_CHAR(time, 'HH24:MI'),
			created_at;
	`

	var appointment Appointment

	err = s.db.QueryRow(
		ctx,
		query,
		telegramID,
		patientName,
		phone,
		serviceID,
		doctorID,
		parsedDate,
		parsedTime,
	).Scan(
		&appointment.ID,
		&appointment.TelegramID,
		&appointment.PatientName,
		&appointment.Phone,
		&appointment.ServiceID,
		&appointment.DoctorID,
		&appointment.Date,
		&appointment.Time,
		&appointment.CreatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError

		if ok := errorAs(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return Appointment{}, fmt.Errorf(
					"это время уже занято другим пациентом",
				)
			}
		}

		return Appointment{}, fmt.Errorf(
			"ошибка создания записи: %w",
			err,
		)
	}

	return appointment, nil
}

func (s *PostgresStorage) GetByTelegramID(
	telegramID int64,
) ([]Appointment, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	query := `
		SELECT
			id,
			telegram_id,
			patient_name,
			phone,
			service_id,
			doctor_id,
			TO_CHAR(date, 'DD.MM.YYYY'),
			TO_CHAR(time, 'HH24:MI'),
			created_at
		FROM appointments
		WHERE telegram_id = $1
		ORDER BY date, time;
	`

	rows, err := s.db.Query(
		ctx,
		query,
		telegramID,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"ошибка получения записей: %w",
			err,
		)
	}

	defer rows.Close()

	var appointments []Appointment

	for rows.Next() {
		var appointment Appointment

		err := rows.Scan(
			&appointment.ID,
			&appointment.TelegramID,
			&appointment.PatientName,
			&appointment.Phone,
			&appointment.ServiceID,
			&appointment.DoctorID,
			&appointment.Date,
			&appointment.Time,
			&appointment.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf(
				"ошибка чтения записи: %w",
				err,
			)
		}

		appointments = append(
			appointments,
			appointment,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"ошибка обработки записей: %w",
			err,
		)
	}

	return appointments, nil
}

func (s *PostgresStorage) IsSlotBusy(
	doctorID int64,
	date string,
	timeValue string,
) (bool, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	parsedDate, err := time.Parse(
		"02.01.2006",
		date,
	)

	if err != nil {
		return false, fmt.Errorf(
			"неверный формат даты: %w",
			err,
		)
	}

	parsedTime, err := time.Parse(
		"15:04",
		timeValue,
	)

	if err != nil {
		return false, fmt.Errorf(
			"неверный формат времени: %w",
			err,
		)
	}

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM appointments
			WHERE doctor_id = $1
			  AND date = $2
			  AND time = $3
		);
	`

	var exists bool

	err = s.db.QueryRow(
		ctx,
		query,
		doctorID,
		parsedDate,
		parsedTime,
	).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf(
			"ошибка проверки времени: %w",
			err,
		)
	}

	return exists, nil
}

func (s *PostgresStorage) GetBusyTimes(
	doctorID int64,
	date string,
) ([]string, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	parsedDate, err := time.Parse(
		"02.01.2006",
		date,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"неверный формат даты: %w",
			err,
		)
	}

	query := `
		SELECT TO_CHAR(time, 'HH24:MI')
		FROM appointments
		WHERE doctor_id = $1
		  AND date = $2
		ORDER BY time;
	`

	rows, err := s.db.Query(
		ctx,
		query,
		doctorID,
		parsedDate,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"ошибка получения занятых времён: %w",
			err,
		)
	}

	defer rows.Close()

	var busyTimes []string

	for rows.Next() {
		var timeValue string

		if err := rows.Scan(&timeValue); err != nil {
			return nil, fmt.Errorf(
				"ошибка чтения занятого времени: %w",
				err,
			)
		}

		busyTimes = append(
			busyTimes,
			timeValue,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"ошибка обработки занятых времён: %w",
			err,
		)
	}

	return busyTimes, nil
}

func errorAs(err error, target **pgconn.PgError) bool {
	return pgxErrorAs(err, target)
}

func pgxErrorAs(err error, target **pgconn.PgError) bool {
	for err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			*target = pgErr
			return true
		}

		type unwrapper interface {
			Unwrap() error
		}

		u, ok := err.(unwrapper)

		if !ok {
			break
		}

		err = u.Unwrap()
	}

	return false
}
