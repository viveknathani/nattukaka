package controllers

import (
	"joystick/services"
	"joystick/shared"

	"github.com/gofiber/fiber/v2"
)

// ServiceController provides methods for service management.
type ServiceController struct {
	serviceService *services.ServiceService
}

// NewServiceController creates a new instance of ServiceController with the provided state.
func NewServiceController(serviceService *services.ServiceService) *ServiceController {
	return &ServiceController{
		serviceService: serviceService,
	}
}

// CreateService creates a new service in the database if all the checks pass.
func (serviceController *ServiceController) CreateService(c *fiber.Ctx) error {
	var serviceCreateRequest shared.ServiceCreateRequest

	if err := shared.ParseBody(c, &serviceCreateRequest); err != nil {
		return shared.SendStandardResponse(c, fiber.StatusBadRequest, nil, err.Error())
	}

	userID := shared.GetUserID(c)

	service, err := serviceController.serviceService.CreateService(userID, &serviceCreateRequest)
	if err != nil {
		return shared.SendStandardResponse(c, fiber.StatusInternalServerError, nil, err.Error())
	}

	return shared.SendStandardResponse(
		c,
		fiber.StatusCreated,
		&map[string]interface{}{"service": service},
		"",
	)
}

// GetAllServices returns all services from the database.
func (serviceController *ServiceController) GetAllServices(c *fiber.Ctx) error {
	services, err := serviceController.serviceService.GetAllServices(shared.GetUserID(c))
	if err != nil {
		return shared.SendStandardResponse(c, fiber.StatusInternalServerError, nil, err.Error())
	}
	return shared.SendStandardResponse(
		c,
		fiber.StatusOK,
		&map[string]interface{}{"services": services},
		"",
	)
}

// GetService returns a service from the database if it exists.
func (serviceController *ServiceController) GetService(c *fiber.Ctx) error {
	service, err := serviceController.serviceService.GetService(
		shared.GetUserID(c),
		c.Params("serviceID"),
	)
	if err != nil {
		return shared.SendStandardResponse(c, fiber.StatusInternalServerError, nil, err.Error())
	}
	return shared.SendStandardResponse(
		c,
		fiber.StatusOK,
		&map[string]interface{}{"service": service},
		"",
	)
}

// UpdateService updates a service in the database if it exists.
func (serviceController *ServiceController) UpdateService(c *fiber.Ctx) error {
	var serviceUpdateRequest shared.ServiceUpdateRequest

	if err := shared.ParseBody(c, &serviceUpdateRequest); err != nil {
		return shared.SendStandardResponse(c, fiber.StatusBadRequest, nil, err.Error())
	}

	userID := shared.GetUserID(c)

	service, err := serviceController.serviceService.UpdateService(
		userID,
		c.Params("serviceID"),
		&serviceUpdateRequest,
	)
	if err != nil {
		return shared.SendStandardResponse(c, fiber.StatusInternalServerError, nil, err.Error())
	}
	return shared.SendStandardResponse(
		c,
		fiber.StatusOK,
		&map[string]interface{}{"service": service},
		"",
	)
}

// DeleteService deletes a service from the database if it exists.
func (serviceController *ServiceController) DeleteService(c *fiber.Ctx) error {
	userID := shared.GetUserID(c)
	serviceID := c.Params("serviceID")

	// Use deployment service to handle complete cleanup
	err := serviceController.serviceService.DeleteService(userID, serviceID)
	if err != nil {
		return shared.SendStandardResponse(c, fiber.StatusInternalServerError, nil, err.Error())
	}

	return shared.SendStandardResponse(
		c,
		fiber.StatusOK,
		&map[string]interface{}{"message": "Service deleted successfully"},
		"",
	)
}
