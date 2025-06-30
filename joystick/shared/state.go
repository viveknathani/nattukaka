package shared

import (
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// State represents the entire state of the application.
// It will hold all the commonly used variables together.
// We can pass this around for usage by controllers and services.
type State struct {
	Logger    *zap.Logger
	Cache     *redis.Client
	Database  *gorm.DB
	Validator *validator.Validate
}
