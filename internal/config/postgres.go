package config

import (
	"context"
	"testing"

	"github.com/whitewolf185/mangaparser/internal/config/flags"

	"github.com/jackc/pgx/v5"
)

// ConnectPostgres подключение к базе данных. Возвращает указатель на коннект к базе и ошибку.
// Перед использованием обязательно нужно сделать парсинг env переменных через вызов функции flags.InitServiceFlags
func ConnectPostgres(ctx context.Context) (*pgx.Conn, error) {
	return pgx.Connect(ctx, GetValue(DbDsn))
}

func TestPostgres(t *testing.T) *pgx.Conn {
	flags.InitServiceFlags()
	db, err := pgx.Connect(context.Background(), GetValue(DbDsn))
	if err != nil {
		t.Fatal(err)
	}
	return db
}