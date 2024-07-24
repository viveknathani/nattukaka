package utils

import "time"

const (
	// PrefixRedisKeyUserOTP ...
	PrefixRedisKeyUserOTP = "OTP:"

	// AgeOfToken ...
	AgeOfToken = time.Hour * 24 * 2

	// PublicIDLength ...
	PublicIDLength = 10

	// UserRoleAdmin ...
	UserRoleAdmin = "ADMIN"

	// PageSize ...
	PageSize = 10

	// ServiceStatusCreated - when you have just created a service
	ServiceStatusCreated = "CREATED"

	// ServiceStatusDeploying - when your web service deployment is in progress
	ServiceStatusDeploying = "DEPLOYING"

	// ServiceStatusDeployed - when your web service is deployed
	ServiceStatusDeployed = "DEPLOYED"

	// ServiceStatusProvisioning - when your database service is being provisioned
	ServiceStatusProvisioning = "PROVISIONING"

	// ServiceStatusAvailable - when your service is available - used for database services
	ServiceStatusAvailable = "AVAILABLE"

	// ServiceStatusUnavailable - when your service does not start
	ServiceStatusUnavailable = "UNAVAILABLE"
)
