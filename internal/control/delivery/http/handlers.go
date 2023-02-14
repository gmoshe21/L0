package http

import (
	"L0/config"
	"L0/internal/control"
	"context"

	"github.com/gofiber/fiber/v2"
)

type controlHandlers struct {//L0
	cfg       *config.Config
	controlUC control.UseCase
}

func NewControlHandlers(cfg *config.Config, controlUC control.UseCase) control.Handlers {
	return &controlHandlers{cfg: cfg, controlUC: controlUC}
}

func (ctrl *controlHandlers) GetOrder() fiber.Handler {
	return func(c *fiber.Ctx) error {
		uid := c.OriginalURL()
		uid = uid[24:]

		if len(uid) == 0 {
			return c.SendStatus(fiber.StatusNoContent)
		}

		result, err := ctrl.controlUC.GetOrder(context.Background(), uid) 
		if err != nil {
			return err
		}

		return c.Send(result)
	}
}