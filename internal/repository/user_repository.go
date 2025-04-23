package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (pr *PostgresRepository) ListUser(ctx context.Context) ([]*User, error) {
	var users []*User
	err := pr.db.SelectContext(ctx, &users, `SELECT * FROM users`)
	if err != nil {
		return nil, fmt.Errorf("error listing users: %w", err)
	}

	return users, nil
}

func (pr *PostgresRepository) CreateUser(ctx context.Context, email, password, role string) (*User, error) {
	var registrationDate = time.Now()
	newID := uuid.New().String()
	_, err := pr.db.ExecContext(
		ctx,
		`INSERT INTO users (id, email, password, role, registration_date) VALUES ($1, $2, $3, $4, $5)`,
		newID,
		email,
		password,
		role,
		registrationDate,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	user := &User{
		ID:               newID,
		Email:            email,
		Password:         password,
		Role:             role,
		RegistrationDate: registrationDate,
	}

	return user, nil
}

func (pr *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := pr.db.GetContext(ctx, &user, `SELECT * FROM users WHERE email = $1`, email)
	if err != nil {
		return nil, fmt.Errorf("error getting user by email: %w", err)
	}

	return &user, nil
}
