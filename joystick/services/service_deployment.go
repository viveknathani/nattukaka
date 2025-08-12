package services

import (
	"context"
	"encoding/json"
	"fmt"
	"joystick/caddyadmin"
	"joystick/shared"
	"strings"
	"time"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/config"
	"github.com/go-git/go-git/v6/storage/memory"
)

// ServiceDeploymentService represents a service deployment service.
type ServiceDeploymentService struct {
	state        *shared.State
	playerClient *PlayerClient
}

// NewServiceDeploymentService creates a new service deployment service.
func NewServiceDeploymentService(
	state *shared.State,
	playerClient *PlayerClient,
) *ServiceDeploymentService {
	return &ServiceDeploymentService{state: state, playerClient: playerClient}
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

	nextNode, nodeErr := serviceDeploymentService.FindNextNode()
	if nodeErr != nil {
		return nil, shared.ErrInternalServerError
	}

	serviceDeployment := &shared.ServiceDeployment{
		ServiceID:   service.ID,
		Commit:      latestCommitHash,
		Status:      shared.ServiceDeploymentStatusQueued,
		ContainerID: "",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		NodeID:      nextNode.ID,
	}

	err = serviceDeploymentService.state.Database.Table("service_deployments").
		Create(&serviceDeployment).
		Error

	if err != nil {
		return nil, shared.ErrInternalServerError
	}

	serviceDeploymentService.ProcessDeployment(serviceDeployment.UUID)

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

// AddOrUpdateRoute updates a route in caddy.
func (serviceDeploymentService *ServiceDeploymentService) AddOrUpdateRoute(
	id string,
	hosts []string,
	upstream string,
) error {
	_, err := caddyadmin.GetRoute(id)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return caddyadmin.AddRoute(id, hosts, upstream)
		}
		return err
	}
	return caddyadmin.UpdateRoute(id, hosts, upstream)
}

func (serviceDeploymentService *ServiceDeploymentService) updateCaddyConfig(
	service *shared.Service,
) error {
	// Parse port mapping to get the host port
	var portMappings []shared.PortMapping
	if err := json.Unmarshal(service.PortMapping, &portMappings); err != nil {
		return fmt.Errorf("failed to parse port mappings: %w", err)
	}

	if len(portMappings) == 0 {
		return fmt.Errorf("no port mappings found for service")
	}

	// Use the first port mapping for the upstream
	hostPort := portMappings[0].HostPort
	upstream := fmt.Sprintf("localhost:%d", hostPort)

	// Create host configuration - using service name as subdomain
	hosts := []string{fmt.Sprintf("%s.localhost", service.Name)}

	// Update or add route in Caddy
	routeID := fmt.Sprintf("service-%d", service.ID)
	return serviceDeploymentService.AddOrUpdateRoute(routeID, hosts, upstream)
}

func (serviceDeploymentService *ServiceDeploymentService) ProcessDeployment(deploymentUUID string) {
	go serviceDeploymentService.processDeploymentAsync(deploymentUUID)
}

