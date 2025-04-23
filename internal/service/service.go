package service

import (
	"context"
	"fmt"
	"time"

	"github.com/DarRo9/pvz_service/config"
	"github.com/DarRo9/pvz_service/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type ReceptionStatus string

const (
	Close      ReceptionStatus = "close"
	InProgress ReceptionStatus = "in_progress"
)

type UserRole string

const (
	UserRoleEmployee  UserRole = "employee"
	UserRoleModerator UserRole = "moderator"
)

type ServiceInterface interface {
	RegisterUser(ctx context.Context, email string, password string, role string) (*repository.User, error)

	CheckPassword(ctx context.Context, user *repository.User, password string) error

	GetUserByEmail(ctx context.Context, email string) (*repository.User, error)

	CreatePVZ(ctx context.Context, city string) (*repository.PVZ, error)

	CloseReception(ctx context.Context, pvzId string) (*repository.Reception, error)

	DeleteProduct(ctx context.Context, pvzId string) (*repository.Product, error)

	CreateReception(ctx context.Context, pvzId string) (*repository.Reception, error)

	ListPVZ(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]*repository.PVZWithReceptions, error)

	CreateProduct(ctx context.Context, receptionID string, productType string) (*repository.Product, error)

	ListAllPVZ(ctx context.Context) ([]*repository.PVZ, error)

	IsValidCity(city string) bool

	IsValidProductType(productType string) bool

	IsValidRole(role UserRole) bool
}

type Service struct {
	repo   repository.Repository
	config *config.Config
}

func NewService(repo repository.Repository, config *config.Config) *Service {
	return &Service{
		repo:   repo,
		config: config,
	}
}

func (s *Service) IsValidCity(city string) bool {
	for _, c := range s.config.Cities {
		if c == city {
			return true
		}
	}
	return false
}

func (s *Service) IsValidProductType(productType string) bool {
	for _, p := range s.config.ProductTypes {
		if p == productType {
			return true
		}
	}
	return false
}

func (s *Service) IsValidRole(role UserRole) bool {
	switch role {
	case UserRoleEmployee, UserRoleModerator:
		return true
	default:
		return false
	}
}

func (s *Service) RegisterUser(ctx context.Context, email string, password string, role string) (*repository.User, error) {
	if !s.IsValidRole(UserRole(role)) {
		return nil, fmt.Errorf("invalid role: %s", role)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}

	user, err := s.repo.CreateUser(ctx, email, string(hashedPassword), role)
	return user, err
}

func (s *Service) CheckPassword(ctx context.Context, user *repository.User, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*repository.User, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	return user, err
}

func (s *Service) CreatePVZ(ctx context.Context, city string) (*repository.PVZ, error) {
	if !s.IsValidCity(city) {
		return nil, fmt.Errorf("invalid city: %s", city)
	}
	pvz, err := s.repo.CreatePVZ(ctx, city)
	return pvz, err
}

func (s *Service) CloseReception(ctx context.Context, pvzId string) (*repository.Reception, error) {
	rc, err := s.repo.CloseReception(ctx, pvzId)
	return rc, err
}

func (s *Service) DeleteProduct(ctx context.Context, pvzId string) (*repository.Product, error) {
	product, err := s.repo.DeleteProduct(ctx, pvzId)
	return product, err
}

func (s *Service) CreateReception(ctx context.Context, pvzId string) (*repository.Reception, error) {
	rc, err := s.repo.CreateReception(ctx, pvzId)
	return rc, err
}

func (s *Service) ListPVZ(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]*repository.PVZWithReceptions, error) {
	pvzs, err := s.repo.ListPVZ(ctx, startDate, endDate, page, limit)

	return pvzs, err
}

func (s *Service) CreateProduct(ctx context.Context, receptionID string, productType string) (*repository.Product, error) {
	if !s.IsValidProductType(productType) {
		return nil, fmt.Errorf("invalid product type: %s", productType)
	}

	product, err := s.repo.CreateProduct(ctx, receptionID, productType)

	return product, err
}

func (s *Service) ListAllPVZ(ctx context.Context) ([]*repository.PVZ, error) {
	pvzs, err := s.repo.ListAllPVZ(ctx)
	return pvzs, err
}
