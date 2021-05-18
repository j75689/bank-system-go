package v1

import (
	"bank-system-go/internal/model"
	"bank-system-go/internal/repository"
	"bank-system-go/internal/service"
	"context"
)

var _ service.TransationService = (*TransationService)(nil)

func NewTransationService(transationRepo repository.TransationRepository) service.TransationService {
	return &TransationService{
		transationRepo: transationRepo,
	}
}

type TransationService struct {
	transationRepo repository.TransationRepository
}

func (svc *TransationService) CreateTransation(ctx context.Context, transation model.Transation) (model.Transation, error) {
	return transation, svc.transationRepo.CreateTransation(ctx, &transation)
}
