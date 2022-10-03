package service

import (
	"context"

	"github.com/gomodule/redigo/redis"
	"github.com/viveknathani/nattukaka/database"
	"github.com/viveknathani/nattukaka/shared"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Service struct {
	Repo      *database.Database
	Conn      redis.Conn
	JwtSecret []byte
	Logger    *zap.Logger
}

func zapReqID(ctx context.Context) zapcore.Field {

	return zapcore.Field{
		Key:    "requestID",
		String: shared.ExtractRequestID(ctx),
		Type:   zapcore.StringType,
	}
}
