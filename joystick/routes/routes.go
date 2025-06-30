package routes

import (
	"joystick/controllers"
	"joystick/services"
	"joystick/shared"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes sets up the routes for the application.
func SetupRoutes(app *fiber.App, state *shared.State) {
	userService := services.NewUserService(state)
	serviceService := services.NewServiceService(state)
	userController := controllers.NewUserController(userService)
	serviceController := controllers.NewServiceController(serviceService)

	authMiddleware := services.GetAuthMiddleware(userService)

	app.Post("/api/v1/users/signup", userController.Signup)
	app.Post("/api/v1/users/login", userController.Login)
	app.Get("/api/v1/services", authMiddleware, serviceController.GetAllServices)
	app.Post("/api/v1/services", authMiddleware, serviceController.CreateService)
	app.Get("/api/v1/services/:serviceID", authMiddleware, serviceController.GetService)
	app.Put("/api/v1/services/:serviceID", authMiddleware, serviceController.UpdateService)
	app.Delete("/api/v1/services/:serviceID", authMiddleware, serviceController.DeleteService)
}
