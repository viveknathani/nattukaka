package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gomodule/redigo/redis"
	"github.com/viveknathani/nattukaka/app"
	"github.com/viveknathani/nattukaka/cache"
	"github.com/viveknathani/nattukaka/config"
	"github.com/viveknathani/nattukaka/database"
	"github.com/viveknathani/nattukaka/service"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var port string = ""
var databaseURL string = ""
var cacheURL string = ""
var jwtSecret string = ""

func init() {

	config.LoadEnvFromFile(".env")
	port = os.Getenv("PORT")
	databaseURL = os.Getenv("DATABASE_URL")
	cacheURL = os.Getenv("CACHE_URL")
	jwtSecret = os.Getenv("JWT_SECRET")
}

// getDatabase will init and return a db
func getDatabase() *database.Database {

	db := &database.Database{}
	err := db.Initialize(databaseURL)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	return db
}

// getCache will return a connection to Redis from the pool
func getCache() (*cache.Cache, redis.Conn) {

	memory := &cache.Cache{}
	memory.Initialize(cacheURL)
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

	// Failure of logger setup should prevent any further operation
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	return logger
}

func main() {
	// Make a channel to capture interrupts, SIGINT, SIGTERM
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Setup logger and database
	logger := getLogger()
	db := getDatabase()
	logger.Info("connected to database!")
	memory, memoryConn := getCache()
	logger.Info("connected to cache!")

	// Run database migrations
	// Read the schema.sql file
	schema, err := os.ReadFile("./database/schema.sql")
	if err != nil {
		logger.Fatal("Failed to read schema file: %v", zap.Error(err))
	}

	// Execute schema.sql
	_, err = db.Exec(string(schema))
	if err != nil {
		logger.Fatal("Failed to execute schema: %v", zap.Error(err))
	}
	logger.Info("ran database migrations!")

	// Init application
	fiberApp := fiber.New()
	application := &app.App{
		Fiber: fiberApp,
		Service: &service.Service{
			Db:        db,
			Cache:     memoryConn,
			JwtSecret: []byte(jwtSecret),
			Logger:    logger,
		},
	}
	application.SetupRoutes()

	// Start server
	go func() {
		logger.Info("starting server...")
		err = application.Fiber.Listen(":" + port)
		if err != nil {
			logger.Fatal("server cannot start!", zap.Error(err))
		}
	}()

	// Wait for signal
	<-done
	shutdown(logger, application, db, memory)
}

func shutdown(logger *zap.Logger, application *app.App, db *database.Database, memory *cache.Cache) {
	logger.Info("shutting down server...")

	err := application.Service.Cache.Close()
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

	application.Fiber.Shutdown()
	logger.Info("goodbye!")
}
