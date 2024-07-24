package types

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
	ID       int    `json:"id,omitempty"`
	PublicID string `json:"id,omitempty"`
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
	ID       int    `json:"id,omitempty"`
	PublicID string `json:"publicId,omitempty"`
	Name     string `json:"name"`
}

// WorkspaceUser represents a user in a workspace with a role
type WorkspaceUser struct {
	ID          int    `json:"id,omitempty"`
	PublicID    string `json:"publicId,omitempty"`
	WorkspaceID int    `json:"workspaceId"`
	UserID      int    `json:"userId"`
	Role        string `json:"role"`
}

// CreateWorkspaceRequest represents the request body for creating a workspace
type CreateWorkspaceRequest struct {
	Name string `json:"name"`
}
