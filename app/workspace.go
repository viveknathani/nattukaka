package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/viveknathani/nattukaka/types"
)

// CreateWorkspaceController ...
func (app *App) CreateWorkspaceController(c *fiber.Ctx) error {

	// Parse body
	var req types.CreateWorkspaceRequest
	if err := c.BodyParser(&req); err != nil {
		return sendResponse(c, fiber.StatusBadRequest, "cannot parse json", nil)
	}

	// Extract email from token
	email := c.Locals("user").(jwt.MapClaims)["email"].(string)

	code, message := app.Service.CreateWorkspace(email, req.Name)

	return sendResponse(c, code, message, nil)
}

// GetWorkspaceUsersController ...
func (app *App) GetWorkspaceUsersController(c *fiber.Ctx) error {

	workspacePublicID := c.Params("workspaceID")
	if workspacePublicID == "" {
		return sendResponse(c, fiber.StatusBadRequest, "workspace ID is required", nil)
	}

	code, message, data := app.Service.GetWorkspaceUsers(workspacePublicID)

	return sendResponse(c, code, message, data)
}

// DeleteWorkspaceController ...
func (app *App) DeleteWorkspaceController(c *fiber.Ctx) error {

	workspacePublicID := c.Params("workspaceID")
	if workspacePublicID == "" {
		return sendResponse(c, fiber.StatusBadRequest, "workspace ID is required", nil)
	}

	// Extract email from token
	email := c.Locals("user").(jwt.MapClaims)["email"].(string)

	code, message := app.Service.DeleteWorkspace(email, workspacePublicID)

	return sendResponse(c, code, message, nil)
}

// GetMyWorkspacesController ...
func (app *App) GetMyWorkspacesController(c *fiber.Ctx) error {

	// Extract email from token
	email := c.Locals("user").(jwt.MapClaims)["email"].(string)

	code, message, data := app.Service.GetMyWorkspaces(email)

	return sendResponse(c, code, message, data)
}
