package service

import (
	"github.com/gomodule/redigo/redis"
	"github.com/viveknathani/nattukaka/database"
	"go.uber.org/zap"
)

// Service encapsulates the nattukaka service
type Service struct {
	Db        *database.Database
	Cache     redis.Conn
	JwtSecret []byte
	Logger    *zap.Logger
}
