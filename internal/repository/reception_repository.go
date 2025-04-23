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

const (
	closeReceptionStatus      = "close"
	inProgressReceptionStatus = "in_progress"
)

func (pr *PostgresRepository) ListReception(ctx context.Context, PVZID string) ([]*Reception, error) {
	var receptions []*Reception
	err := pr.db.SelectContext(ctx, &receptions, `SELECT * FROM reception WHERE pvz_id = $1`, PVZID)
	if err != nil {
		return nil, fmt.Errorf("error listing receptions: %w", err)
	}

	return receptions, nil
}

func (pr *PostgresRepository) CreateReception(ctx context.Context, PVZID string) (*Reception, error) {
	rc := &Reception{}
	err := pr.ExecTx(
		ctx,
		func(tx *sqlx.Tx) error {
			var lastReceptionStatus string
			err := tx.QueryRowContext(
				ctx, `
				SELECT status FROM reception 
				ORDER BY execution_date DESC 
				LIMIT 1
				FOR UPDATE`,
			).Scan(&lastReceptionStatus)

			isNoReceptions := errors.Is(err, sql.ErrNoRows)
			if err != nil && !isNoReceptions {
				return fmt.Errorf("error getting last reception status: %w", err)
			}

			if !isNoReceptions && lastReceptionStatus != "close" {
				return fmt.Errorf("last reception is not closed: %s", lastReceptionStatus)
			}

			executionDate := time.Now()
			newID := uuid.New().String()
			_, err = tx.ExecContext(ctx,
				`INSERT INTO reception (id, execution_date, pvz_id, status)
				VALUES ($1, $2, $3, $4)`,
				newID, executionDate, PVZID, inProgressReceptionStatus,
			)

			if err != nil {
				return fmt.Errorf("error inserting reception: %w", err)
			}
			rc.ID = newID
			rc.ExecutionDate = executionDate
			rc.PVZID = PVZID
			rc.Status = inProgressReceptionStatus
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error creating reception: %w", err)
	}
	return rc, nil
}

func (pr *PostgresRepository) CloseReception(ctx context.Context, PVZID string) (*Reception, error) {
	var lastReception Reception
	err := pr.db.QueryRowContext(
		ctx,
		`SELECT id, execution_date, pvz_id, status FROM reception 
		ORDER BY execution_date DESC 
		LIMIT 1`,
	).Scan(
		&lastReception.ID,
		&lastReception.ExecutionDate,
		&lastReception.PVZID,
		&lastReception.Status,
	)

	isNoReceptions := errors.Is(err, sql.ErrNoRows)
	if err != nil && !isNoReceptions {
		return nil, fmt.Errorf("error getting last reception status: %w", err)
	}

	if isNoReceptions {
		return nil, fmt.Errorf("no receptions found")
	}

	if lastReception.Status == closeReceptionStatus {
		return nil, fmt.Errorf("last reception is already closed")
	}

	_, err = pr.db.ExecContext(ctx,
		`UPDATE reception
		SET status = $1
		WHERE id = $2`,
		closeReceptionStatus,
		lastReception.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("error updating reception status: %w", err)
	}

	lastReception.Status = closeReceptionStatus
	return &lastReception, nil
}
