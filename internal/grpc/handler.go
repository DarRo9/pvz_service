package grpc

import (
	"context"
	"log"

	"github.com/DarRo9/pvz_service/internal/grpc/pvz/pvz_v1"
	"github.com/DarRo9/pvz_service/internal/service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCHandler struct {
	service *service.Service
	pvz_v1.UnimplementedPVZServiceServer
}

func NewGRPCHandler(service *service.Service) *GRPCHandler {
	return &GRPCHandler{
		service: service,
	}
}

func (h *GRPCHandler) GetPVZList(ctx context.Context, _ *pvz_v1.GetPVZListRequest) (*pvz_v1.GetPVZListResponse, error) {
	log.Println("Got request in GetPVZList")
	pvzs, err := h.service.ListAllPVZ(ctx)
	if err != nil {
		log.Printf("Error getting PVZ list: %v", err)
		return nil, err
	}

	response := &pvz_v1.GetPVZListResponse{}
	for _, p := range pvzs {
		response.Pvzs = append(response.Pvzs, &pvz_v1.PVZ{
			Id:               p.ID,
			RegistrationDate: timestamppb.New(p.RegistrationDate),
			City:             p.City,
		})
	}
	return response, nil
}
