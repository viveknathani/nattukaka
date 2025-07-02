package shared

import "time"

const (
	// EnvDevelopment ...
	EnvDevelopment = "development"

	// EnvProduction ...
	EnvProduction = "production"

	// TokenAge ...
	TokenAge = time.Hour * 24 * 2

	// ServiceDeploymentStatusQueued ...
	ServiceDeploymentStatusQueued = "QUEUED"

	// ServiceDeploymentStatusRunning ...
	ServiceDeploymentStatusRunning = "RUNNING"

	// ServiceDeploymentStatusFailed ...
	ServiceDeploymentStatusFailed = "FAILED"

	// ServiceDeploymentStatusSuccess ...
	ServiceDeploymentStatusSuccess = "SUCCESS"
)
