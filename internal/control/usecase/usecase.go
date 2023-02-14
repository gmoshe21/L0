package usecase

import (
	"L0/config"
	"L0/internal/control"
	"L0/internal/models"
	"sync"

	"context"
	"encoding/json"
)

type controlUC struct {
	cfg						*config.Config
	mu                      *sync.RWMutex
	controlRepo				control.Repository
	data					[]models.Order
}

func NewControlUseCase( cfg *config.Config, controlRepo control.Repository) control.UseCase {
	return &controlUC{cfg: cfg, mu: &sync.RWMutex{},controlRepo: controlRepo}
}

func (c *controlUC) NewOrder(ctx context.Context, params models.Order) error {
	c.mu.Lock()
	c.data = append(c.data, params)
	c.mu.Unlock()
	return c.controlRepo.NewOrder(ctx, params)
}

func (c *controlUC) DataRecovery(ctx context.Context) error {
	allData, err := c.controlRepo.DataRecovery(ctx)
	if err != nil {
		return err
	}

	if len(allData) == 0 {
		return nil
	}
	
	c.mu.Lock()
	err = json.Unmarshal(allData, &c.data)
	if err != nil {
		return err
	}
	c.mu.Unlock()

	return nil
}

func (c *controlUC) GetOrder(ctx context.Context, uid string) (result []byte, err error) {
	return c.controlRepo.GetOrder(ctx, uid)
}