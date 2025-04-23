package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DarRo9/pvz_service/internal/repository"
	"github.com/DarRo9/pvz_service/internal/service"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) IsValidRole(role service.UserRole) bool {
	args := m.Called(role)
	return args.Bool(0)
}

func (s *MockService) IsValidCity(city string) bool {
	return true
}

func (s *MockService) IsValidProductType(city string) bool {
	return true
}

func (m *MockService) GetUserByEmail(ctx context.Context, email string) (*repository.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*repository.User), args.Error(1)
}

func (m *MockService) CheckPassword(ctx context.Context, user *repository.User, password string) error {
	args := m.Called(ctx, user, password)
	return args.Error(0)
}

func (m *MockService) CreateProduct(ctx context.Context, receptionID, productType string) (*repository.Product, error) {
	args := m.Called(ctx, receptionID, productType)
	return args.Get(0).(*repository.Product), args.Error(1)
}

func (m *MockService) ListPVZ(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]*repository.PVZWithReceptions, error) {
	args := m.Called(ctx, startDate, endDate, page, limit)
	return args.Get(0).([]*repository.PVZWithReceptions), args.Error(1)
}

func (m *MockService) CreatePVZ(ctx context.Context, city string) (*repository.PVZ, error) {
	args := m.Called(ctx, city)
	return args.Get(0).(*repository.PVZ), args.Error(1)
}

func (m *MockService) CloseReception(ctx context.Context, pvzID string) (*repository.Reception, error) {
	args := m.Called(ctx, pvzID)
	return args.Get(0).(*repository.Reception), args.Error(1)
}

func (m *MockService) DeleteProduct(ctx context.Context, pvzID string) (*repository.Product, error) {
	args := m.Called(ctx, pvzID)
	return args.Get(0).(*repository.Product), args.Error(1)
}

func (m *MockService) CreateReception(ctx context.Context, pvzID string) (*repository.Reception, error) {
	args := m.Called(ctx, pvzID)
	return args.Get(0).(*repository.Reception), args.Error(1)
}

func (m *MockService) RegisterUser(ctx context.Context, email, password, role string) (*repository.User, error) {
	args := m.Called(ctx, email, password, role)
	return args.Get(0).(*repository.User), args.Error(1)
}

func (m *MockService) ListAllPVZ(ctx context.Context) ([]*repository.PVZ, error) {
	return nil, nil
}

