package utils

import "time"

const (
	// PrefixRedisKeyUserOTP ...
	PrefixRedisKeyUserOTP = "OTP:"

	// AgeOfToken ...
	AgeOfToken = time.Hour * 24 * 2

	// PublicIDLength ...
	PublicIDLength = 10

	// UserRoleAdmin
	UserRoleAdmin = "ADMIN"
)
