package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"

	"migrations/internal/model"
)

type DB struct {
	pool *pgxpool.Pool
}

func NewDB(dsn string) (*DB, error) {
	if err := runMigrations(dsn); err != nil {
		return nil, fmt.Errorf("failed to run DB migrations: %w", err)
	}
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create a connection pool: %w", err)
	}
	return &DB{
		pool: pool,
	}, nil
}

func runMigrations(dsn string) error {
	const migrationsPath = "../db/migrations"

	m, err := migrate.New(fmt.Sprintf("file://%s", migrationsPath), dsn)
	if err != nil {
		return fmt.Errorf("failed to get a new migrate instance: %w", err)
	}
	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("failed to apply migrations to the DB: %w", err)
		}
	}
	return nil
}

func (db *DB) PutEmployee(ctx context.Context, emp *model.Employee) error {
	tag, err := db.pool.Exec(
		ctx,
		`INSERT INTO employees(first_name, last_name, salary, position, email)
		VALUES ($1, $2, $3, (SELECT id FROM positions WHERE title=$4), $5)`,
		emp.FirstName, emp.LastName, emp.Salary, emp.Position, emp.Email,
	)
	if err != nil {

		return fmt.Errorf("failed to store employee: %w", err)
	}
	rowsAffectedCount := tag.RowsAffected()
	if rowsAffectedCount != 1 {
		return fmt.Errorf("expected one row to be affected, actually affected %d", rowsAffectedCount)
	}
	return nil
}
