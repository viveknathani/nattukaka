package app

import (
	"github.com/gofiber/fiber/v2"
)

// CreateVolumeController handles creating a new volume
func (app *App) CreateVolumeController(c *fiber.Ctx) error {
	code, message, data := app.Service.CreateVolume()
	return sendResponse(c, code, message, data)
}

// AttachVolumeToServiceController handles attaching a volume to a service
func (app *App) AttachVolumeToServiceController(c *fiber.Ctx) error {
	serviceID := c.Params("serviceID")
	volumeID := c.Params("volumeID")

	err := app.Service.AttachVolumeToService(serviceID, volumeID)
	if err != nil {
		app.Service.Logger.Error(err.Error())
		return sendResponse(c, fiber.StatusInternalServerError, "something went wrong", nil)
	}

	return sendResponse(c, fiber.StatusOK, "volume attached to service", nil)
}

// DetachVolumeFromServiceController handles detaching a volume from a service
func (app *App) DetachVolumeFromServiceController(c *fiber.Ctx) error {
	serviceID := c.Params("serviceID")
	volumeID := c.Query("volumeID")

	err := app.Service.DetachVolumeFromService(serviceID, volumeID)
	if err != nil {
		return sendResponse(c, fiber.StatusInternalServerError, "something went wrong", nil)
	}

	return sendResponse(c, fiber.StatusOK, "volume detached from service", nil)
}

// GetVolumeByIDController retrieves a volume by its public ID
func (app *App) GetVolumeByIDController(c *fiber.Ctx) error {
	publicID := c.Params("volumeID")

	volume, err := app.Service.GetVolumeByID(publicID)
	if err != nil {
		return sendResponse(c, fiber.StatusInternalServerError, "something went wrong", nil)
	}

	if volume == nil {
		return sendResponse(c, fiber.StatusNotFound, "volume not found", nil)
	}

	return sendResponse(c, fiber.StatusOK, "", volume)
}

// GetVolumesByWorkspaceController retrieves all volumes for a workspace
func (app *App) GetVolumesByWorkspaceController(c *fiber.Ctx) error {
	// workspaceID := c.Params("volumeID")
	// page := c.Query("page", "1")

	volumes, err := app.Service.GetVolumesByWorkspace(0, 0)
	if err != nil {
		return sendResponse(c, fiber.StatusInternalServerError, "something went wrong", nil)
	}

	return sendResponse(c, fiber.StatusOK, "", volumes)
}

// DeleteVolumeController handles deleting a volume
func (app *App) DeleteVolumeController(c *fiber.Ctx) error {
	publicID := c.Params("volumeID")

	err := app.Service.DeleteVolume(publicID)
	if err != nil {
		return sendResponse(c, fiber.StatusInternalServerError, "something went wrong", nil)
	}

	return sendResponse(c, fiber.StatusOK, "volume deleted", nil)
}
