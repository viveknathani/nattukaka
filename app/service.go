package app

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/viveknathani/nattukaka/types"
)

// CreateServiceController handles the creation of a new service
func (app *App) CreateServiceController(c *fiber.Ctx) error {
	// Parse body
	var req types.CreateServiceRequest
	if err := c.BodyParser(&req); err != nil {
		return sendResponse(c, fiber.StatusBadRequest, "cannot parse json", nil)
	}

	// Basic validation
	if req.Base == nil || (req.Web == nil && req.Database == nil) {
		return sendResponse(c, fiber.StatusBadRequest, "invalid request body", nil)
	}

	// Prepare service data
	service := &types.Service{
		PublicID:       req.Base.PublicID,
		Name:           req.Base.Name,
		Status:         req.Base.Status,
		Type:           req.Base.Type,
		Runtime:        req.Base.Runtime,
		WorkspaceID:    req.Base.WorkspaceID,
		CreatedBy:      req.Base.CreatedBy,
		InstanceTypeID: req.Base.InstanceTypeID,
		InternalURL:    req.Base.InternalURL,
		ExternalURL:    req.Base.ExternalURL,
	}

	var webService *types.WebService
	var databaseService *types.DatabaseService

	if req.Base.Type == "WEB" {
		webService = &types.WebService{
			PublicID:         req.Web.PublicID,
			Repository:       req.Web.Repository,
			Branch:           req.Web.Branch,
			RootDirectory:    req.Web.RootDirectory,
			BuildCommand:     req.Web.BuildCommand,
			PreDeployCommand: req.Web.PreDeployCommand,
			StartCommand:     req.Web.StartCommand,
			HealthCheckPath:  req.Web.HealthCheckPath,
			Environment:      req.Web.Environment,
		}
	} else if req.Base.Type == "DATABASE" {
		databaseService = &types.DatabaseService{
			PublicID: req.Database.PublicID,
		}
	}

	// Create service
	err := app.Service.CreateService(service, webService, databaseService)
	if err != nil {
		return sendResponse(c, fiber.StatusInternalServerError, "something went wrong", nil)
	}

	return sendResponse(c, fiber.StatusCreated, "service created", nil)
}

// GetServiceByIDController handles fetching a service by its public ID
func (app *App) GetServiceByIDController(c *fiber.Ctx) error {
	publicID := c.Params("serviceID")

	service, err := app.Service.GetServiceByID(publicID)
	if err != nil {
		return sendResponse(c, fiber.StatusInternalServerError, "something went wrong", nil)
	}

	return sendResponse(c, fiber.StatusOK, "", service)
}

// GetServicesByWorkspaceController handles fetching all services in a workspace, paginated
func (app *App) GetServicesByWorkspaceController(c *fiber.Ctx) error {
	workspaceID := c.Params("workspaceID")
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		return sendResponse(c, fiber.StatusBadRequest, "invalid page number", nil)
	}

	workspaceIDInt, err := strconv.Atoi(workspaceID)
	if err != nil {
		return sendResponse(c, fiber.StatusBadRequest, "invalid workspace ID", nil)
	}

	services, err := app.Service.GetServicesByWorkspace(workspaceIDInt, page)
	if err != nil {
		return sendResponse(c, fiber.StatusInternalServerError, "something went wrong", nil)
	}

	return sendResponse(c, fiber.StatusOK, "", services)
}

// UpdateServiceController handles updating a service
func (app *App) UpdateServiceController(c *fiber.Ctx) error {
	// Parse body
	var req types.UpdateServiceRequest
	if err := c.BodyParser(&req); err != nil {
		return sendResponse(c, fiber.StatusBadRequest, "cannot parse json", nil)
	}

	publicID := c.Params("serviceID")

	// Fetch the existing service
	service, err := app.Service.GetServiceByID(publicID)
	if err != nil {
		return sendResponse(c, fiber.StatusInternalServerError, "something went wrong", nil)
	}

	// Update base properties
	if req.Base != nil {
		if req.Base.Name != "" {
			service.Name = req.Base.Name
		}
	}

	var webService *types.WebService
	var databaseService *types.DatabaseService

	if service.Type == "WEB" && req.Web != nil {
		webService = &types.WebService{
			PublicID:         service.PublicID,
			Repository:       req.Web.Repository,
			Branch:           req.Web.Branch,
			RootDirectory:    req.Web.RootDirectory,
			BuildCommand:     req.Web.BuildCommand,
			PreDeployCommand: req.Web.PreDeployCommand,
			StartCommand:     req.Web.StartCommand,
			HealthCheckPath:  req.Web.HealthCheckPath,
			Environment:      req.Web.Environment,
		}
	} else if service.Type == "DATABASE" && req.Database != nil {
		databaseService = &types.DatabaseService{
			PublicID: service.PublicID,
		}
	}

	// Update service
	err = app.Service.UpdateService(service, webService, databaseService)
	if err != nil {
		return sendResponse(c, fiber.StatusInternalServerError, "something went wrong", nil)
	}

	return sendResponse(c, fiber.StatusOK, "service updated", nil)
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
