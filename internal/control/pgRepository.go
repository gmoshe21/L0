package control

import (
	"L0/internal/models"
	"context"
)

type Repository interface {
	NewOrder(ctx context.Context, params models.Order) error
	DataRecovery(ctx context.Context) (result []byte, err error)
	GetOrder(ctx context.Context, uid string) (result []byte, err error)
}
