package control

import (
	"github.com/gofiber/fiber/v2"
)

type Handlers interface {
	GetOrder() fiber.Handler
}