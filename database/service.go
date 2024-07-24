package database

import (
	"database/sql"
	"time"

	"github.com/viveknathani/nattukaka/types"
)

// SQL statements as constants
const (
	statementInsertService             = `insert into services (public_id, name, status, type, runtime, workspace_id, created_by, instance_type_id, internal_url, external_url) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) returning id`
	statementInsertWebService          = `insert into web_services (public_id, service_id, repository, branch, root_directory, build_command, pre_deploy_command, start_command, health_check_path, environment) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	statementInsertDatabaseService     = `insert into database_services (public_id, service_id) values ($1, $2)`
	statementSelectServiceByID         = `select id, public_id, name, status, type, runtime, workspace_id, created_by, last_deployed_at, created_at, instance_type_id, internal_url, external_url from services where public_id = $1`
	statementSelectServicesByWorkspace = `select id, public_id, name, status, type, runtime, workspace_id, created_by, last_deployed_at, created_at, instance_type_id, internal_url, external_url from services where workspace_id = $1 limit 10 offset $2`
	statementUpdateService             = `update services set name = $1 where public_id = $2`
	statementUpdateWebService          = `update web_services set repository = $1, branch = $2, root_directory = $3, build_command = $4, pre_deploy_command = $5, start_command = $6, health_check_path = $7, environment = $8 where public_id = $9`
	statementDeleteService             = `delete from services where public_id = $1`
	statementDeleteWebService          = `delete from web_services where service_id = $1`
	statementDeleteDatabaseService     = `delete from database_services where service_id = $1`
	statementUpdateServiceStatus       = `update services set status = $1 where public_id = $2`
	statementMarkServiceAsDeployed     = `update services set status = $1, last_deployed_at = $2 where public_id = $3`
)

// InsertService inserts a new service into the database
func (db *Database) InsertService(service *types.Service) (int, error) {
	var id int
	err := db.execWithTransaction(statementInsertService, &id, service.PublicID, service.Name, service.Status, service.Type, service.Runtime, service.WorkspaceID, service.CreatedBy, service.InstanceTypeID, service.InternalURL, service.ExternalURL)
	return id, err
}

// InsertWebService inserts a new web service into the database
func (db *Database) InsertWebService(webService *types.WebService) error {
	return db.execWithTransaction(statementInsertWebService, webService.PublicID, webService.ServiceID, webService.Repository, webService.Branch, webService.RootDirectory, webService.BuildCommand, webService.PreDeployCommand, webService.StartCommand, webService.HealthCheckPath, webService.Environment)
}

// InsertDatabaseService inserts a new database service into the database
func (db *Database) InsertDatabaseService(databaseService *types.DatabaseService) error {
	return db.execWithTransaction(statementInsertDatabaseService, databaseService.PublicID, databaseService.ServiceID)
}

// GetServiceByID fetches a service by its public ID
func (db *Database) GetServiceByID(publicID string) (*types.Service, error) {
	var service types.Service
	err := db.query(statementSelectServiceByID, func(rows *sql.Rows) error {
		if rows.Next() {
			err := rows.Scan(&service.ID, &service.PublicID, &service.Name, &service.Status, &service.Type, &service.Runtime, &service.WorkspaceID, &service.CreatedBy, &service.LastDeployedAt, &service.CreatedAt, &service.InstanceTypeID, &service.InternalURL, &service.ExternalURL)
			if err != nil {
				return err
			}
		}
		return nil
	}, publicID)
	if err != nil {
		return nil, err
	}
	return &service, nil
}

// GetServicesByWorkspace fetches all services in a workspace, paginated
func (db *Database) GetServicesByWorkspace(workspaceID int, offset int) ([]types.Service, error) {
	var services []types.Service
	err := db.query(statementSelectServicesByWorkspace, func(rows *sql.Rows) error {
		for rows.Next() {
			var service types.Service
			err := rows.Scan(&service.ID, &service.PublicID, &service.Name, &service.Status, &service.Type, &service.Runtime, &service.WorkspaceID, &service.CreatedBy, &service.LastDeployedAt, &service.CreatedAt, &service.InstanceTypeID, &service.InternalURL, &service.ExternalURL)
			if err != nil {
				return err
			}
			services = append(services, service)
		}
		return nil
	}, workspaceID, offset)
	return services, err
}

// UpdateService updates a service in the database
func (db *Database) UpdateService(service *types.Service) error {
	return db.execWithTransaction(statementUpdateService, service.Name, service.PublicID)
}

// UpdateWebService updates a web service in the database
func (db *Database) UpdateWebService(webService *types.WebService) error {
	return db.execWithTransaction(statementUpdateWebService, webService.Repository, webService.Branch, webService.RootDirectory, webService.BuildCommand, webService.PreDeployCommand, webService.StartCommand, webService.HealthCheckPath, webService.Environment, webService.PublicID)
}

// DeleteService deletes a service from the database
func (db *Database) DeleteService(publicID string) error {
	service, err := db.GetServiceByID(publicID)
	if err != nil {
		return err
	}
	err = db.execWithTransaction(statementDeleteService, publicID)
	if err != nil {
		return err
	}
	if service.Type == "WEB" {
		err = db.execWithTransaction(statementDeleteWebService, service.ID)
	} else if service.Type == "DATABASE" {
		err = db.execWithTransaction(statementDeleteDatabaseService, service.ID)
	}
	return err
}

// UpdateServiceStatus updates the status of a service
func (db *Database) UpdateServiceStatus(publicID string, status string) error {
	return db.execWithTransaction(statementUpdateServiceStatus, status, publicID)
}

// MarkServiceAsDeployed marks a service as deployed
func (db *Database) MarkServiceAsDeployed(publicID string, status string) error {
	return db.execWithTransaction(statementMarkServiceAsDeployed, status, time.Now(), publicID)
}
