package http

import (
	"L0/internal/control"

	"github.com/gofiber/fiber/v2"
)

func MapAPIRoutes(group fiber.Router, h control.Handlers) {
	group.Get("/get_order", h.GetOrder())
}
