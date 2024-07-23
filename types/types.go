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
	PublicID string `json:"publicId,omitempty"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

// VerifyOTPResponse represents the response body for response of verify OTP
type VerifyOTPResponse struct {
	User      User   `json:"user"`
	AuthToken string `json:"authToken"`
}
