package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/viveknathani/nattukaka/types"
)

// SignupController ...
func (app *App) SignupController(c *fiber.Ctx) error {

	// Parse body
	var req types.SignupRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}

	code, message := app.Service.Signup(&req)

	return sendResponse(c, code, message, nil)
}
