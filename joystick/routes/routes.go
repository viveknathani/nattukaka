package routes

import (
	"embed"
	"io/fs"
	"joystick/controllers"
	"joystick/services"
	"joystick/shared"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

// SetupRoutes sets up the routes for the application.
func SetupRoutes(app *fiber.App, state *shared.State, staticFiles embed.FS) {
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

	// Create filesystem from embedded files
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal("Failed to create sub filesystem for static files: " + err.Error())
	}

	// Static files from embedded filesystem  
	app.Use("/", filesystem.New(filesystem.Config{
		Root: http.FS(staticFS),
	}))

	// Route redirects to remove .html suffix
	app.Get("/login", func(c *fiber.Ctx) error {
		data, err := staticFiles.ReadFile("static/login.html")
		if err != nil {
			return err
		}
		c.Set("Content-Type", "text/html")
		return c.Send(data)
	})
	app.Get("/dashboard", func(c *fiber.Ctx) error {
		data, err := staticFiles.ReadFile("static/dashboard.html")
		if err != nil {
			return err
		}
		c.Set("Content-Type", "text/html")
		return c.Send(data)
	})

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
