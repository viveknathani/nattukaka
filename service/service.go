package service

import (
	"github.com/gofiber/fiber/v2"
	"github.com/viveknathani/nattukaka/types"
	"github.com/viveknathani/nattukaka/utils"
)

// CreateService creates a new service
func (srv *Service) CreateService(email string, req *types.CreateServiceRequest) (int, string) {

	// Validate payload
	if req.Base == nil || (req.Web == nil && req.Database == nil) {
		return fiber.StatusBadRequest, "invalid request body"
	}

	// Validate workspace
	workspace, err := srv.Db.GetWorkspaceByPublicID(req.WorkspaceID)
	if err != nil {
		return fiber.StatusInternalServerError, err.Error()
	}
	if workspace == nil {
		return fiber.StatusBadRequest, "workspace does not exist"
	}

	// Validate creator
	user, err := srv.Db.GetUserByEmail(email)
	if err != nil {
		return fiber.StatusInternalServerError, err.Error()
	}
	if user == nil {
		return fiber.StatusBadRequest, "user does not exist"
	}

	// Prepare service data
	publicIDForService, err := utils.GeneratePublicId("service")
	if err != nil {
		return fiber.StatusInternalServerError, err.Error()
	}
	service := &types.Service{
		PublicID:       publicIDForService,
		Name:           req.Base.Name,
		Status:         utils.ServiceStatusCreated,
		Type:           req.Base.Type,
		Runtime:        req.Base.Runtime,
		WorkspaceID:    workspace.ID,
		CreatedBy:      user.ID,
		InstanceTypeID: req.Base.InstanceTypeID,
		InternalURL:    "", // TO BE GENERATED
		ExternalURL:    "", // TO BE GENERATED
	}

	var webService *types.WebService
	var databaseService *types.DatabaseService

	if req.Base.Type == "WEB" {
		publicIDForWebService, err := utils.GeneratePublicId("web_service")
		if err != nil {
			return fiber.StatusInternalServerError, err.Error()
		}
		webService = &types.WebService{
			PublicID:         publicIDForWebService,
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
		publicIDForDatabaseService, err := utils.GeneratePublicId("database_service")
		if err != nil {
			return fiber.StatusInternalServerError, err.Error()
		}
		databaseService = &types.DatabaseService{
			PublicID: publicIDForDatabaseService,
		}
	}

	serviceID, err := srv.Db.InsertService(service)
	if err != nil {
		srv.Logger.Error(err.Error())
		return fiber.StatusInternalServerError, err.Error()
	}

	if service.Type == "WEB" {
		webService.ServiceID = serviceID
		err = srv.Db.InsertWebService(webService)
		if err != nil {
			srv.Logger.Error(err.Error())
			return fiber.StatusInternalServerError, err.Error()
		}
	} else if service.Type == "DATABASE" {
		databaseService.ServiceID = serviceID
		err = srv.Db.InsertDatabaseService(databaseService)
		if err != nil {
			srv.Logger.Error(err.Error())
			return fiber.StatusInternalServerError, err.Error()
		}
	}

	return fiber.StatusCreated, "service created!"
}

// GetServiceByID fetches a service by its public ID
func (srv *Service) GetServiceByID(publicID string) (int, string, *types.Service) {
	service, err := srv.Db.GetServiceByID(publicID)
	if err != nil {
		srv.Logger.Error(err.Error())
		return fiber.StatusInternalServerError, "something went wrong", nil
	}

	if service.ID == 0 {
		return fiber.StatusNotFound, "not found", nil
	}

	return fiber.StatusOK, "", service
}

// GetServicesByWorkspace fetches all services in a workspace
func (srv *Service) GetServicesByWorkspace(workspaceID string, page int) (int, string, []types.Service) {
	// Validate workspace
	workspace, err := srv.Db.GetWorkspaceByPublicID(workspaceID)
	if err != nil {
		return fiber.StatusInternalServerError, err.Error(), nil
	}
	if workspace == nil {
		return fiber.StatusBadRequest, "workspace does not exist", nil
	}

	// Get services
	offset := (page - 1) * utils.PageSize
	services, err := srv.Db.GetServicesByWorkspace(workspace.ID, utils.PageSize, offset)
	if err != nil {
		srv.Logger.Error(err.Error())
		return fiber.StatusInternalServerError, err.Error(), nil
	}
	return fiber.StatusOK, "", services
}

// UpdateService updates a service
func (srv *Service) UpdateService(publicIDService string, req *types.UpdateServiceRequest) (int, string) {
	// Fetch the existing service
	code, message, service := srv.GetServiceByID(publicIDService)
	if code != fiber.StatusOK {
		return code, message
	}

	// Update base properties
	if req.Base != nil && req.Base.Name != "" {
		service.Name = req.Base.Name
	}

	var webService *types.WebService

	if service.Type == "WEB" && req.Web != nil {
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
	}

	err := srv.Db.UpdateService(service)
	if err != nil {
		srv.Logger.Error(err.Error())
		return fiber.StatusInternalServerError, err.Error()
	}

	if service.Type == "WEB" {
		err = srv.Db.UpdateWebService(webService)
		if err != nil {
			srv.Logger.Error(err.Error())
			return fiber.StatusInternalServerError, err.Error()
		}
	}

	return fiber.StatusOK, ""
}

// DeleteService deletes a service
func (srv *Service) DeleteService(publicID string) error {
	err := srv.Db.DeleteService(publicID)
	if err != nil {
		srv.Logger.Error(err.Error())
		return err
	}
	return nil
}

// UpdateServiceStatus updates the status of a service
func (srv *Service) UpdateServiceStatus(publicID string, status string) error {
	err := srv.Db.UpdateServiceStatus(publicID, status)
	if err != nil {
		srv.Logger.Error(err.Error())
		return err
	}
	return nil
}

// MarkServiceAsDeployed marks a service as deployed
func (srv *Service) MarkServiceAsDeployed(publicID string, status string) error {
	err := srv.Db.MarkServiceAsDeployed(publicID, status)
	if err != nil {
		srv.Logger.Error(err.Error())
		return err
	}
	return nil
}
