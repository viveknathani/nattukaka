package services

import (
	"joystick/shared"
)

// ServiceService provides methods for service management
type ServiceService struct {
	state *shared.State
}

// NewServiceService creates a new instance of ServiceService with the provided state.
func NewServiceService(state *shared.State) *ServiceService {
	return &ServiceService{
		state: state,
	}
}

// CreateService creates a new service in the database if all the checks pass.
func (serviceService *ServiceService) CreateService(
	userID int,
	serviceCreateRequest *shared.ServiceCreateRequest,
) (*shared.Service, *shared.JoyStickError) {
	service := &shared.Service{
		Name:          serviceCreateRequest.Name,
		RepositoryURL: serviceCreateRequest.RepositoryURL,
		Branch:        serviceCreateRequest.Branch,
		EnvVars:       serviceCreateRequest.EnvVars,
		PortMapping:   serviceCreateRequest.PortMapping,
		OwnerID:       userID,
	}

	err := serviceService.state.Database.Table("services").Create(service).Error
	if err != nil {
		serviceService.state.Logger.Error("error creating service: " + err.Error())
		return nil, shared.ErrInternalServerError
	}

	return service, nil
}

// GetService returns a service from the database if it exists.
func (serviceService *ServiceService) GetService(
	userID int,
	serviceID string,
) (*shared.Service, *shared.JoyStickError) {
	var service *shared.Service
	err := serviceService.state.Database.Table("services").
		Where("uuid = ? and owner_id = ?", serviceID, userID).
		First(&service).
		Error
	if err != nil {
		return nil, shared.ErrServiceNotFound
	}
	return service, nil
}

// GetAllServices returns all services from the database.
func (serviceService *ServiceService) GetAllServices(
	userID int,
) ([]shared.Service, *shared.JoyStickError) {
	var services []shared.Service
	err := serviceService.state.Database.Table("services").
		Where("owner_id = ?", userID).
		Find(&services).
		Error
	if err != nil {
		return nil, shared.ErrServiceNotFound
	}
	return services, nil
}

// UpdateService updates a service in the database if it exists.
func (serviceService *ServiceService) UpdateService(
	userID int,
	serviceID string,
	serviceUpdateRequest *shared.ServiceUpdateRequest,
) (*shared.Service, *shared.JoyStickError) {
	var service *shared.Service
	err := serviceService.state.Database.Table("services").
		Where("uuid = ? and owner_id = ?", serviceID, userID).
		First(&service).
		Error
	if err != nil {
		return nil, shared.ErrServiceNotFound
	}

	service.Name = serviceUpdateRequest.Name
	service.RepositoryURL = serviceUpdateRequest.RepositoryURL
	service.Branch = serviceUpdateRequest.Branch
	service.EnvVars = serviceUpdateRequest.EnvVars
	service.PortMapping = serviceUpdateRequest.PortMapping

	err = serviceService.state.Database.Table("services").
		Where("uuid = ? and owner_id = ?", serviceID, userID).
		UpdateColumns(service).
		Error
	if err != nil {
		serviceService.state.Logger.Error("error updating service: " + err.Error())
		return nil, shared.ErrInternalServerError
	}
	return service, nil
}

// DeleteService deletes a service from the database if it exists.
func (serviceService *ServiceService) DeleteService(
	userID int,
	serviceID string,
) (*shared.Service, *shared.JoyStickError) {
	var service *shared.Service
	err := serviceService.state.Database.Table("services").
		Where("uuid = ? and owner_id = ?", serviceID, userID).
		First(&service).
		Error
	if err != nil {
		return nil, shared.ErrServiceNotFound
	}

	err = serviceService.state.Database.Table("services").
		Where("uuid = ? and owner_id = ?", serviceID, userID).
		Delete(&service).
		Error
	if err != nil {
		serviceService.state.Logger.Error("error deleting service: " + err.Error())
		return nil, shared.ErrInternalServerError
	}
	return service, nil
}
