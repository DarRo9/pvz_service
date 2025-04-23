package handler

import (
	"testing"
	"time"

	"github.com/DarRo9/pvz_service/internal/repository"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/assert"
)

func TestProductRepositoryToHTTP(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name     string
		input    *repository.Product
		expected *Product
	}{
		{
			name: "successful conversion",
			input: &repository.Product{
				ID:            "550e8400-e29b-41d4-a716-446655440000",
				ReceptionId:   "550e8400-e29b-41d4-a716-446655440001",
				Type:          "Electronics",
				ReceptionDate: now,
			},
			expected: &Product{
				Id:          func() *uuid.UUID { u, _ := uuid.Parse("550e8400-e29b-41d4-a716-446655440000"); return &u }(),
				ReceptionId: func() uuid.UUID { u, _ := uuid.Parse("550e8400-e29b-41d4-a716-446655440001"); return u }(),
				Type:        ProductType("Electronics"),
				DateTime:    &now,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := productRepositoryToHTTP(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestPVZRepositoryToHTTP(t *testing.T) {
	testCases := []struct {
		name     string
		input    *repository.PVZ
		expected *PVZ
	}{
		{
			name: "successful conversion",
			input: &repository.PVZ{
				ID:   "550e8400-e29b-41d4-a716-446655440000",
				City: "Moscow",
			},
			expected: &PVZ{
				Id:   func() *uuid.UUID { u, _ := uuid.Parse("550e8400-e29b-41d4-a716-446655440000"); return &u }(),
				City: PVZCity("Moscow"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := pvzRepositoryToHTTP(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestReceptionRepositoryToHTTP(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name     string
		input    *repository.Reception
		expected *Reception
	}{
		{
			name: "successful conversion",
			input: &repository.Reception{
				ID:            "550e8400-e29b-41d4-a716-446655440000",
				PVZID:         "550e8400-e29b-41d4-a716-446655440001",
				ExecutionDate: now,
				Status:        "Completed",
			},
			expected: &Reception{
				Id:       func() *uuid.UUID { u, _ := uuid.Parse("550e8400-e29b-41d4-a716-446655440000"); return &u }(),
				PvzId:    func() uuid.UUID { u, _ := uuid.Parse("550e8400-e29b-41d4-a716-446655440001"); return u }(),
				DateTime: now,
				Status:   ReceptionStatus("Completed"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := receptionRepositoryToHTTP(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestPVZWithReceptionsRepositoryToHTTP(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name     string
		input    *repository.PVZWithReceptions
		expected *PVZWithReceptions
	}{
		{
			name: "successful conversion with empty receptions",
			input: &repository.PVZWithReceptions{
				PVZ: &repository.PVZ{
					ID:   "550e8400-e29b-41d4-a716-446655440000",
					City: "Moscow",
				},
				Receptions: []*repository.ReceptionWithProducts{},
			},
			expected: &PVZWithReceptions{
				PVZ: &PVZ{
					Id:   func() *uuid.UUID { u, _ := uuid.Parse("550e8400-e29b-41d4-a716-446655440000"); return &u }(),
					City: PVZCity("Moscow"),
				},
				Receptions: []*ReceptionWithProducts{},
			},
		},
		{
			name: "successful conversion with receptions and products",
			input: &repository.PVZWithReceptions{
				PVZ: &repository.PVZ{
					ID:   "550e8400-e29b-41d4-a716-446655440000",
					City: "Moscow",
				},
				Receptions: []*repository.ReceptionWithProducts{
					{
						Reception: &repository.Reception{
							ID:            "550e8400-e29b-41d4-a716-446655440001",
							PVZID:         "550e8400-e29b-41d4-a716-446655440000",
							ExecutionDate: now,
							Status:        "InProgress",
						},
						Products: []*repository.Product{
							{
								ID:            "550e8400-e29b-41d4-a716-446655440002",
								ReceptionId:   "550e8400-e29b-41d4-a716-446655440001",
								Type:          "Clothing",
								ReceptionDate: now,
							},
						},
					},
				},
			},
			expected: &PVZWithReceptions{
				PVZ: &PVZ{
					Id:   func() *uuid.UUID { u, _ := uuid.Parse("550e8400-e29b-41d4-a716-446655440000"); return &u }(),
					City: PVZCity("Moscow"),
				},
				Receptions: []*ReceptionWithProducts{
					{
						Reception: &Reception{
							Id:       func() *uuid.UUID { u, _ := uuid.Parse("550e8400-e29b-41d4-a716-446655440001"); return &u }(),
							PvzId:    func() uuid.UUID { u, _ := uuid.Parse("550e8400-e29b-41d4-a716-446655440000"); return u }(),
							DateTime: now,
							Status:   ReceptionStatus("InProgress"),
						},
						Products: []*Product{
							{
								Id:          func() *uuid.UUID { u, _ := uuid.Parse("550e8400-e29b-41d4-a716-446655440002"); return &u }(),
								ReceptionId: func() uuid.UUID { u, _ := uuid.Parse("550e8400-e29b-41d4-a716-446655440001"); return u }(),
								Type:        ProductType("Clothing"),
								DateTime:    &now,
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := pvzWithReceptionsRepositoryToHTTP(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestUserRepositoryToHTTP(t *testing.T) {
	testCases := []struct {
		name     string
		input    *repository.User
		expected *User
	}{
		{
			name: "successful conversion",
			input: &repository.User{
				ID:    "550e8400-e29b-41d4-a716-446655440000",
				Email: "test@example.com",
				Role:  "Admin",
			},
			expected: &User{
				Id:    func() *uuid.UUID { u, _ := uuid.Parse("550e8400-e29b-41d4-a716-446655440000"); return &u }(),
				Email: openapi_types.Email("test@example.com"),
				Role:  UserRole("Admin"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := userRepositoryToHTTP(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}
