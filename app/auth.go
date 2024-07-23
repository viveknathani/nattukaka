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
		return sendResponse(c, fiber.StatusBadRequest, "cannot parse json", nil)
	}

	code, message := app.Service.Signup(&req)

	return sendResponse(c, code, message, nil)
}

// SendOTPController ...
func (app *App) SendOTPController(c *fiber.Ctx) error {

	// Parse body
	var req types.SendOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return sendResponse(c, fiber.StatusBadRequest, "cannot parse json", nil)
	}

	code, message := app.Service.SendOTP(&req)

	return sendResponse(c, code, message, nil)
}
