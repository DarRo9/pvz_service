package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	ExecTx(ctx context.Context, fn func(*sqlx.Tx) error) error

	// PVZ
	ListPVZ(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]*PVZWithReceptions, error)
	ListAllPVZ(ctx context.Context) ([]*PVZ, error)
	CreatePVZ(ctx context.Context, city string) (*PVZ, error)

	// Reception
	CreateReception(ctx context.Context, PVZID string) (*Reception, error)
	CloseReception(ctx context.Context, PVZID string) (*Reception, error)
	ListReception(ctx context.Context, PVZID string) ([]*Reception, error)

	// Product
	ListProducts(ctx context.Context, receptionID string) ([]*Product, error)
	CreateProduct(ctx context.Context, receptionID string, productType string) (*Product, error)
	DeleteProduct(ctx context.Context, PVZID string) (*Product, error)

	// User
	ListUser(ctx context.Context) ([]*User, error)
	CreateUser(ctx context.Context, email, password, role string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}

type PostgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (pr *PostgresRepository) ExecTx(ctx context.Context, fn func(*sqlx.Tx) error) error {
	tx, err := pr.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

	err = fn(tx)
	if err != nil {
		return fmt.Errorf("error executing transaction function: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}
