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

func sendResponse(c *fiber.Ctx, statusCode int, message string, data *fiber.Map) error {
	response := fiber.Map{
		"message": nil,
		"data":    data,
	}
	if message != "" {
		response["message"] = message
	}
	return c.Status(statusCode).JSON(response)
}
