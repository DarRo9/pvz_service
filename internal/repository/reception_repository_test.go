package repository

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestListReception(t *testing.T) {

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
	}

	query := `SELECT * FROM reception WHERE pvz_id = $1`

	testCases := []struct {
		name string
		test func(*testing.T, Repository, sqlmock.Sqlmock)
	}{
		{
			name: "Success",
			test: func(t *testing.T, r Repository, mock sqlmock.Sqlmock) {

				rows := sqlmock.NewRows(
					[]string{"id", "pvz_id", "status", "execution_date"},
				)
				for _, rc := range receptons {
					rows.AddRow(rc.ID, rc.PVZID, rc.Status, rc.ExecutionDate)
				}

				mock.ExpectQuery(
					query,
				).WillReturnRows(rows)

				result, err := r.ListReception(context.Background(), "1")
				require.NoError(t, err)
				require.Equal(t, result, receptons)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "Error listing",
			test: func(t *testing.T, r Repository, mock sqlmock.Sqlmock) {

				mock.ExpectQuery(
					query,
				).WillReturnError(fmt.Errorf("error listing receptions"))

				_, err := r.ListReception(context.Background(), "1")
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

func TestCreateReception(t *testing.T) {
	query1 := `SELECT status FROM reception 
		ORDER BY execution_date DESC 
		LIMIT 1
		FOR UPDATE`
	query2 := `INSERT INTO reception (id, execution_date, pvz_id, status)
		VALUES ($1, $2, $3, $4)`

	testCases := []struct {
		name string
		test func(*testing.T, Repository, sqlmock.Sqlmock)
	}{
		{
			name: "Successful create with no previous receptions",
			test: func(t *testing.T, r Repository, mock sqlmock.Sqlmock) {

				mock.ExpectBegin()
				mock.ExpectQuery(
					query1,
				).WillReturnError(
					sql.ErrNoRows,
				)
				mock.ExpectExec(
					query2,
				).WillReturnResult(
					sqlmock.NewResult(1, 1),
				)
				mock.ExpectCommit()

				rc, err := r.CreateReception(context.Background(), "1")
				require.NoError(t, err)
				require.Equal(t, inProgressReceptionStatus, rc.Status)
				require.Equal(t, "1", rc.PVZID)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "Successful create with previous closed reception",
			test: func(t *testing.T, r Repository, mock sqlmock.Sqlmock) {

				mock.ExpectBegin()
				mock.ExpectQuery(
					query1,
				).WillReturnRows(
					sqlmock.NewRows([]string{"status"}).AddRow(closeReceptionStatus),
				)
				mock.ExpectExec(
					query2,
				).WillReturnResult(
					sqlmock.NewResult(1, 1),
				)
				mock.ExpectCommit()

				rc, err := r.CreateReception(context.Background(), "1")
				require.NoError(t, err)
				require.Equal(t, inProgressReceptionStatus, rc.Status)
				require.Equal(t, "1", rc.PVZID)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "Error creating reception with previous in_progress reception",
			test: func(t *testing.T, r Repository, mock sqlmock.Sqlmock) {

				mock.ExpectBegin()
				mock.ExpectQuery(
					query1,
				).WillReturnRows(
					sqlmock.NewRows([]string{"status"}).AddRow(inProgressReceptionStatus),
				)
				mock.ExpectRollback()

				_, err := r.CreateReception(context.Background(), "1")
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "Error querying last reception status",
			test: func(t *testing.T, r Repository, mock sqlmock.Sqlmock) {

				mock.ExpectBegin()
				mock.ExpectQuery(
					query1,
				).WillReturnError(
					fmt.Errorf("error getting last reception status"),
				)
				mock.ExpectRollback()

				_, err := r.CreateReception(context.Background(), "1")
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "Error inserting reception",
			test: func(t *testing.T, r Repository, mock sqlmock.Sqlmock) {

				mock.ExpectBegin()
				mock.ExpectQuery(
					query1,
				).WillReturnError(
					sql.ErrNoRows,
				)
				mock.ExpectExec(
					query2,
				).WillReturnError(
					fmt.Errorf("error inserting reception"),
				)
				mock.ExpectRollback()

				_, err := r.CreateReception(context.Background(), "1")
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

func TestCloseReception(t *testing.T) {
	testCases := []struct {
		name string
		test func(*testing.T, Repository, sqlmock.Sqlmock)
	}{
		{
			name: "Success",
			test: func(t *testing.T, r Repository, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(
					`SELECT id, execution_date, pvz_id, status FROM reception 
					ORDER BY execution_date DESC 
					LIMIT 1`,
				).WillReturnRows(
					sqlmock.NewRows([]string{"id", "execution_date", "pvz_id", "status"}).AddRow(
						1,
						dummyDate,
						1,
						inProgressReceptionStatus,
					),
				)
				mock.ExpectExec(
					`UPDATE reception
					SET status = $1
					WHERE id = $2`,
				).WillReturnResult(
					sqlmock.NewResult(1, 1),
				)

				rc, err := r.CloseReception(context.Background(), "1")
				require.NoError(t, err)
				require.Equal(t, &Reception{
					ID:            "1",
					PVZID:         "1",
					Status:        closeReceptionStatus,
					ExecutionDate: dummyDate,
				}, rc)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "Error closing reception with no receptions",
			test: func(t *testing.T, r Repository, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(
					`SELECT id, execution_date, pvz_id, status FROM reception 
					ORDER BY execution_date DESC 
					LIMIT 1`,
				).WillReturnError(
					sql.ErrNoRows,
				)

				_, err := r.CloseReception(context.Background(), "1")
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "Error closing reception with already closed reception",
			test: func(t *testing.T, r Repository, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(
					`SELECT id, execution_date, pvz_id, status FROM reception 
					ORDER BY execution_date DESC 
					LIMIT 1`,
				).WillReturnRows(
					sqlmock.NewRows([]string{"id", "execution_date", "pvz_id", "status"}).AddRow(
						1,
						dummyDate,
						1,
						closeReceptionStatus,
					),
				)

				_, err := r.CloseReception(context.Background(), "1")
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "Error querying last reception",
			test: func(t *testing.T, r Repository, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(
					`SELECT id, execution_date, pvz_id, status FROM reception 
					ORDER BY execution_date DESC 
					LIMIT 1`,
				).WillReturnError(
					fmt.Errorf("error getting last reception status"),
				)

				_, err := r.CloseReception(context.Background(), "1")
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "Error updating reception",
			test: func(t *testing.T, r Repository, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(
					`SELECT id, execution_date, pvz_id, status FROM reception 
					ORDER BY execution_date DESC 
					LIMIT 1`,
				).WillReturnRows(
					sqlmock.NewRows([]string{"id", "execution_date", "pvz_id", "status"}).AddRow(
						1,
						dummyDate,
						1,
						inProgressReceptionStatus,
					),
				)
				mock.ExpectExec(
					`UPDATE reception
					SET status = $1
					WHERE id = $2`,
				).WillReturnError(
					fmt.Errorf("error updating reception status"),
				)

				_, err := r.CloseReception(context.Background(), "1")
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
