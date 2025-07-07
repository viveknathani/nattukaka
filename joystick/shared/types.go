package shared

import (
	"encoding/json"
	"time"
)

// User represents a user in the system.
type User struct {
	ID       int    `json:"id"       gorm:"primaryKey"`
	UUID     string `json:"uuid"     gorm:"type:uuid;default:uuid_generate_v4()"`
	Username string `json:"username" gorm:"not null;unique"`
	Password string `json:"password" gorm:"bytea;not null"`
}

// UserSignupRequest represents a user signup request.
type UserSignupRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// UserLoginRequest represents a user login request.
type UserLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// ServiceCreateRequest represents a service creation request.
type ServiceCreateRequest struct {
	Name          string          `json:"name"          validate:"required"`
	RepositoryURL string          `json:"repositoryUrl" validate:"required"`
	Branch        string          `json:"branch"        validate:"required"`
	EnvVars       json.RawMessage `json:"envVars"       validate:"required"`
	PortMapping   json.RawMessage `json:"portMapping"   validate:"required"`
}

// ServiceUpdateRequest represents a service update request.
type ServiceUpdateRequest struct {
	Name          string          `json:"name"          validate:"required"`
	RepositoryURL string          `json:"repositoryUrl" validate:"required"`
	Branch        string          `json:"branch"        validate:"required"`
	EnvVars       json.RawMessage `json:"envVars"       validate:"required"`
	PortMapping   json.RawMessage `json:"portMapping"   validate:"required"`
}

// PortMapping represents a port mapping in the system.
type PortMapping struct {
	ContainerPort int `json:"containerPort"`
	HostPort      int `json:"hostPort"`
}

// Service represents a service in the system.
type Service struct {
	ID            int             `json:"id"            gorm:"id;primaryKey"`
	UUID          string          `json:"uuid"          gorm:"type:uuid;default:uuid_generate_v4()"`
	Name          string          `json:"name"          gorm:"name;not null;unique"`
	RepositoryURL string          `json:"repositoryUrl" gorm:"repository_url;not null"`
	Branch        string          `json:"branch"        gorm:"branch;not null"`
	EnvVars       json.RawMessage `json:"envVars"       gorm:"env_vars;not null"`
	PortMapping   json.RawMessage `json:"portMapping"   gorm:"port_mapping;not null"`
	OwnerID       int             `json:"ownerId"       gorm:"owner_id;not null"`
	CreatedAt     time.Time       `json:"createdAt"     gorm:"created_at;not null"`
	UpdatedAt     time.Time       `json:"updatedAt"     gorm:"updated_at;not null"`
}

// ServiceDeployment represents a service deployment in the system.
type ServiceDeployment struct {
	ID          int       `json:"id"          gorm:"id;primaryKey"`
	UUID        string    `json:"uuid"        gorm:"type:uuid;default:uuid_generate_v4()"`
	ServiceID   int       `json:"serviceId"   gorm:"service_id;not null"`
	Commit      string    `json:"commit"      gorm:"commit;not null"`
	Status      string    `json:"status"      gorm:"status;not null"`
	ContainerID string    `json:"containerId" gorm:"container_id;not null"`
	NodeID      int       `json:"nodeId"      gorm:"node_id;not null"`
	CreatedAt   time.Time `json:"createdAt"   gorm:"created_at;not null"`
	UpdatedAt   time.Time `json:"updatedAt"   gorm:"updated_at;not null"`
}

// ServiceDeploymentUpdateRequest represents a service deployment update request.
// This body can accept status and containerId as optional fields.
// If not provided, won't update those fields.
type ServiceDeploymentUpdateRequest struct {
	Status      string `json:"status"`
	ContainerID string `json:"containerId"`
}

// Node represents a node in the system.
type Node struct {
	ID        int       `json:"id"          gorm:"id;primaryKey"`
	UUID      string    `json:"uuid"        gorm:"type:uuid;default:uuid_generate_v4()"`
	Name      string    `json:"name"        gorm:"name;not null;unique"`
	IP        string    `json:"ip"          gorm:"ip;not null"`
	CreatedAt time.Time `json:"createdAt"   gorm:"created_at;not null"`
	UpdatedAt time.Time `json:"updatedAt"   gorm:"updated_at;not null"`
}
