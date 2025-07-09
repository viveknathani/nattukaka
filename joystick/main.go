package main

import (
	"log"
	"os"
	"os/signal"
	"player/gen"
	"syscall"

	"joystick/database"
	"joystick/routes"
	"joystick/shared"

	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"
)

func loadEnviromentVariables() {
	if os.Getenv("ENV") == shared.EnvDevelopment {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}
}

func connectToDatabaseAndRunMigrations() *gorm.DB {
	db := database.NewDatabase(os.Getenv("DATABASE_URL"))
	connection, err := db.DB()
	if err != nil {
		log.Fatal("Error in getting the connection object!", err)
	}
	database.RunMigrations("database/migrations", connection)
	return db
}

func main() {
	grpcClient, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("failed to connect to player service: ", err)
	}
	defer grpcClient.Close()

	gen.NewPlayerServiceClient(grpcClient)
	// Use .env file in development, and the process enviroment in production
	loadEnviromentVariables()

	// Setup the database
	db := connectToDatabaseAndRunMigrations()

	// Setup cache
	options, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatal("failed to parse redis url: ", err)
	}
	redisClient := redis.NewClient(options)

	// Init the fiber app
	app := fiber.New(
		fiber.Config{
			BodyLimit: 10 * 1024 * 1024,
		},
	)

	// Define the shared state of the application
	state := shared.State{
		Logger:    shared.NewLogger(),
		Cache:     redisClient,
		Database:  db,
		Validator: validator.New(),
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	// Setup routes
	routes.SetupRoutes(app, &state)

	// Make a channel to listen for key events that can kill the process
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Put in the port, start the server!
	port := os.Getenv("PORT")
	go func() {
		err := app.Listen(":" + port)
		if err != nil {
			state.Logger.Error(err.Error())
			os.Exit(1)
		}
	}()
	state.Logger.Info("server is up! ⚡️")

	// Wait for the process to be killed
	<-done
	shutdown(&state)
}

func shutdown(state *shared.State) {

	sqlDB, err := state.Database.DB()
	if err != nil {
		state.Logger.Error(err.Error())
	} else {
		if err := sqlDB.Close(); err != nil {
			state.Logger.Error(err.Error())
		}
	}

	err = state.Cache.Close()
	if err != nil {
		state.Logger.Error(err.Error())
	}

	state.Logger.Info("goodbye!")
}