func TestHTTPHandler_PostDummyLogin(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    PostDummyLoginJSONBody
		mockSetup      func(*MockService)
		expectedStatus int
		expectToken    bool
	}{
		{
			name: "successful login with valid role",
			requestBody: PostDummyLoginJSONBody{
				Role: "employee",
			},
			mockSetup: func(ms *MockService) {
				ms.On("IsValidRole", service.UserRole("employee")).Return(true)
			},
			expectedStatus: http.StatusOK,
			expectToken:    true,
		},
		{
			name: "invalid role",
			requestBody: PostDummyLoginJSONBody{
				Role: "invalid",
			},
			mockSetup: func(ms *MockService) {
				ms.On("IsValidRole", service.UserRole("invalid")).Return(false)
			},
			expectedStatus: http.StatusBadRequest,
			expectToken:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockService)
			tt.mockSetup(mockService)
			handler := NewHTTPHandler(mockService)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/dummy_login", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			handler.PostDummyLogin(w, req)

			resp := w.Result()
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.expectToken {
				var tokenResp Token
				err := json.NewDecoder(resp.Body).Decode(&tokenResp)
				assert.NoError(t, err)
				assert.NotEmpty(t, tokenResp)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestHTTPHandler_PostLogin(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    PostLoginJSONBody
		mockSetup      func(*MockService)
		expectedStatus int
		expectToken    bool
	}{
		{
			name: "successful login",
			requestBody: PostLoginJSONBody{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(ms *MockService) {
				user := &repository.User{
					ID:    "user123",
					Email: "test@example.com",
					Role:  "employee",
				}
				ms.On("GetUserByEmail", mock.Anything, "test@example.com").Return(user, nil)
				ms.On("CheckPassword", mock.Anything, user, "password123").Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectToken:    true,
		},
		{
			name: "invalid password",
			requestBody: PostLoginJSONBody{
				Email:    "test@example.com",
				Password: "wrong",
			},
			mockSetup: func(ms *MockService) {
				user := &repository.User{
					ID:    "user123",
					Email: "test@example.com",
					Role:  "employee",
				}
				ms.On("GetUserByEmail", mock.Anything, "test@example.com").Return(user, nil)
				ms.On("CheckPassword", mock.Anything, user, "wrong").Return(errors.New("invalid password"))
			},
			expectedStatus: http.StatusUnauthorized,
			expectToken:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockService)
			tt.mockSetup(mockService)
			handler := NewHTTPHandler(mockService)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			handler.PostLogin(w, req)

			resp := w.Result()
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.expectToken {
				var tokenResp Token
				err := json.NewDecoder(resp.Body).Decode(&tokenResp)
				assert.NoError(t, err)
				assert.NotEmpty(t, tokenResp)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestHTTPHandler_PostProducts(t *testing.T) {
	UUID := uuid.New()
	tests := []struct {
		name           string
		requestBody    Product
		mockSetup      func(*MockService)
		expectedStatus int
		withAuth       bool
	}{
		{
			name: "successful product creation",
			requestBody: Product{
				ReceptionId: UUID,
				Type:        "electronics",
			},
			mockSetup: func(ms *MockService) {
				product := &repository.Product{
					ID:          "product123",
					ReceptionId: UUID.String(),
					Type:        "electronics",
				}
				ms.On("CreateProduct", mock.Anything, UUID.String(), "electronics").Return(product, nil)
			},
			expectedStatus: http.StatusCreated,
			withAuth:       true,
		},
		{
			name: "unauthorized access",
			requestBody: Product{
				ReceptionId: UUID,
				Type:        "electronics",
			},
			mockSetup:      func(ms *MockService) {},
			expectedStatus: http.StatusUnauthorized,
			withAuth:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockService)
			tt.mockSetup(mockService)
			handler := NewHTTPHandler(mockService)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/products", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			if tt.withAuth {
				claims := jwt.MapClaims{"role": "employee"}
				ctx := context.WithValue(req.Context(), "user", claims)
				req = req.WithContext(ctx)
			}

			handler.PostProducts(w, req)

			resp := w.Result()
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.expectedStatus == http.StatusCreated {
				var productResp Product
				err := json.NewDecoder(resp.Body).Decode(&productResp)
				assert.NoError(t, err)
				assert.Equal(t, tt.requestBody.Type, productResp.Type)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestHTTPHandler_GetPvz(t *testing.T) {
	page_1 := 1
	limit_10 := 10

	page_0 := 0
	tests := []struct {
		name           string
		params         GetPvzParams
		mockSetup      func(*MockService)
		expectedStatus int
		withAuth       bool
	}{
		{
			name: "successful list with pagination",
			params: GetPvzParams{
				Page:  &page_1,
				Limit: &limit_10,
			},
			mockSetup: func(ms *MockService) {
				pvzs := []*repository.PVZWithReceptions{
					{
						PVZ: &repository.PVZ{
							ID:   "pvz1",
							City: "Moscow",
						},
						Receptions: []*repository.ReceptionWithProducts{},
					},
				}
				ms.On("ListPVZ", mock.Anything, (*time.Time)(nil), (*time.Time)(nil), 1, 10).Return(pvzs, nil)
			},
			expectedStatus: http.StatusOK,
			withAuth:       true,
		},
		{
			name: "invalid page parameter",
			params: GetPvzParams{
				Page:  &page_0,
				Limit: &limit_10,
			},
			mockSetup:      func(ms *MockService) {},
			expectedStatus: http.StatusBadRequest,
			withAuth:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockService)
			tt.mockSetup(mockService)
			handler := NewHTTPHandler(mockService)

			req := httptest.NewRequest("GET", "/pvz", nil)
			w := httptest.NewRecorder()

			// Add query parameters
			q := req.URL.Query()
			if tt.params.Page != nil {
				q.Add("page", fmt.Sprintf("%d", *tt.params.Page))
			}
			if tt.params.Limit != nil {
				q.Add("limit", fmt.Sprintf("%d", *tt.params.Limit))
			}
			req.URL.RawQuery = q.Encode()

			if tt.withAuth {
				claims := jwt.MapClaims{"role": "employee"}
				ctx := context.WithValue(req.Context(), "user", claims)
				req = req.WithContext(ctx)
			}

			handler.GetPvz(w, req, tt.params)

			resp := w.Result()
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.expectedStatus == http.StatusOK {
				var pvzs []*PVZWithReceptions
				err := json.NewDecoder(resp.Body).Decode(&pvzs)
				assert.NoError(t, err)
				assert.NotEmpty(t, pvzs)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestHTTPHandler_PostPvz(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    PVZ
		mockSetup      func(*MockService)
		expectedStatus int
		withAuth       bool
	}{
		{
			name: "successful PVZ creation",
			requestBody: PVZ{
				City: "Moscow",
			},
			mockSetup: func(ms *MockService) {
				pvz := &repository.PVZ{
					ID:   "pvz123",
					City: "Moscow",
				}
				ms.On("CreatePVZ", mock.Anything, "Moscow").Return(pvz, nil)
			},
			expectedStatus: http.StatusCreated,
			withAuth:       true,
		},
		{
			name: "forbidden for non-moderator",
			requestBody: PVZ{
				City: "Moscow",
			},
			mockSetup:      func(ms *MockService) {},
			expectedStatus: http.StatusForbidden,
			withAuth:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockService)
			tt.mockSetup(mockService)
			handler := NewHTTPHandler(mockService)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/pvz", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			if tt.withAuth {
				role := "employee"
				if tt.expectedStatus == http.StatusCreated {
					role = "moderator"
				}
				claims := jwt.MapClaims{"role": role}
				ctx := context.WithValue(req.Context(), "user", claims)
				req = req.WithContext(ctx)
			}

			handler.PostPvz(w, req)

			resp := w.Result()
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.expectedStatus == http.StatusCreated {
				var pvzResp PVZ
				err := json.NewDecoder(resp.Body).Decode(&pvzResp)
				assert.NoError(t, err)
				assert.Equal(t, tt.requestBody.City, pvzResp.City)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestHTTPHandler_PostRegister(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    PostRegisterJSONBody
		mockSetup      func(*MockService)
		expectedStatus int
	}{
		{
			name: "successful registration",
			requestBody: PostRegisterJSONBody{
				Email:    "new@example.com",
				Password: "password123",
				Role:     "employee",
			},
			mockSetup: func(ms *MockService) {
				user := &repository.User{
					ID:    "user123",
					Email: "new@example.com",
					Role:  "employee",
				}
				ms.On("RegisterUser", mock.Anything, "new@example.com", "password123", "employee").Return(user, nil)
			},
			expectedStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockService)
			tt.mockSetup(mockService)
			handler := NewHTTPHandler(mockService)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			handler.PostRegister(w, req)

			resp := w.Result()
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var userResp User
			err := json.NewDecoder(resp.Body).Decode(&userResp)
			assert.NoError(t, err)
			assert.Equal(t, tt.requestBody.Email, userResp.Email)

			mockService.AssertExpectations(t)
		})
	}
}

func TestValidateRole(t *testing.T) {
	tests := []struct {
		name           string
		ctx            context.Context
		allowedRoles   []string
		expectedResult bool
		expectedStatus int
	}{
		{
			name: "valid role",
			ctx: context.WithValue(context.Background(), "user", jwt.MapClaims{
				"role": "employee",
			}),
			allowedRoles:   []string{"employee", "moderator"},
			expectedResult: true,
			expectedStatus: 0,
		},
		{
			name: "invalid role",
			ctx: context.WithValue(context.Background(), "user", jwt.MapClaims{
				"role": "guest",
			}),
			allowedRoles:   []string{"employee", "moderator"},
			expectedResult: false,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "no user in context",
			ctx:            context.Background(),
			allowedRoles:   []string{"employee"},
			expectedResult: false,
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			result := validateRole(tt.ctx, w, tt.allowedRoles)
			assert.Equal(t, tt.expectedResult, result)

			if tt.expectedStatus != 0 {
				resp := w.Result()
				assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			}
		})
	}
}
