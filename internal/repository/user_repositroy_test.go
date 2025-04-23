package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	const query = `INSERT INTO users (id, email, password, role, registration_date) VALUES ($1, $2, $3, $4, $5)`

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

				user, err := r.CreateUser(context.Background(), "example@mail.com", "pass", "employee")
				require.NoError(t, err)
				require.Equal(t, "example@mail.com", user.Email)
				require.Equal(t, "pass", user.Password)
				require.Equal(t, "employee", user.Role)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "Error",
			test: func(t *testing.T, r Repository, mock sqlmock.Sqlmock) {
				mock.ExpectExec(
					query,
				).WillReturnError(
					fmt.Errorf("error inserting user"),
				)

				_, err := r.CreateUser(context.Background(), "example@mail.com", "pass", "employee")
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

func TestListUser(t *testing.T) {

	const query = `SELECT * FROM users`

	users := []*User{
		{
			ID:               "1",
			Email:            "example@mail.com",
			Password:         "pass",
			Role:             "employee",
			RegistrationDate: dummyDate,
		},
		{
			ID:               "2",
			Email:            "dummy@mail.com",
			Password:         "qwerty",
			Role:             "moderator",
			RegistrationDate: dummyDate,
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
					[]string{"id", "email", "password", "role", "registration_date"},
				)
				for _, user := range users {
					rows.AddRow(user.ID, user.Email, user.Password, user.Role, user.RegistrationDate)
				}

				mock.ExpectQuery(
					query,
				).WillReturnRows(rows)

				result, err := r.ListUser(context.Background())
				require.NoError(t, err)
				require.Equal(t, result, users)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "Error listing",
			test: func(t *testing.T, r Repository, mock sqlmock.Sqlmock) {

				mock.ExpectQuery(
					query,
				).WillReturnError(fmt.Errorf("error listing users"))

				_, err := r.ListUser(context.Background())
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

func TestGetUserByEmail(t *testing.T) {

	const query = `SELECT * FROM users WHERE email = $1`

	testCases := []struct {
		name string
		test func(*testing.T, Repository, sqlmock.Sqlmock)
	}{
		{
			name: "Success",
			test: func(t *testing.T, r Repository, mock sqlmock.Sqlmock) {

				mock.ExpectQuery(
					query,
				).WillReturnRows(
					sqlmock.NewRows(
						[]string{"id", "email", "password", "role", "registration_date"},
					).AddRow(
						1,
						"example@mail.com",
						"pass",
						"employee",
						dummyDate,
					),
				)

				result, err := r.GetUserByEmail(context.Background(), "example@mail.com")
				require.NoError(t, err)
				require.Equal(t, "1", result.ID)
				require.Equal(t, "example@mail.com", result.Email)
				require.Equal(t, "pass", result.Password)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "Error",
			test: func(t *testing.T, r Repository, mock sqlmock.Sqlmock) {

				mock.ExpectQuery(
					query,
				).WillReturnError(
					fmt.Errorf("error getting user by email"),
				)
				_, err := r.GetUserByEmail(context.Background(), "example@mail.com")
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
