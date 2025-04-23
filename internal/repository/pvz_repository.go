package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

func (pr *PostgresRepository) ListAllPVZ(ctx context.Context) ([]*PVZ, error) {
	var pvzList []*PVZ
	err := pr.db.SelectContext(ctx, &pvzList, `SELECT * FROM pvz`)
	if err != nil {
		return nil, fmt.Errorf("error listing pvz: %w", err)
	}

	return pvzList, nil
}

func (pr *PostgresRepository) ListPVZ(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]*PVZWithReceptions, error) {
	offset := (page - 1) * limit

	var pvzList []*PVZ

	query := `
        SELECT DISTINCT p.id, p.registration_date, p.city
        FROM pvz p
        JOIN reception r ON p.id = r.pvz_id
    `

	args := make([]interface{}, 0)
	whereClause := ""

	if startDate != nil && endDate != nil {
		whereClause = "WHERE r.execution_date BETWEEN $1 AND $2"
		args = append(args, startDate, endDate)
	} else if startDate != nil {
		whereClause = "WHERE r.execution_date >= $1"
		args = append(args, startDate)
	} else if endDate != nil {
		whereClause = "WHERE r.execution_date <= $1"
		args = append(args, endDate)
	}

	query += " " + whereClause + " ORDER BY p.registration_date DESC"

	offsetPos := len(args) + 1
	limitPos := len(args) + 2
	query += fmt.Sprintf(" OFFSET $%d LIMIT $%d", offsetPos, limitPos)
	args = append(args, offset, limit)

	err := pr.db.SelectContext(ctx, &pvzList, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error listing pvz: %w", err)
	}

	pvzIDs := make([]string, len(pvzList))
	for i := range pvzList {
		pvzIDs[i] = pvzList[i].ID
	}

	var rcList []*Reception
	err = pr.db.SelectContext(
		ctx,
		&rcList,
		`SELECT id, execution_date, pvz_id, status
        FROM reception
        WHERE pvz_id = ANY($1)`,
		pq.Array(pvzIDs),
	)
	if err != nil {
		return nil, fmt.Errorf("error listing receptions: %w", err)
	}

	receptionsByPVZ := make(map[string][]*Reception)
	for _, rc := range rcList {
		receptionsByPVZ[rc.PVZID] = append(receptionsByPVZ[rc.PVZID], rc)
	}

	rcIDs := make([]string, len(rcList))
	for i := range rcList {
		rcIDs[i] = rcList[i].ID
	}

	productList := make([]*Product, 0)
	err = pr.db.SelectContext(
		ctx,
		&productList,
		`SELECT id, type, reception_date, reception_id
		FROM product
		WHERE reception_id = ANY($1)`,
		pq.Array(rcIDs),
	)
	if err != nil {
		return nil, fmt.Errorf("error listing products: %w", err)
	}

	productsByReception := make(map[string][]*Product)
	for _, product := range productList {
		productsByReception[product.ReceptionId] = append(productsByReception[product.ReceptionId], product)
	}

	pvzWithReceptions := make([]*PVZWithReceptions, len(pvzList))
	for i, pvz := range pvzList {
		pvzWithReceptions[i] = &PVZWithReceptions{
			PVZ:        pvz,
			Receptions: make([]*ReceptionWithProducts, len(receptionsByPVZ[pvz.ID])),
		}
		for j, rc := range receptionsByPVZ[pvz.ID] {
			pvzWithReceptions[i].Receptions[j] = &ReceptionWithProducts{
				Reception: rc,
				Products:  productsByReception[rc.ID],
			}
		}
	}

	return pvzWithReceptions, nil
}

func (pr *PostgresRepository) CreatePVZ(ctx context.Context, city string) (*PVZ, error) {
	var registrationDate = time.Now()
	newID := uuid.New().String()
	_, err := pr.db.ExecContext(
		ctx,
		`INSERT INTO pvz (id, city, registration_date) VALUES ($1, $2, $3)`,
		newID,
		city,
		registrationDate,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating pvz: %w", err)
	}

	pvz := &PVZ{
		ID:               newID,
		RegistrationDate: registrationDate,
		City:             city,
	}

	return pvz, nil
}
