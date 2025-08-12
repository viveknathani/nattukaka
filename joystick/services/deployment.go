package services

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"joystick/caddyadmin"
// 	"joystick/shared"
// )

// // DeploymentService handles deployment operations
// type DeploymentService struct {
// 	state                    *shared.State
// 	serviceDeploymentService *ServiceDeploymentService
// 	serviceService           *ServiceService
// 	playerClient             *PlayerClient
// }

// // NewDeploymentService creates a new deployment service
// func NewDeploymentService(
// 	state *shared.State,
// 	serviceDeploymentService *ServiceDeploymentService,
// 	serviceService *ServiceService,
// 	playerClient *PlayerClient,
// ) *DeploymentService {
// 	return &DeploymentService{
// 		state:                    state,
// 		serviceDeploymentService: serviceDeploymentService,
// 		serviceService:           serviceService,
// 		playerClient:             playerClient,
// 	}
// }

// // ProcessDeployment handles the complete deployment workflow
// func (ds *DeploymentService) ProcessDeployment(deploymentUUID string) {
// 	go ds.processDeploymentAsync(deploymentUUID)
// }

// func (ds *DeploymentService) processDeploymentAsync(deploymentUUID string) {
// 	ctx := context.Background()

// 	// Get deployment details
// 	var deployment *shared.ServiceDeployment
// 	err := ds.state.Database.Table("service_deployments").
// 		Where("uuid = ?", deploymentUUID).
// 		First(&deployment).
// 		Error
// 	if err != nil {
// 		ds.state.Logger.Error("failed to get deployment: " + err.Error())
// 		return
// 	}

// 	// Get service details
// 	var service *shared.Service
// 	err = ds.state.Database.Table("services").
// 		Where("id = ?", deployment.ServiceID).
// 		First(&service).
// 		Error
// 	if err != nil {
// 		ds.state.Logger.Error("failed to get service: " + err.Error())
// 		ds.updateDeploymentStatus(deploymentUUID, shared.ServiceDeploymentStatusFailed)
// 		return
// 	}

// 	// Update status to BUILDING
// 	ds.updateDeploymentStatus(deploymentUUID, shared.ServiceDeploymentStatusBuilding)

// 	// Step 1: Build image
// 	imageTag, err := ds.playerClient.BuildImage(
// 		ctx,
// 		service.Name,
// 		service.RepositoryURL,
// 		service.Branch,
// 		deployment.Commit,
// 	)
// 	if err != nil {
// 		ds.state.Logger.Error("failed to build image: " + err.Error())
// 		ds.updateDeploymentStatus(deploymentUUID, shared.ServiceDeploymentStatusFailed)
// 		return
// 	}

// 	ds.state.Logger.Info(fmt.Sprintf("Built image: %s", imageTag))

// 	// Update status to DEPLOYING
// 	ds.updateDeploymentStatus(deploymentUUID, shared.ServiceDeploymentStatusDeploying)

// 	// Step 2: Stop and remove old container if it exists
// 	// Get the latest successful deployment to find the current container
// 	var latestDeployment *shared.ServiceDeployment
// 	latestErr := ds.state.Database.Table("service_deployments").
// 		Where("service_id = ? AND status = ? AND container_id != ''", deployment.ServiceID, shared.ServiceDeploymentStatusSuccess).
// 		Order("created_at desc").
// 		First(&latestDeployment).
// 		Error

// 	if latestErr == nil && latestDeployment.ContainerID != "" {
// 		if err := ds.playerClient.StopContainer(ctx, latestDeployment.ContainerID); err != nil {
// 			ds.state.Logger.Error("failed to stop old container: " + err.Error())
// 			// Continue anyway, the container might already be stopped
// 		}

// 		if err := ds.playerClient.RemoveContainer(ctx, latestDeployment.ContainerID); err != nil {
// 			ds.state.Logger.Error("failed to remove old container: " + err.Error())
// 			// Continue anyway, the container might already be removed
// 		}
// 	}

// 	// Step 3: Parse environment variables
// 	var envVars map[string]string
// 	if err := json.Unmarshal(service.EnvVars, &envVars); err != nil {
// 		ds.state.Logger.Error("failed to parse env vars: " + err.Error())
// 		ds.updateDeploymentStatus(deploymentUUID, shared.ServiceDeploymentStatusFailed)
// 		return
// 	}

// 	// Step 4: Create and start new container
// 	containerID, err := ds.playerClient.CreateAndStartContainer(ctx, imageTag, envVars)
// 	if err != nil {
// 		ds.state.Logger.Error("failed to create and start container: " + err.Error())
// 		ds.updateDeploymentStatus(deploymentUUID, shared.ServiceDeploymentStatusFailed)
// 		return
// 	}

// 	ds.state.Logger.Info(fmt.Sprintf("Started container: %s", containerID))

