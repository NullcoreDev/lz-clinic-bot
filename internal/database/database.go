package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

func CreateDatabaseIfNotExists() error {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/postgres",
		user,
		password,
		host,
		port,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return fmt.Errorf("ошибка подключения к PostgreSQL: %w", err)
	}
	defer conn.Close(ctx)

	dbName := os.Getenv("DB_NAME")

	var exists bool

	err = conn.QueryRow(
		ctx,
		"SELECT EXISTS (SELECT FROM pg_database WHERE datname = $1)",
		dbName,
	).Scan(&exists)

	if err != nil {
		return fmt.Errorf("ошибка проверки базы: %w", err)
	}

	if exists {
		return nil
	}

	_, err = conn.Exec(
		ctx,
		fmt.Sprintf(`CREATE DATABASE "%s"`, dbName),
	)

	if err != nil {
		return fmt.Errorf("ошибка создания базы: %w", err)
	}

	return nil
}

func Connect() (*pgx.Conn, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		user,
		password,
		host,
		port,
		dbName,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе: %w", err)
	}

	if err := conn.Ping(ctx); err != nil {
		conn.Close(ctx)
		return nil, fmt.Errorf("база не отвечает: %w", err)
	}

	return conn, nil
}