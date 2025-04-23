package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestCreatePVZ(t *testing.T) {
	query := `INSERT INTO pvz (id, city, registration_date) VALUES ($1, $2, $3)`

	testCases := []struct {
		name string
		test func(*testing.T, Repository, sqlmock.Sqlmock)
	}{
		{
			name: "Success",
			test: func(t *testing.T, r Repository, mock sqlmock.Sqlmock) {
				mock.ExpectExec(
					query,
				).WillReturnResult(
					sqlmock.NewResult(1, 1),
				)
				pvz, err := r.CreatePVZ(context.Background(), "Moscow")
				require.NoError(t, err)
				require.Equal(t, pvz.City, "Moscow")

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "Error inserting",
			test: func(t *testing.T, r Repository, mock sqlmock.Sqlmock) {
				mock.ExpectExec(
					query,
				).WillReturnError(
					fmt.Errorf("error inserting pvz"),
				)

				_, err := r.CreatePVZ(context.Background(), "Moscow")
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			withMockRepository(t, func(r Repository, mock sqlmock.Sqlmock) {
				tc.test(t, r, mock)
			})
		})
	}
}

func TestListPVZ(t *testing.T) {
	pvzs := []*PVZ{
		{
			ID:               "1",
			City:             "Moscow",
			RegistrationDate: dummyDate,
		},
		{
			ID:               "2",
			City:             "Saint Petersburg",
			RegistrationDate: dummyDate,
		},
	}
	receptons := []*Reception{
		&Reception{
			ID:            "1",
			PVZID:         "1",
			Status:        "in_progress",
			ExecutionDate: dummyDate,
		},
		&Reception{
			ID:            "2",
			PVZID:         "1",
			Status:        "close",
			ExecutionDate: dummyDate,
		},
		&Reception{
			ID:            "3",
			PVZID:         "2",
			Status:        "close",
			ExecutionDate: dummyDate,
		},
	}
	products := []*Product{
		&Product{
			ID:            "1",
			ReceptionDate: dummyDate,
			ReceptionId:   "1",
			Type:          "product_type",
		},
		&Product{
			ID:            "2",
			ReceptionDate: dummyDate,
			ReceptionId:   "1",
			Type:          "product_type_2",
		},
		&Product{
			ID:            "2",
			ReceptionDate: dummyDate,
			ReceptionId:   "3",
			Type:          "product_type_2",
		},
	}

	query1 := `SELECT DISTINCT p.id, p.registration_date, p.city
        FROM pvz p
        JOIN reception r ON p.id = r.pvz_id
		ORDER BY p.registration_date DESC
		OFFSET $1 LIMIT $2
		`

	query2 := `SELECT id, execution_date, pvz_id, status
        FROM reception
        WHERE pvz_id = ANY($1)`

	query3 := `SELECT id, type, reception_date, reception_id
		FROM product
		WHERE reception_id = ANY($1)`

	expected := []*PVZWithReceptions{
		{
			PVZ: pvzs[0],
			Receptions: []*ReceptionWithProducts{
				{
					Reception: receptons[0],
					Products: []*Product{
						products[0],
						products[1],
					},
				},
				{
					Reception: receptons[1],
					Products:  nil,
				},
			},
		},
		{
			PVZ: pvzs[1],
			Receptions: []*ReceptionWithProducts{
				{
					Reception: receptons[2],
					Products: []*Product{
						products[2],
					},
				},
			},
		},
	}

	testCases := []struct {
		name string
		test func(*testing.T, Repository, sqlmock.Sqlmock)
	}{
		{
			name: "Success",
			test: func(t *testing.T, r Repository, mock sqlmock.Sqlmock) {

				rows := sqlmock.NewRows(
					[]string{"id", "city", "registration_date"},
				)
				for _, pvz := range pvzs {
					rows.AddRow(pvz.ID, pvz.City, pvz.RegistrationDate)
				}

				mock.ExpectQuery(
					query1,
				).WillReturnRows(rows)

				rows2 := sqlmock.NewRows(
					[]string{"id", "execution_date", "pvz_id", "status"},
				)
				for _, rc := range receptons {
					rows2.AddRow(rc.ID, rc.ExecutionDate, rc.PVZID, rc.Status)
				}

				mock.ExpectQuery(
					query2,
				).WillReturnRows(rows2)

				rows3 := sqlmock.NewRows(
					[]string{"id", "type", "reception_date", "reception_id"},
				)

				for _, rc := range products {
					rows3.AddRow(rc.ID, rc.Type, rc.ReceptionDate, rc.ReceptionId)
				}
				mock.ExpectQuery(
					query3,
				).WillReturnRows(rows3)

				result, err := r.ListPVZ(context.Background(), nil, nil, 1, 10)
				require.NoError(t, err)
				require.Equal(t, expected, result)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			withMockRepository(t, func(r Repository, mock sqlmock.Sqlmock) {
				tc.test(t, r, mock)
			})
		})
	}
}

func TestListAllPVZ(t *testing.T) {

	pvzs := []*PVZ{
		{
			ID:               "1",
			City:             "Moscow",
			RegistrationDate: dummyDate,
		},
		{
			ID:               "2",
			City:             "Saint Petersburg",
			RegistrationDate: dummyDate,
		},
		{
			ID:               "3",
			City:             "Kazan",
			RegistrationDate: dummyDate,
		},
	}

	query := `SELECT * FROM pvz`

	testCases := []struct {
		name string
		test func(*testing.T, Repository, sqlmock.Sqlmock)
	}{
		{
			name: "Success",
			test: func(t *testing.T, r Repository, mock sqlmock.Sqlmock) {

				rows := sqlmock.NewRows(
					[]string{"id", "city", "registration_date"},
				)
				for _, pvz := range pvzs {
					rows.AddRow(pvz.ID, pvz.City, pvz.RegistrationDate)
				}

				mock.ExpectQuery(
					query,
				).WillReturnRows(rows)

				result, err := r.ListAllPVZ(context.Background())
				require.NoError(t, err)
				require.Equal(t, result, pvzs)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "Error listing",
			test: func(t *testing.T, r Repository, mock sqlmock.Sqlmock) {

				mock.ExpectQuery(
					query,
				).WillReturnError(fmt.Errorf("error listing pvz"))

				_, err := r.ListAllPVZ(context.Background())
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			withMockRepository(t, func(r Repository, mock sqlmock.Sqlmock) {
				tc.test(t, r, mock)
			})
		})
	}
}
