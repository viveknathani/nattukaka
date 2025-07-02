package services

import (
	"joystick/shared"
	"time"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/config"
	"github.com/go-git/go-git/v6/storage/memory"
)

// ServiceDeploymentService represents a service deployment service.
type ServiceDeploymentService struct {
	state *shared.State
}

// NewServiceDeploymentService creates a new service deployment service.
func NewServiceDeploymentService(state *shared.State) *ServiceDeploymentService {
	return &ServiceDeploymentService{state: state}
}

// CreateServiceDeployment creates a new service deployment.
func (serviceDeploymentService *ServiceDeploymentService) CreateServiceDeployment(
	serviceID string,
) (*shared.ServiceDeployment, *shared.JoyStickError) {
	var service *shared.Service
	err := serviceDeploymentService.state.Database.Table("services").
		Where("uuid = ?", serviceID).
		First(&service).
		Error
	if err != nil {
		return nil, shared.ErrServiceNotFound
	}

	latestCommitHash, fetchErr := serviceDeploymentService.getLatestCommitHash(
		service.RepositoryURL,
		service.Branch,
	)
	if fetchErr != nil {
		return nil, shared.ErrInternalServerError
	}

	serviceDeployment := &shared.ServiceDeployment{
		ServiceID:   service.ID,
		Commit:      latestCommitHash,
		Status:      shared.ServiceDeploymentStatusQueued,
		ContainerID: "",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = serviceDeploymentService.state.Database.Table("service_deployments").
		Create(&serviceDeployment).
		Error

	if err != nil {
		return nil, shared.ErrInternalServerError
	}
	return serviceDeployment, nil
}

// UpdateServiceDeployment updates a service deployment.
func (serviceDeploymentService *ServiceDeploymentService) UpdateServiceDeployment(
	serviceDeploymentID string,
	serviceDeploymentUpdateRequest *shared.ServiceDeploymentUpdateRequest,
) (*shared.ServiceDeployment, *shared.JoyStickError) {
	var serviceDeployment *shared.ServiceDeployment
	err := serviceDeploymentService.state.Database.Table("service_deployments").
		Where("uuid = ?", serviceDeploymentID).
		First(&serviceDeployment).
		Error
	if err != nil {
		return nil, shared.ErrServiceDeploymentNotFound
	}

	if serviceDeploymentUpdateRequest.Status != "" {
		serviceDeployment.Status = serviceDeploymentUpdateRequest.Status
	}
	if serviceDeploymentUpdateRequest.ContainerID != "" {
		serviceDeployment.ContainerID = serviceDeploymentUpdateRequest.ContainerID
	}

	err = serviceDeploymentService.state.Database.Table("service_deployments").
		Where("uuid = ?", serviceDeploymentID).
		UpdateColumns(serviceDeployment).
		Error
	if err != nil {
		serviceDeploymentService.state.Logger.Error(
			"error updating service deployment: " + err.Error(),
		)
		return nil, shared.ErrInternalServerError
	}
	return serviceDeployment, nil
}

// GetAllServiceDeployments returns all service deployments, ordered by created_at in descending order, limited to top 10.
func (serviceDeploymentService *ServiceDeploymentService) GetAllServiceDeployments() ([]shared.ServiceDeployment, *shared.JoyStickError) {
	var serviceDeployments []shared.ServiceDeployment
	err := serviceDeploymentService.state.Database.Table("service_deployments").
		Order("created_at desc").
		Limit(10).
		Find(&serviceDeployments).
		Error
	if err != nil {
		return nil, shared.ErrInternalServerError
	}
	return serviceDeployments, nil
}

// getLatestCommitHash returns the latest commit hash of a repository.
func (serviceDeploymentService *ServiceDeploymentService) getLatestCommitHash(
	repositoryURL string,
	branch string,
) (string, *shared.JoyStickError) {
	remote := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{repositoryURL},
	})

	refs, err := remote.List(&git.ListOptions{})

	targetRef := "refs/heads/" + branch
	if err != nil {
		serviceDeploymentService.state.Logger.Error("error listing remote refs: " + err.Error())
		return "", shared.ErrInternalServerError
	}

	for _, ref := range refs {
		if ref.Name().String() == targetRef {
			return ref.Hash().String(), nil
		}
	}

	serviceDeploymentService.state.Logger.Error(
		"error getting latest commit hash: " + "no ref found for branch " + branch,
	)
	return "", shared.ErrInternalServerError
}
