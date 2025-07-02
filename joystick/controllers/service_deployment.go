package controllers

import (
	"joystick/services"
	"joystick/shared"

	"github.com/gofiber/fiber/v2"
)

// ServiceDeploymentController represents a service deployment controller.
type ServiceDeploymentController struct {
	serviceDeploymentService *services.ServiceDeploymentService
}

// NewServiceDeploymentController creates a new instance of ServiceDeploymentController with the provided state.
func NewServiceDeploymentController(
	serviceDeploymentService *services.ServiceDeploymentService,
) *ServiceDeploymentController {
	return &ServiceDeploymentController{
		serviceDeploymentService: serviceDeploymentService,
	}
}

// CreateServiceDeployment creates a new service deployment.
func (serviceDeploymentController *ServiceDeploymentController) CreateServiceDeployment(
	c *fiber.Ctx,
) error {
	serviceID := c.Params("serviceID")

	// Create the service deployment
	serviceDeployment, err := serviceDeploymentController.serviceDeploymentService.CreateServiceDeployment(
		serviceID,
	)
	if err != nil {
		return shared.SendStandardResponse(c, fiber.StatusInternalServerError, nil, err.Error())
	}

	// Return the service deployment
	return shared.SendStandardResponse(
		c,
		fiber.StatusOK,
		&map[string]interface{}{"deployment": serviceDeployment},
		"",
	)
}

// UpdateServiceDeployment updates a service deployment.
func (serviceDeploymentController *ServiceDeploymentController) UpdateServiceDeployment(
	c *fiber.Ctx,
) error {
	serviceDeploymentID := c.Params("deploymentID")

	var serviceDeploymentUpdateRequest shared.ServiceDeploymentUpdateRequest

	if err := shared.ParseBody(c, &serviceDeploymentUpdateRequest); err != nil {
		return shared.SendStandardResponse(c, fiber.StatusBadRequest, nil, err.Error())
	}

	// Update the service deployment
	serviceDeployment, err := serviceDeploymentController.serviceDeploymentService.UpdateServiceDeployment(
		serviceDeploymentID,
		&serviceDeploymentUpdateRequest,
	)
	if err != nil {
		return shared.SendStandardResponse(c, fiber.StatusInternalServerError, nil, err.Error())
	}

	// Return the service deployment
	return shared.SendStandardResponse(
		c,
		fiber.StatusOK,
		&map[string]interface{}{"deployment": serviceDeployment},
		"",
	)
}

// GetAllServiceDeployments returns all service deployments.
func (serviceDeploymentController *ServiceDeploymentController) GetAllServiceDeployments(
	c *fiber.Ctx,
) error {
	// Get all service deployments
	serviceDeployments, err := serviceDeploymentController.serviceDeploymentService.GetAllServiceDeployments()
	if err != nil {
		return shared.SendStandardResponse(c, fiber.StatusInternalServerError, nil, err.Error())
	}

	// Return the service deployments
	return shared.SendStandardResponse(
		c,
		fiber.StatusOK,
		&map[string]interface{}{"deployments": serviceDeployments},
		"",
	)
}
