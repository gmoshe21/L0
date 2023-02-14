package control

import (
	"L0/internal/models"
	"context"
)

type UseCase interface {
	NewOrder(ctx context.Context, params models.Order) error
	DataRecovery(ctx context.Context) error
	GetOrder(ctx context.Context, uid string) (result []byte, err error)
}
