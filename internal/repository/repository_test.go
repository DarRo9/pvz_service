package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

var dummyDate time.Time = time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)

func withMockRepository(t *testing.T, fn func(Repository, sqlmock.Sqlmock)) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "postgres")
	var r Repository = NewPostgresRepository(db)
	fn(r, mock)
}

func TestExecTx(t *testing.T) {
	testCases := []struct {
		name string
		test func(*testing.T, Repository, sqlmock.Sqlmock)
	}{
		{
			name: "Error starting transaction",
			test: func(t *testing.T, r Repository, mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(
					fmt.Errorf("error starting transaction"),
				)

				err := r.ExecTx(context.Background(), func(tx *sqlx.Tx) error {
					return nil
				},
				)
				require.Error(t, err)
				mock.ExpectationsWereMet()
			},
		},
		{
			name: "Error committing transaction",
			test: func(t *testing.T, r Repository, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectCommit().WillReturnError(
					fmt.Errorf("error committing transaction"),
				)
				err := r.ExecTx(context.Background(), func(tx *sqlx.Tx) error {
					return nil
				},
				)
				require.Error(t, err)
				mock.ExpectationsWereMet()
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
