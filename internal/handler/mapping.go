package handler

import (
	"github.com/DarRo9/pvz_service/internal/repository"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func productRepositoryToHTTP(product *repository.Product) *Product {
	id, _ := uuid.Parse(product.ID)
	receptionId, _ := uuid.Parse(product.ReceptionId)
	return &Product{
		DateTime:    &product.ReceptionDate,
		Id:          &id,
		ReceptionId: receptionId,
		Type:        ProductType(product.Type),
	}
}

func pvzRepositoryToHTTP(pvz *repository.PVZ) *PVZ {
	id, _ := uuid.Parse(pvz.ID)
	return &PVZ{
		Id:   &id,
		City: PVZCity(pvz.City),
	}
}

func receptionRepositoryToHTTP(reception *repository.Reception) *Reception {
	id, _ := uuid.Parse(reception.ID)
	pvzId, _ := uuid.Parse(reception.PVZID)
	return &Reception{
		Id:       &id,
		PvzId:    pvzId,
		DateTime: reception.ExecutionDate,
		Status:   ReceptionStatus(reception.Status),
	}
}

func pvzWithReceptionsRepositoryToHTTP(pvz *repository.PVZWithReceptions) *PVZWithReceptions {
	receptions := make([]*ReceptionWithProducts, len(pvz.Receptions))
	for i, reception := range pvz.Receptions {
		receptions[i] = &ReceptionWithProducts{
			Reception: receptionRepositoryToHTTP(reception.Reception),
			Products:  make([]*Product, len(reception.Products)),
		}
		for j, product := range reception.Products {
			receptions[i].Products[j] = productRepositoryToHTTP(product)
		}
	}
	return &PVZWithReceptions{
		PVZ:        pvzRepositoryToHTTP(pvz.PVZ),
		Receptions: receptions,
	}
}

func userRepositoryToHTTP(user *repository.User) *User {
	id, _ := uuid.Parse(user.ID)
	return &User{
		Id:    &id,
		Email: openapi_types.Email(user.Email),
		Role:  UserRole(user.Role),
	}
}
