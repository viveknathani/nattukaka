package app

import (
	"github.com/gofiber/fiber/v2"
)

// GetAllInstanceTypesController handles the request to fetch all instance types
func (app *App) GetAllInstanceTypesController(c *fiber.Ctx) error {
	instanceTypes, err := app.Service.GetAllInstanceTypes()
	if err != nil {
		return sendResponse(c, fiber.StatusInternalServerError, "something went wrong", nil)
	}

	return sendResponse(c, fiber.StatusOK, "", instanceTypes)
}
