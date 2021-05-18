package service

import (
	"bank-system-go/internal/model"
	"context"
)

type TransationService interface {
	CreateTransation(ctx context.Context, transation model.Transation) (model.Transation, error)
}
