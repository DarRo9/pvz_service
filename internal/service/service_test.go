package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DarRo9/pvz_service/config"
	"github.com/DarRo9/pvz_service/internal/repository"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// MockRepository реализует интерфейс repository.Repository для тестов
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateUser(ctx context.Context, email, password, role string) (*repository.User, error) {
	args := m.Called(ctx, email, password, role)
	return args.Get(0).(*repository.User), args.Error(1)
}

func (m *MockRepository) GetUserByEmail(ctx context.Context, email string) (*repository.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*repository.User), args.Error(1)
}

func (m *MockRepository) CreatePVZ(ctx context.Context, city string) (*repository.PVZ, error) {
	args := m.Called(ctx, city)
	return args.Get(0).(*repository.PVZ), args.Error(1)
}

func (m *MockRepository) CloseReception(ctx context.Context, pvzId string) (*repository.Reception, error) {
	args := m.Called(ctx, pvzId)
	return args.Get(0).(*repository.Reception), args.Error(1)
}

func (m *MockRepository) DeleteProduct(ctx context.Context, pvzId string) (*repository.Product, error) {
	args := m.Called(ctx, pvzId)
	return args.Get(0).(*repository.Product), args.Error(1)
}

func (m *MockRepository) CreateReception(ctx context.Context, pvzId string) (*repository.Reception, error) {
	args := m.Called(ctx, pvzId)
	return args.Get(0).(*repository.Reception), args.Error(1)
}

func (m *MockRepository) ListPVZ(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]*repository.PVZWithReceptions, error) {
	args := m.Called(ctx, startDate, endDate, page, limit)
	return args.Get(0).([]*repository.PVZWithReceptions), args.Error(1)
}

func (m *MockRepository) CreateProduct(ctx context.Context, receptionID, productType string) (*repository.Product, error) {
	args := m.Called(ctx, receptionID, productType)
	return args.Get(0).(*repository.Product), args.Error(1)
}

func (m *MockRepository) ListAllPVZ(ctx context.Context) ([]*repository.PVZ, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*repository.PVZ), args.Error(1)
}

func (m *MockRepository) ListProducts(ctx context.Context, receptionID string) ([]*repository.Product, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*repository.Product), args.Error(1)
}

func (m *MockRepository) ListReception(ctx context.Context, PVZID string) ([]*repository.Reception, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*repository.Reception), args.Error(1)
}

func (m *MockRepository) ListUser(ctx context.Context) ([]*repository.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*repository.User), args.Error(1)
}

func (m *MockRepository) ExecTx(ctx context.Context, fn func(*sqlx.Tx) error) error {
	args := m.Called(ctx)
	return args.Error(1)
}

func TestService_IsValidCity(t *testing.T) {
	tests := []struct {
		name     string
		city     string
		config   *config.Config
		expected bool
	}{
		{
			name:     "valid city",
			city:     "Moscow",
			config:   &config.Config{Cities: []string{"Moscow", "Saint Petersburg"}},
			expected: true,
		},
		{
			name:     "invalid city",
			city:     "Paris",
			config:   &config.Config{Cities: []string{"Moscow", "Saint Petersburg"}},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{config: tt.config}
			assert.Equal(t, tt.expected, s.IsValidCity(tt.city))
		})
	}
}

func TestService_IsValidProductType(t *testing.T) {
	tests := []struct {
		name        string
		productType string
		config      *config.Config
		expected    bool
	}{
		{
			name:        "valid product type",
			productType: "electronics",
			config:      &config.Config{ProductTypes: []string{"electronics", "clothing"}},
			expected:    true,
		},
		{
			name:        "invalid product type",
			productType: "food",
			config:      &config.Config{ProductTypes: []string{"electronics", "clothing"}},
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{config: tt.config}
			assert.Equal(t, tt.expected, s.IsValidProductType(tt.productType))
		})
	}
}

func TestService_IsValidRole(t *testing.T) {
	tests := []struct {
		name     string
		role     UserRole
		expected bool
	}{
		{
			name:     "valid employee role",
			role:     UserRoleEmployee,
			expected: true,
		},
		{
			name:     "valid moderator role",
			role:     UserRoleModerator,
			expected: true,
		},
		{
			name:     "invalid role",
			role:     "admin",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{}
			assert.Equal(t, tt.expected, s.IsValidRole(tt.role))
		})
	}
}