// 	// Step 5: Update deployment with container ID
// 	_, updateErr := ds.serviceDeploymentService.UpdateServiceDeployment(deploymentUUID, &shared.ServiceDeploymentUpdateRequest{
// 		ContainerID: containerID,
// 		Status:      shared.ServiceDeploymentStatusRunning,
// 	})
// 	if updateErr != nil {
// 		ds.state.Logger.Error("failed to update deployment: " + updateErr.Error())
// 		// Don't mark as failed since container is running
// 	}

// 	// Step 6: Update Caddy configuration
// 	if err := ds.updateCaddyConfig(service, deployment); err != nil {
// 		ds.state.Logger.Error("failed to update caddy config: " + err.Error())
// 		// Don't fail deployment for Caddy issues, log and continue
// 	}

// 	// Mark deployment as successful
// 	ds.updateDeploymentStatus(deploymentUUID, shared.ServiceDeploymentStatusSuccess)
// 	ds.state.Logger.Info(fmt.Sprintf("Deployment %s completed successfully", deploymentUUID))
// }

// func (ds *DeploymentService) updateDeploymentStatus(deploymentUUID, status string) {
// 	_, err := ds.serviceDeploymentService.UpdateServiceDeployment(deploymentUUID, &shared.ServiceDeploymentUpdateRequest{
// 		Status: status,
// 	})
// 	if err != nil {
// 		ds.state.Logger.Error("failed to update deployment status: " + err.Error())
// 	}
// }

// func (ds *DeploymentService) updateCaddyConfig(service *shared.Service, deployment *shared.ServiceDeployment) error {
// 	// Parse port mapping to get the host port
// 	var portMappings []shared.PortMapping
// 	if err := json.Unmarshal(service.PortMapping, &portMappings); err != nil {
// 		return fmt.Errorf("failed to parse port mappings: %w", err)
// 	}

// 	if len(portMappings) == 0 {
// 		return fmt.Errorf("no port mappings found for service")
// 	}

// 	// Use the first port mapping for the upstream
// 	hostPort := portMappings[0].HostPort
// 	upstream := fmt.Sprintf("localhost:%d", hostPort)

// 	// Create host configuration - using service name as subdomain
// 	hosts := []string{fmt.Sprintf("%s.localhost", service.Name)}

// 	// Update or add route in Caddy
// 	routeID := fmt.Sprintf("service-%d", service.ID)
// 	return ds.serviceDeploymentService.AddOrUpdateRoute(routeID, hosts, upstream)
// }

// // ProcessServiceDeletion handles service deletion with cleanup
// func (ds *DeploymentService) ProcessServiceDeletion(serviceUUID string, userID int) error {
// 	ctx := context.Background()

// 	// Get service details
// 	service, err := ds.serviceService.GetService(userID, serviceUUID)
// 	if err != nil {
// 		return fmt.Errorf("failed to get service: %w", err)
// 	}

// 	// Get latest deployment for this service
// 	var deployment *shared.ServiceDeployment
// 	dbErr := ds.state.Database.Table("service_deployments").
// 		Where("service_id = ?", service.ID).
// 		Order("created_at desc").
// 		First(&deployment).
// 		Error

// 	// If there's a running container, stop and remove it
// 	if dbErr == nil && deployment.ContainerID != "" {
// 		if err := ds.playerClient.StopContainer(ctx, deployment.ContainerID); err != nil {
// 			ds.state.Logger.Error("failed to stop container during deletion: " + err.Error())
// 		}

// 		if err := ds.playerClient.RemoveContainer(ctx, deployment.ContainerID); err != nil {
// 			ds.state.Logger.Error("failed to remove container during deletion: " + err.Error())
// 		}
// 	}

// 	// Remove Caddy configuration
// 	routeID := fmt.Sprintf("service-%d", service.ID)
// 	if err := ds.removeCaddyRoute(routeID); err != nil {
// 		ds.state.Logger.Error("failed to remove caddy route: " + err.Error())
// 	}

// 	// Delete all deployments for this service
// 	dbErr = ds.state.Database.Table("service_deployments").
// 		Where("service_id = ?", service.ID).
// 		Delete(&shared.ServiceDeployment{}).
// 		Error
// 	if dbErr != nil {
// 		return fmt.Errorf("failed to delete service deployments: %w", dbErr)
// 	}

// 	// Delete the service
// 	_, serviceErr := ds.serviceService.DeleteService(userID, serviceUUID)
// 	if serviceErr != nil {
// 		return fmt.Errorf("failed to delete service: %w", serviceErr)
// 	}

// 	ds.state.Logger.Info(fmt.Sprintf("Service %s deleted successfully", serviceUUID))
// 	return nil
// }

// func (ds *DeploymentService) removeCaddyRoute(routeID string) error {
// 	return caddyadmin.DeleteRoute(routeID)
// }
