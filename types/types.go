package types

import "time"

// SignupRequest represents the request body for user signup
type SignupRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// SendOTPRequest represents the request body for sending OTP
type SendOTPRequest struct {
	Email string `json:"email"`
}

// VerifyOTPRequest represents the request body for verifying OTP
type VerifyOTPRequest struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

// User represents a user profile
type User struct {
	ID       int    `json:"-"`
	PublicID string `json:"publicId,omitempty"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

// VerifyOTPResponse represents the response body for response of verify OTP
type VerifyOTPResponse struct {
	User      User   `json:"user"`
	AuthToken string `json:"authToken"`
}

// Workspace represents a workspace in the system
type Workspace struct {
	ID       int    `json:"-"`
	PublicID string `json:"publicId,omitempty"`
	Name     string `json:"name"`
}

// WorkspaceUser represents a user in a workspace with a role
type WorkspaceUser struct {
	ID          int    `json:"-"`
	PublicID    string `json:"publicId,omitempty"`
	WorkspaceID int    `json:"workspaceId"`
	UserID      int    `json:"userId"`
	Role        string `json:"role"`
}

// CreateWorkspaceRequest represents the request body for creating a workspace
type CreateWorkspaceRequest struct {
	Name string `json:"name"`
}

// InstanceType represents an instance type
type InstanceType struct {
	ID       int     `json:"id"`
	PublicID string  `json:"publicId"`
	Name     string  `json:"name"`
	CPU      float64 `json:"cpu"`
	Memory   float64 `json:"memory"`
	Disk     float64 `json:"disk"`
}

// CreateServiceRequest represents a request to create a new service
type CreateServiceRequest struct {
	Base     *Service         `json:"base"`
	Web      *WebService      `json:"web,omitempty"`
	Database *DatabaseService `json:"database,omitempty"`
}

// UpdateServiceRequest represents a request to update an existing service
type UpdateServiceRequest struct {
	Base     *Service         `json:"base,omitempty"`
	Web      *WebService      `json:"web,omitempty"`
	Database *DatabaseService `json:"database,omitempty"`
}

// Service represents a service entity
type Service struct {
	ID             int       `json:"-"`                        // internal ID, not exported
	PublicID       string    `json:"publicId"`                 // public ID of the service
	Name           string    `json:"name"`                     // name of the service
	Status         string    `json:"status"`                   // status of the service
	Type           string    `json:"type"`                     // type of the service (WEB or DATABASE)
	Runtime        string    `json:"runtime"`                  // runtime of the service
	WorkspaceID    int       `json:"workspaceId"`              // ID of the workspace the service belongs to
	CreatedBy      int       `json:"createdBy"`                // ID of the user who created the service
	LastDeployedAt time.Time `json:"lastDeployedAt,omitempty"` // last deployed timestamp
	CreatedAt      time.Time `json:"createdAt"`                // creation timestamp
	InstanceTypeID int       `json:"instanceTypeId"`           // ID of the instance type
	InternalURL    string    `json:"internalUrl,omitempty"`    // internal URL of the service
	ExternalURL    string    `json:"externalUrl,omitempty"`    // external URL of the service
}

// WebService represents a web service entity
type WebService struct {
	PublicID         string `json:"publicId"`         // public ID of the web service
	ServiceID        int    `json:"-"`                // internal service ID, not exported
	Repository       string `json:"repository"`       // repository URL
	Branch           string `json:"branch"`           // branch name
	RootDirectory    string `json:"rootDirectory"`    // root directory path
	BuildCommand     string `json:"buildCommand"`     // build command
	PreDeployCommand string `json:"preDeployCommand"` // pre-deploy command
	StartCommand     string `json:"startCommand"`     // start command
	HealthCheckPath  string `json:"healthCheckPath"`  // health check path
	Environment      string `json:"environment"`      // environment variables
}

// DatabaseService represents a database service entity
type DatabaseService struct {
	PublicID  string `json:"publicId"` // public ID of the database service
	ServiceID int    `json:"-"`        // internal service ID, not exported
}

// Volume represents a volume entry
type Volume struct {
	ID       int    `json:"id"`       // ID of the volume
	PublicID string `json:"publicId"` // Public ID of the volume
}

// CreateVolumeRequest represents the request payload for creating a volume
type CreateVolumeRequest struct {
	PublicID string `json:"publicId"` // Public ID of the volume
}

// AttachVolumeRequest represents the request payload for attaching a volume to a service
type AttachVolumeRequest struct {
	VolumeID int `json:"volumeId"` // ID of the volume
}

// DetachVolumeRequest represents the request payload for detaching a volume from a service
type DetachVolumeRequest struct {
	VolumeID int `json:"volumeId"` // ID of the volume
}

// VolumePaginationRequest represents the request payload for paginated volume retrieval
type VolumePaginationRequest struct {
	Page  int `json:"page"`  // Page number
	Limit int `json:"limit"` // Number of items per page
}
