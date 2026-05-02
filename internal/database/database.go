package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
}

func New(dsn string) (*DB, error) {
	if dsn == "" {
		return nil, fmt.Errorf("dsn is empty")
	}

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse dsn: %w", err)
	}

	config.MinConns = 2
	config.MaxConns = 10

	config.MaxConnLifetime = 30 * time.Minute
	config.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}

	slog.Info("database connection established")
	return &DB{pool: pool}, nil
}

// Query выполняет SQL-запрос, который возвращает набор строк.
// Принимает контекст запроса ctx, текст SQL и позиционные параметры args.
// Возвращает pgx.Rows для чтения результата и ошибку выполнения (если произошла).
func (d *DB) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return d.pool.Query(ctx, sql, args...)
}

// QueryRow выполняет SQL-запрос, который ожидает одну строку результата.
// Принимает контекст запроса ctx, текст SQL и позиционные параметры args.
// Возвращает pgx.Row, из которого можно считать значения; фактическая ошибка
// будет доступна при вызове Scan.
func (d *DB) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return d.pool.QueryRow(ctx, sql, args...)
}

// Exec выполняет SQL-команду, которая не возвращает строки (например, INSERT/UPDATE/DELETE).
// Принимает контекст запроса ctx, текст SQL и позиционные параметры args.
// Возвращает pgconn.CommandTag (метаданные о выполненной команде) и ошибку выполнения (если произошла).
func (d *DB) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	return d.pool.Exec(ctx, sql, args...)
}

func (d *DB) Close(_ context.Context) error {
	d.pool.Close()
	return nil
}