func (serviceDeploymentService *ServiceDeploymentService) processDeploymentAsync(
	deploymentUUID string,
) {
	ctx := context.Background()

	// Get deployment details
	var deployment *shared.ServiceDeployment
	err := serviceDeploymentService.state.Database.Table("service_deployments").
		Where("uuid = ?", deploymentUUID).
		First(&deployment).
		Error
	if err != nil {
		serviceDeploymentService.state.Logger.Error("failed to get deployment: " + err.Error())
		return
	}

	// Get service details
	var service *shared.Service
	err = serviceDeploymentService.state.Database.Table("services").
		Where("id = ?", deployment.ServiceID).
		First(&service).
		Error
	if err != nil {
		serviceDeploymentService.state.Logger.Error("failed to get service: " + err.Error())
		serviceDeploymentService.UpdateServiceDeployment(
			deploymentUUID,
			&shared.ServiceDeploymentUpdateRequest{
				Status: shared.ServiceDeploymentStatusFailed,
			},
		)
		return
	}

	// Update status to BUILDING
	serviceDeploymentService.UpdateServiceDeployment(
		deploymentUUID,
		&shared.ServiceDeploymentUpdateRequest{
			Status: shared.ServiceDeploymentStatusBuilding,
		},
	)

	// Step 1: Build image
	imageTag, err := serviceDeploymentService.playerClient.BuildImage(
		ctx,
		service.Name,
		service.RepositoryURL,
		service.Branch,
		deployment.Commit,
	)
	if err != nil {
		serviceDeploymentService.state.Logger.Error("failed to build image: " + err.Error())
		serviceDeploymentService.UpdateServiceDeployment(
			deploymentUUID,
			&shared.ServiceDeploymentUpdateRequest{
				Status: shared.ServiceDeploymentStatusFailed,
			},
		)
		return
	}

	serviceDeploymentService.state.Logger.Info(fmt.Sprintf("Built image: %s", imageTag))

	// Update status to DEPLOYING
	serviceDeploymentService.UpdateServiceDeployment(
		deploymentUUID,
		&shared.ServiceDeploymentUpdateRequest{
			Status: shared.ServiceDeploymentStatusDeploying,
		},
	)

	// Step 2: Stop and remove old container if it exists
	// Get the latest successful deployment to find the current container
	var latestDeployment *shared.ServiceDeployment
	latestErr := serviceDeploymentService.state.Database.Table("service_deployments").
		Where("service_id = ? AND status = ? AND container_id != ''", deployment.ServiceID, shared.ServiceDeploymentStatusSuccess).
		Order("created_at desc").
		First(&latestDeployment).
		Error

	if latestErr == nil && latestDeployment.ContainerID != "" {
		if err := serviceDeploymentService.playerClient.StopContainer(ctx, latestDeployment.ContainerID); err != nil {
			serviceDeploymentService.state.Logger.Error(
				"failed to stop old container: " + err.Error(),
			)
			// continue anyway, the container might already be stopped
		}

		if err := serviceDeploymentService.playerClient.RemoveContainer(ctx, latestDeployment.ContainerID); err != nil {
			serviceDeploymentService.state.Logger.Error(
				"failed to remove old container: " + err.Error(),
			)
			// continue anyway, the container might already be removed
		}
	}

	// Step 3: Parse environment variables
	var envVars map[string]string
	if err := json.Unmarshal(service.EnvVars, &envVars); err != nil {
		serviceDeploymentService.state.Logger.Error("failed to parse env vars: " + err.Error())
		serviceDeploymentService.UpdateServiceDeployment(
			deploymentUUID,
			&shared.ServiceDeploymentUpdateRequest{
				Status: shared.ServiceDeploymentStatusFailed,
			},
		)
		return
	}

	// Step 4: Create and start new container
	containerID, err := serviceDeploymentService.playerClient.CreateAndStartContainer(
		ctx,
		imageTag,
		envVars,
	)
	if err != nil {
		serviceDeploymentService.state.Logger.Error(
			"failed to create and start container: " + err.Error(),
		)
		serviceDeploymentService.UpdateServiceDeployment(
			deploymentUUID,
			&shared.ServiceDeploymentUpdateRequest{
				Status: shared.ServiceDeploymentStatusFailed,
			},
		)
		return
	}

	serviceDeploymentService.state.Logger.Info(fmt.Sprintf("Started container: %s", containerID))

	// Step 5: Update deployment with container ID
	_, updateErr := serviceDeploymentService.UpdateServiceDeployment(
		deploymentUUID,
		&shared.ServiceDeploymentUpdateRequest{
			ContainerID: containerID,
			Status:      shared.ServiceDeploymentStatusRunning,
		},
	)
	if updateErr != nil {
		serviceDeploymentService.state.Logger.Error(
			"failed to update deployment: " + updateErr.Error(),
		)
		// Don't mark as failed since container is running
	}

	// Step 6: Update Caddy configuration
	if err := serviceDeploymentService.updateCaddyConfig(service); err != nil {
		serviceDeploymentService.state.Logger.Error("failed to update caddy config: " + err.Error())
		// don't fail deployment for Caddy issues, log and continue
	}

	// Mark deployment as successful
	serviceDeploymentService.UpdateServiceDeployment(
		deploymentUUID,
		&shared.ServiceDeploymentUpdateRequest{
			Status: shared.ServiceDeploymentStatusSuccess,
		},
	)
	serviceDeploymentService.state.Logger.Info(
		fmt.Sprintf("Deployment %s completed successfully", deploymentUUID),
	)
}

func (serviceDeploymentService *ServiceDeploymentService) RemoveLatestDeployment(
	serviceId int,
) error {
	ctx := context.Background()

	// Get latest deployment for this service
	var deployment *shared.ServiceDeployment
	dbErr := serviceDeploymentService.state.Database.Table("service_deployments").
		Where("service_id = ?", serviceId).
		Order("created_at desc").
		First(&deployment).
		Error

	// If there's a running container, stop and remove it
	if dbErr == nil && deployment.ContainerID != "" {
		if err := serviceDeploymentService.playerClient.StopContainer(ctx, deployment.ContainerID); err != nil {
			serviceDeploymentService.state.Logger.Error(
				"failed to stop container during deletion: " + err.Error(),
			)
		}

		if err := serviceDeploymentService.playerClient.RemoveContainer(ctx, deployment.ContainerID); err != nil {
			serviceDeploymentService.state.Logger.Error(
				"failed to remove container during deletion: " + err.Error(),
			)
		}
	}

	// Remove Caddy configuration
	routeID := fmt.Sprintf("service-%d", serviceId)
	if err := caddyadmin.DeleteRoute(routeID); err != nil {
		serviceDeploymentService.state.Logger.Error("failed to remove caddy route: " + err.Error())
	}

	return nil
}

// FindNextNode returns the next node that can run a service deployment.
// Right now, it just returns the first node in the database.
// A more sophisticated design would involve:
// 1. Checking if the node is online
// 2. Checking if the node has enough resources to run the service deployment
func (serviceDeploymentService *ServiceDeploymentService) FindNextNode() (*shared.Node, *shared.JoyStickError) {
	var node *shared.Node
	err := serviceDeploymentService.state.Database.Table("nodes").
		Order("created_at desc").
		Limit(1).
		First(&node).
		Error
	if err != nil {
		return nil, shared.ErrInternalServerError
	}
	return node, nil
}