func TestService_RegisterUser(t *testing.T) {
	tests := []struct {
		name       string
		email      string
		password   string
		role       string
		config     *config.Config
		mockSetup  func(*MockRepository, string)
		expectErr  bool
		errMessage string
	}{
		{
			name:     "successful registration",
			email:    "test@example.com",
			password: "password",
			role:     "employee",
			config:   &config.Config{},
			mockSetup: func(mr *MockRepository, hashedPassword string) {
				mr.On("CreateUser", mock.Anything, "test@example.com", mock.AnythingOfType("string"), "employee").
					Return(&repository.User{
						Email:    "test@example.com",
						Password: hashedPassword, // Возвращаем тот же хеш, что и получили
						Role:     "employee",
					}, nil)
			},
			expectErr: false,
		},
		{
			name:       "invalid role",
			email:      "test@example.com",
			password:   "password",
			role:       "invalid",
			config:     &config.Config{},
			mockSetup:  func(mr *MockRepository, _ string) {},
			expectErr:  true,
			errMessage: "invalid role: invalid",
		},
		{
			name:     "repository error",
			email:    "test@example.com",
			password: "password",
			role:     "employee",
			config:   &config.Config{},
			mockSetup: func(mr *MockRepository, _ string) {
				mr.On("CreateUser", mock.Anything, "test@example.com", mock.AnythingOfType("string"), "employee").
					Return(&repository.User{}, errors.New("repository error"))
			},
			expectErr:  true,
			errMessage: "repository error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockRepository{}

			// Генерируем хеш пароля один раз для теста
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(tt.password), bcrypt.DefaultCost)
			assert.NoError(t, err, "failed to hash password")

			tt.mockSetup(mockRepo, string(hashedPassword))

			s := NewService(mockRepo, tt.config)
			user, err := s.RegisterUser(context.Background(), tt.email, tt.password, tt.role)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMessage)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.email, user.Email)
				assert.Equal(t, tt.role, user.Role)

				// Проверяем что пароль действительно хеширован
				err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(tt.password))
				assert.NoError(t, err, "password should be hashed correctly")
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_CheckPassword(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)
	user := &repository.User{Password: string(hashedPassword)}

	tests := []struct {
		name      string
		user      *repository.User
		password  string
		expectErr bool
	}{
		{
			name:      "correct password",
			user:      user,
			password:  "correct",
			expectErr: false,
		},
		{
			name:      "incorrect password",
			user:      user,
			password:  "wrong",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{}
			err := s.CheckPassword(context.Background(), tt.user, tt.password)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_CreatePVZ(t *testing.T) {
	tests := []struct {
		name       string
		city       string
		config     *config.Config
		mockSetup  func(*MockRepository)
		expectErr  bool
		errMessage string
	}{
		{
			name:   "successful creation",
			city:   "Moscow",
			config: &config.Config{Cities: []string{"Moscow", "Saint Petersburg"}},
			mockSetup: func(mr *MockRepository) {
				mr.On("CreatePVZ", mock.Anything, "Moscow").
					Return(&repository.PVZ{}, nil)
			},
			expectErr: false,
		},
		{
			name:       "invalid city",
			city:       "Paris",
			config:     &config.Config{Cities: []string{"Moscow", "Saint Petersburg"}},
			mockSetup:  func(mr *MockRepository) {},
			expectErr:  true,
			errMessage: "invalid city: Paris",
		},
		{
			name:   "repository error",
			city:   "Moscow",
			config: &config.Config{Cities: []string{"Moscow", "Saint Petersburg"}},
			mockSetup: func(mr *MockRepository) {
				mr.On("CreatePVZ", mock.Anything, "Moscow").
					Return(&repository.PVZ{}, errors.New("repository error"))
			},
			expectErr:  true,
			errMessage: "repository error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockRepository{}
			tt.mockSetup(mockRepo)

			s := NewService(mockRepo, tt.config)
			_, err := s.CreatePVZ(context.Background(), tt.city)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMessage)
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_CreateProduct(t *testing.T) {
	tests := []struct {
		name        string
		receptionID string
		productType string
		config      *config.Config
		mockSetup   func(*MockRepository)
		expectErr   bool
		errMessage  string
	}{
		{
			name:        "successful creation",
			receptionID: "123",
			productType: "electronics",
			config:      &config.Config{ProductTypes: []string{"electronics", "clothing"}},
			mockSetup: func(mr *MockRepository) {
				mr.On("CreateProduct", mock.Anything, "123", "electronics").
					Return(&repository.Product{}, nil)
			},
			expectErr: false,
		},
		{
			name:        "invalid product type",
			receptionID: "123",
			productType: "food",
			config:      &config.Config{ProductTypes: []string{"electronics", "clothing"}},
			mockSetup:   func(mr *MockRepository) {},
			expectErr:   true,
			errMessage:  "invalid product type: food",
		},
		{
			name:        "repository error",
			receptionID: "123",
			productType: "electronics",
			config:      &config.Config{ProductTypes: []string{"electronics", "clothing"}},
			mockSetup: func(mr *MockRepository) {
				mr.On("CreateProduct", mock.Anything, "123", "electronics").
					Return(&repository.Product{}, errors.New("repository error"))
			},
			expectErr:  true,
			errMessage: "repository error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockRepository{}
			tt.mockSetup(mockRepo)

			s := NewService(mockRepo, tt.config)
			_, err := s.CreateProduct(context.Background(), tt.receptionID, tt.productType)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMessage)
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_GetUserByEmail(t *testing.T) {
	mockRepo := &MockRepository{}
	mockRepo.On("GetUserByEmail", mock.Anything, "test@example.com").
		Return(&repository.User{Email: "test@example.com"}, nil)

	s := NewService(mockRepo, &config.Config{})
	user, err := s.GetUserByEmail(context.Background(), "test@example.com")

	assert.NoError(t, err)
	assert.Equal(t, "test@example.com", user.Email)
	mockRepo.AssertExpectations(t)
}

func TestService_CloseReception(t *testing.T) {
	mockRepo := &MockRepository{}
	mockRepo.On("CloseReception", mock.Anything, "123").
		Return(&repository.Reception{ID: "123", Status: "close"}, nil)

	s := NewService(mockRepo, &config.Config{})
	reception, err := s.CloseReception(context.Background(), "123")

	assert.NoError(t, err)
	assert.Equal(t, "close", reception.Status)
	mockRepo.AssertExpectations(t)
}

func TestService_DeleteProduct(t *testing.T) {
	mockRepo := &MockRepository{}
	mockRepo.On("DeleteProduct", mock.Anything, "123").
		Return(&repository.Product{ID: "123"}, nil)

	s := NewService(mockRepo, &config.Config{})
	product, err := s.DeleteProduct(context.Background(), "123")

	assert.NoError(t, err)
	assert.Equal(t, "123", product.ID)
	mockRepo.AssertExpectations(t)
}

func TestService_CreateReception(t *testing.T) {
	mockRepo := &MockRepository{}
	mockRepo.On("CreateReception", mock.Anything, "123").
		Return(&repository.Reception{ID: "456", PVZID: "123"}, nil)

	s := NewService(mockRepo, &config.Config{})
	reception, err := s.CreateReception(context.Background(), "123")

	assert.NoError(t, err)
	assert.Equal(t, "123", reception.PVZID)
	mockRepo.AssertExpectations(t)
}

func TestService_ListPVZ(t *testing.T) {
	now := time.Now()
	startDate := now.AddDate(0, -1, 0)
	endDate := now

	expectedPVZs := []*repository.PVZWithReceptions{
		{
			PVZ: &repository.PVZ{ID: "1", City: "Moscow"},
		},
	}

	mockRepo := &MockRepository{}
	mockRepo.On("ListPVZ", mock.Anything, &startDate, &endDate, 1, 10).
		Return(expectedPVZs, nil)

	s := NewService(mockRepo, &config.Config{})
	pvzs, err := s.ListPVZ(context.Background(), &startDate, &endDate, 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, expectedPVZs, pvzs)
	mockRepo.AssertExpectations(t)
}

func TestService_ListAllPVZ(t *testing.T) {
	expectedPVZs := []*repository.PVZ{
		{ID: "1", City: "Moscow"},
		{ID: "2", City: "Saint Petersburg"},
	}

	mockRepo := &MockRepository{}
	mockRepo.On("ListAllPVZ", mock.Anything).
		Return(expectedPVZs, nil)

	s := NewService(mockRepo, &config.Config{})
	pvzs, err := s.ListAllPVZ(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, expectedPVZs, pvzs)
	mockRepo.AssertExpectations(t)
}
