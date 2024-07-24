package app

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/viveknathani/nattukaka/types"
)

// CreateServiceController handles the creation of a new service
func (app *App) CreateServiceController(c *fiber.Ctx) error {
	// Parse body
	var req types.CreateServiceRequest
	if err := c.BodyParser(&req); err != nil {
		return sendResponse(c, fiber.StatusBadRequest, "cannot parse json", nil)
	}

	// Extract email from token
	email := c.Locals("user").(jwt.MapClaims)["email"].(string)

	// Create service
	code, message := app.Service.CreateService(email, &req)

	return sendResponse(c, code, message, nil)
}

// GetServiceByIDController handles fetching a service by its public ID
func (app *App) GetServiceByIDController(c *fiber.Ctx) error {
	publicID := c.Params("serviceID")

	code, message, data := app.Service.GetServiceByID(publicID)
	return sendResponse(c, code, message, data)
}

// GetServicesByWorkspaceController handles fetching all services in a workspace, paginated
func (app *App) GetServicesByWorkspaceController(c *fiber.Ctx) error {
	workspaceID := c.Params("workspaceID")
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		return sendResponse(c, fiber.StatusInternalServerError, "something went wrong", nil)
	}

	code, message, data := app.Service.GetServicesByWorkspace(workspaceID, page)
	return sendResponse(c, code, message, data)
}

// UpdateServiceController handles updating a service
func (app *App) UpdateServiceController(c *fiber.Ctx) error {
	// Parse body
	var req types.UpdateServiceRequest
	if err := c.BodyParser(&req); err != nil {
		return sendResponse(c, fiber.StatusBadRequest, "cannot parse json", nil)
	}

	// Update service
	code, message := app.Service.UpdateService(c.Params("serviceID"), &req)
	return sendResponse(c, code, message, nil)
}

// DeleteServiceController handles deleting a service
func (app *App) DeleteServiceController(c *fiber.Ctx) error {
	publicID := c.Params("serviceID")

	err := app.Service.DeleteService(publicID)
	if err != nil {
		return sendResponse(c, fiber.StatusInternalServerError, "something went wrong", nil)
	}

	return sendResponse(c, fiber.StatusOK, "service deleted", nil)
}
