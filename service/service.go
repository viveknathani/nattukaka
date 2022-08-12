package service

import (
	"github.com/gomodule/redigo/redis"
	"github.com/viveknathani/nattukaka/repository"
)

type Service struct {
	Repo repository.Repository
	Conn redis.Conn
}
