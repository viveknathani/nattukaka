package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/viveknathani/nattukaka/cache"
	"github.com/viveknathani/nattukaka/database"
	"github.com/viveknathani/nattukaka/httpkaka"
	"github.com/viveknathani/nattukaka/service"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var port string = ""
var databaseServer string = ""
var redisServer string = ""
var jwtSecret string = ""

func init() {

	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }
	port = os.Getenv("PORT")
	databaseUser := os.Getenv("DATABASE_USER")
	databasePassword := os.Getenv("DATABASE_PASSWORD")
	databaseHost := os.Getenv("DATABASE_HOST")
	databasePort := os.Getenv("DATABASE_PORT")
	databaseName := os.Getenv("DATABASE_NAME")
	databaseServer = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", databaseUser, databasePassword, databaseHost, databasePort, databaseName)
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisServer = fmt.Sprintf("%s:%s", redisHost, redisPort)
	jwtSecret = os.Getenv("JWT_SECRET")
}

// getDatabase will init and return a db
func getDatabase() *database.Database {

	db := &database.Database{}
	err := db.Initialize(databaseServer)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	return db
}

// getCache will return a connection to Redis from the pool
func getCache() (*cache.Cache, redis.Conn) {

	memory := &cache.Cache{}
	memory.Initialize(redisServer, "", "")
	memoryConn := memory.Pool.Get()
	return memory, memoryConn
}

// getLogger will configure and return a uber/zap logger
func getLogger() *zap.Logger {

	cfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevel(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "message",
			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,
			TimeKey:     "ts",
			EncodeTime:  zapcore.EpochMillisTimeEncoder,
		},
	}

	logger, err := cfg.Build()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	return logger
}

func main() {

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	logger := getLogger()
	db := getDatabase()
	memory, memoryConn := getCache()
	srv := &httpkaka.Server{
		Service: &service.Service{
			Repo:      db,
			Conn:      memoryConn,
			JwtSecret: []byte(jwtSecret),
			Logger:    logger,
		},
		Router: mux.NewRouter(),
	}
	srv.SetupRoutes()
	go func() {
		err := http.ListenAndServe(":"+port, srv)
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
	}()
	<-done
	shutdown(srv, db, memory)
}

func shutdown(srv *httpkaka.Server, db *database.Database, memory *cache.Cache) {

	err := srv.Service.Conn.Close()
	if err != nil {
		fmt.Print(err)
	}
	err = memory.Close()
	if err != nil {
		fmt.Print(err)
	}
	err = db.Close()
	if err != nil {
		fmt.Print(err)
	}
}
