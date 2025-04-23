package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func (pr *PostgresRepository) ListProducts(ctx context.Context, receptionID string) ([]*Product, error) {
	var products []*Product
	err := pr.db.SelectContext(ctx, &products, `SELECT * FROM product WHERE reception_id = $1`, receptionID)
	if err != nil {
		return nil, fmt.Errorf("error listing products: %w", err)
	}

	return products, nil
}

func (pr *PostgresRepository) CreateProduct(ctx context.Context, receptionID string, productType string) (*Product, error) {
	product := &Product{}
	err := pr.ExecTx(
		ctx,
		func(tx *sqlx.Tx) error {
			var lastReception Reception
			err := tx.QueryRowContext(
				ctx,
				`SELECT id, execution_date, pvz_id, status FROM reception 
				ORDER BY execution_date DESC 
				LIMIT 1
				FOR UPDATE`,
			).Scan(
				&lastReception.ID,
				&lastReception.ExecutionDate,
				&lastReception.PVZID,
				&lastReception.Status,
			)
			isNoReceptions := errors.Is(err, sql.ErrNoRows)
			if err != nil && !isNoReceptions {
				return fmt.Errorf("error getting last reception: %w", err)
			}

			if isNoReceptions {
				return fmt.Errorf("no receptions found")
			}

			if lastReception.Status == closeReceptionStatus {
				return fmt.Errorf("last reception is closed")
			}

			receptionDate := time.Now()
			newID := uuid.New().String()
			_, err = tx.ExecContext(ctx,
				`INSERT INTO product (id, reception_date, reception_id, type)
				VALUES ($1, $2, $3, $4)`,
				newID, receptionDate, lastReception.ID, productType,
			)

			if err != nil {
				return fmt.Errorf("error inserting reception: %w", err)
			}
			product.ID = newID
			product.ReceptionDate = receptionDate
			product.Type = productType
			product.ReceptionId = lastReception.ID
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error creating product: %w", err)
	}

	return product, nil
}

func (pr *PostgresRepository) DeleteProduct(ctx context.Context, PVZID string) (*Product, error) {
	product := &Product{}
	err := pr.ExecTx(
		ctx,
		func(tx *sqlx.Tx) error {
			var lastReception Reception
			err := tx.QueryRowContext(
				ctx,
				`SELECT id, execution_date, pvz_id, status FROM reception
				WHERE pvz_id = $1
				ORDER BY execution_date DESC
				LIMIT 1
				FOR UPDATE`,
				PVZID,
			).Scan(
				&lastReception.ID,
				&lastReception.ExecutionDate,
				&lastReception.PVZID,
				&lastReception.Status,
			)
			isNoReceptions := errors.Is(err, sql.ErrNoRows)
			if err != nil && !isNoReceptions {
				return fmt.Errorf("error getting last reception: %w", err)
			}

			if isNoReceptions {
				return fmt.Errorf("no receptions found")
			}

			if lastReception.Status == closeReceptionStatus {
				return fmt.Errorf("last reception is closed")
			}

			err = tx.QueryRowContext(ctx,
				`DELETE FROM product
				WHERE id = (
					SELECT id
					FROM product
					WHERE reception_id = $1
					ORDER BY reception_date DESC
					LIMIT 1
				)
				RETURNING id, type, reception_date, reception_id`,
				lastReception.ID,
			).Scan(
				&product.ID,
				&product.Type,
				&product.ReceptionDate,
				&product.ReceptionId,
			)

			isNoProducts := errors.Is(err, sql.ErrNoRows)
			if err != nil && !isNoProducts {
				return fmt.Errorf("error deleting product: %w", err)
			}

			if isNoProducts {
				return fmt.Errorf("no products found")
			}

			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error creating product: %w", err)
	}

	return product, nil
}
