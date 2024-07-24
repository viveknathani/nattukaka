package service

import (
	"github.com/viveknathani/nattukaka/types"
)

// CreateService creates a new service
func (srv *Service) CreateService(service *types.Service, webService *types.WebService, databaseService *types.DatabaseService) error {
	serviceID, err := srv.Db.InsertService(service)
	if err != nil {
		srv.Logger.Error(err.Error())
		return err
	}

	if service.Type == "WEB" {
		webService.ServiceID = serviceID
		err = srv.Db.InsertWebService(webService)
		if err != nil {
			srv.Logger.Error(err.Error())
			return err
		}
	} else if service.Type == "DATABASE" {
		databaseService.ServiceID = serviceID
		err = srv.Db.InsertDatabaseService(databaseService)
		if err != nil {
			srv.Logger.Error(err.Error())
			return err
		}
	}

	return nil
}

// GetServiceByID fetches a service by its public ID
func (srv *Service) GetServiceByID(publicID string) (*types.Service, error) {
	service, err := srv.Db.GetServiceByID(publicID)
	if err != nil {
		srv.Logger.Error(err.Error())
		return nil, err
	}
	return service, nil
}

// GetServicesByWorkspace fetches all services in a workspace
func (srv *Service) GetServicesByWorkspace(workspaceID int, page int) ([]types.Service, error) {
	offset := (page - 1) * 10
	services, err := srv.Db.GetServicesByWorkspace(workspaceID, offset)
	if err != nil {
		srv.Logger.Error(err.Error())
		return nil, err
	}
	return services, nil
}

// UpdateService updates a service
func (srv *Service) UpdateService(service *types.Service, webService *types.WebService, databaseService *types.DatabaseService) error {
	err := srv.Db.UpdateService(service)
	if err != nil {
		srv.Logger.Error(err.Error())
		return err
	}

	if service.Type == "WEB" {
		err = srv.Db.UpdateWebService(webService)
		if err != nil {
			srv.Logger.Error(err.Error())
			return err
		}
	} else if service.Type == "DATABASE" {
		// No-op for now
	}

	return nil
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
