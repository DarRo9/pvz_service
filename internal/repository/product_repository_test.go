package repository

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestListProduct(t *testing.T) {
	const query = `SELECT * FROM product WHERE reception_id = $1`

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
	}

	testCases := []struct {
		name string
		test func(*testing.T, Repository, sqlmock.Sqlmock)
	}{
		{
			name: "Error",
			test: func(t *testing.T, r Repository, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(
					query,
				).WillReturnError(
					fmt.Errorf("error listing products"),
				)
				_, err := r.ListProducts(context.Background(), "1")
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "Success",
			test: func(t *testing.T, r Repository, mock sqlmock.Sqlmock) {

				rows := sqlmock.NewRows(
					[]string{"id", "reception_date", "reception_id", "type"},
				)
				for _, rc := range products {
					rows.AddRow(rc.ID, rc.ReceptionDate, rc.ReceptionId, rc.Type)
				}
				mock.ExpectQuery(
					query,
				).WillReturnRows(
					rows,
				)

				result, err := r.ListProducts(context.Background(), "1")
				require.NoError(t, err)
				require.Equal(t, result, products)

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

func TestCreateProduct(t *testing.T) {
	const query1 = `SELECT id, execution_date, pvz_id, status FROM reception 
		ORDER BY execution_date DESC 
		LIMIT 1
		FOR UPDATE`
	const query2 = `INSERT INTO product (id, reception_date, reception_id, type)
		VALUES ($1, $2, $3, $4)`

	testCases := []struct {
		name string
		test func(*testing.T, Repository, sqlmock.Sqlmock)
	}{
		{
			name: "Success",
			test: func(t *testing.T, r Repository, mock sqlmock.Sqlmock) {

				mock.ExpectBegin()
				mock.ExpectQuery(
					query1,
				).WillReturnRows(
					sqlmock.NewRows([]string{"id", "execution_date", "pvz_id", "status"}).AddRow(
						1,
						dummyDate,
						1,
						inProgressReceptionStatus,
					),
				)
				mock.ExpectExec(
					query2,
				).WillReturnResult(
					sqlmock.NewResult(1, 1),
				)
				mock.ExpectCommit()

				result, err := r.CreateProduct(context.Background(), "1", "product_type")
				require.NoError(t, err)
				require.Equal(t, "product_type", result.Type)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "Error creating product with no receptions",
			test: func(t *testing.T, r Repository, mock sqlmock.Sqlmock) {

				mock.ExpectBegin()
				mock.ExpectQuery(
					query1,
				).WillReturnError(
					sql.ErrNoRows,
				)
				mock.ExpectRollback()

				_, err := r.CreateProduct(context.Background(), "1", "product_type")
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "Error creating product with closed reception",
			test: func(t *testing.T, r Repository, mock sqlmock.Sqlmock) {

				mock.ExpectBegin()
				mock.ExpectQuery(
					query1,
				).WillReturnRows(
					sqlmock.NewRows([]string{"id", "execution_date", "pvz_id", "status"}).AddRow(
						1,
						dummyDate,
						1,
						closeReceptionStatus,
					),
				)
				mock.ExpectRollback()

				_, err := r.CreateProduct(context.Background(), "1", "product_type")
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "Error creating product with query error",
			test: func(t *testing.T, r Repository, mock sqlmock.Sqlmock) {

				mock.ExpectBegin()
				mock.ExpectQuery(
					query1,
				).WillReturnError(
					fmt.Errorf("error getting last reception"),
				)
				mock.ExpectRollback()

				_, err := r.CreateProduct(context.Background(), "1", "product_type")
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "Error creating product with insert error",
			test: func(t *testing.T, r Repository, mock sqlmock.Sqlmock) {

				mock.ExpectBegin()
				mock.ExpectQuery(
					query1,
				).WillReturnRows(
					sqlmock.NewRows([]string{"id", "execution_date", "pvz_id", "status"}).AddRow(
						1,
						dummyDate,
						1,
						inProgressReceptionStatus,
					),
				)
				mock.ExpectExec(
					query2,
				).WillReturnError(
					fmt.Errorf("error inserting product"),
				)

				mock.ExpectRollback()

				_, err := r.CreateProduct(context.Background(), "1", "product_type")
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

func TestDeleteProduct(t *testing.T) {

	const query1 = `SELECT id, execution_date, pvz_id, status FROM reception
		WHERE pvz_id = $1
		ORDER BY execution_date DESC
		LIMIT 1
		FOR UPDATE`
	const query2 = `DELETE FROM product
		WHERE id = (
			SELECT id
			FROM product
			WHERE reception_id = $1
			ORDER BY reception_date DESC
			LIMIT 1
		)
		RETURNING id, type, reception_date, reception_id`

	testCases := []struct {
		name string
		test func(*testing.T, Repository, sqlmock.Sqlmock)
	}{
		{
			name: "Success",
			test: func(t *testing.T, r Repository, mock sqlmock.Sqlmock) {

				mock.ExpectBegin()
				mock.ExpectQuery(query1).WillReturnRows(
					sqlmock.NewRows([]string{"id", "execution_date", "pvz_id", "status"}).AddRow(
						1,
						dummyDate,
						1,
						inProgressReceptionStatus,
					),
				)
				mock.ExpectQuery(
					query2,
				).WillReturnRows(
					sqlmock.NewRows([]string{"id", "type", "reception_date", "reception_id"}).AddRow(
						1,
						"product_type",
						dummyDate,
						1,
					),
				)
				mock.ExpectCommit()

				result, err := r.DeleteProduct(context.Background(), "1")
				require.NoError(t, err)
				require.Equal(t, "1", result.ID)
				require.Equal(t, "product_type", result.Type)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "Error deleting product with no receptions",
			test: func(t *testing.T, r Repository, mock sqlmock.Sqlmock) {

				mock.ExpectBegin()
				mock.ExpectQuery(
					query1,
				).WillReturnError(
					sql.ErrNoRows,
				)
				mock.ExpectRollback()

				_, err := r.DeleteProduct(context.Background(), "1")
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "Error deleting product with closed reception",
			test: func(t *testing.T, r Repository, mock sqlmock.Sqlmock) {

				mock.ExpectBegin()
				mock.ExpectQuery(
					query1,
				).WillReturnRows(
					sqlmock.NewRows([]string{"id", "execution_date", "pvz_id", "status"}).AddRow(
						1,
						dummyDate,
						1,
						closeReceptionStatus,
					),
				)
				mock.ExpectRollback()

				_, err := r.DeleteProduct(context.Background(), "1")
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "Error deleting product with no products",
			test: func(t *testing.T, r Repository, mock sqlmock.Sqlmock) {

				mock.ExpectBegin()
				mock.ExpectQuery(
					query1,
				).WillReturnRows(
					sqlmock.NewRows([]string{"id", "execution_date", "pvz_id", "status"}).AddRow(
						1,
						dummyDate,
						1,
						inProgressReceptionStatus,
					),
				)
				mock.ExpectQuery(
					query2,
				).WillReturnError(
					sql.ErrNoRows,
				)
				mock.ExpectRollback()

				_, err := r.DeleteProduct(context.Background(), "1")
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "Error deleting product with query error",
			test: func(t *testing.T, r Repository, mock sqlmock.Sqlmock) {

				mock.ExpectBegin()
				mock.ExpectQuery(
					query1,
				).WillReturnError(
					fmt.Errorf("error getting last reception"),
				)
				mock.ExpectRollback()

				_, err := r.DeleteProduct(context.Background(), "1")
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "Error deleting product with delete error",
			test: func(t *testing.T, r Repository, mock sqlmock.Sqlmock) {

				mock.ExpectBegin()
				mock.ExpectQuery(
					query1,
				).WillReturnRows(
					sqlmock.NewRows([]string{"id", "execution_date", "pvz_id", "status"}).AddRow(
						1,
						dummyDate,
						1,
						inProgressReceptionStatus,
					),
				)
				mock.ExpectQuery(
					query2,
				).WillReturnError(
					fmt.Errorf("error deleting product"),
				)
				mock.ExpectRollback()

				_, err := r.DeleteProduct(context.Background(), "1")
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
