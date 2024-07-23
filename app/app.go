package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/viveknathani/nattukaka/service"
)

// App encapsulates the entire service
type App struct {
	Fiber   *fiber.App
	Service *service.Service
}
