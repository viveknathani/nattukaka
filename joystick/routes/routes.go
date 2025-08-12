package routes

import (
	"joystick/controllers"
	"joystick/services"
	"joystick/shared"
	"log"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes sets up the routes for the application.
func SetupRoutes(app *fiber.App, state *shared.State) {
	// Create player client - using localhost:50051 for now
	playerClient, err := services.NewPlayerClient("localhost:50051")
	if err != nil {
		state.Logger.Error("Failed to create player client: " + err.Error())
		log.Fatal("failed to create player client: " + err.Error())
		return
	}
	if playerClient == nil {
		state.Logger.Error("Player client is nil, continuing without it")
	}

	userService := services.NewUserService(state)
	serviceDeploymentService := services.NewServiceDeploymentService(state, playerClient)
	serviceService := services.NewServiceService(state, serviceDeploymentService)

	userController := controllers.NewUserController(userService)
	serviceController := controllers.NewServiceController(serviceService)
	serviceDeploymentController := controllers.NewServiceDeploymentController(
		serviceDeploymentService,
	)

	authMiddleware := services.GetAuthMiddleware(userService)

	app.Post("/api/v1/users/signup", userController.Signup)
	app.Post("/api/v1/users/login", userController.Login)
	app.Get("/api/v1/services", authMiddleware, serviceController.GetAllServices)
	app.Post("/api/v1/services", authMiddleware, serviceController.CreateService)
	app.Get("/api/v1/services/:serviceID", authMiddleware, serviceController.GetService)
	app.Put("/api/v1/services/:serviceID", authMiddleware, serviceController.UpdateService)
	app.Delete("/api/v1/services/:serviceID", authMiddleware, serviceController.DeleteService)
	app.Post(
		"/api/v1/services/:serviceID/deployments",
		authMiddleware,
		serviceDeploymentController.CreateServiceDeployment,
	)
	app.Get(
		"/api/v1/services/:serviceID/deployments",
		authMiddleware,
		serviceDeploymentController.GetAllServiceDeployments,
	)
	app.Patch(
		"/api/v1/services/:serviceID/deployments/:deploymentID",
		authMiddleware,
		serviceDeploymentController.UpdateServiceDeployment,
	)
}
