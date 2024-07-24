package app

// SetupRoutes will add all the endpoints to the fiber app
func (app *App) SetupRoutes() {
	app.Fiber.Get("/", app.HomeController)
	app.Fiber.Post("/api/v1/auth/signup", app.SignupController)
	app.Fiber.Post("/api/v1/auth/sendOTP", app.SendOTPController)
	app.Fiber.Post("/api/v1/auth/verifyOTP", app.VerifyOTPController)
	app.Fiber.Get("/api/v1/users/whoami", app.verifyTokeMiddleware, app.GetProfileController)
	app.Fiber.Get("/api/v1/users/workspaces", app.verifyTokeMiddleware, app.GetMyWorkspacesController)
	app.Fiber.Post("/api/v1/workspaces", app.verifyTokeMiddleware, app.CreateWorkspaceController)
	app.Fiber.Get("/api/v1/workspaces/:workspaceID/users", app.verifyTokeMiddleware, app.GetWorkspaceUsersController)
	app.Fiber.Delete("/api/v1/workspaces/:workspaceID", app.verifyTokeMiddleware, app.DeleteWorkspaceController)
	app.Fiber.Get("/api/v1/instanceTypes", app.GetAllInstanceTypesController)
	app.Fiber.Post("/api/v1/services", app.verifyTokeMiddleware, app.CreateServiceController)
	app.Fiber.Get("/api/v1/services/workspace/:workspaceID", app.verifyTokeMiddleware, app.GetServicesByWorkspaceController)
	app.Fiber.Get("/api/v1/services/:serviceID", app.verifyTokeMiddleware, app.GetServiceByIDController)
	app.Fiber.Put("/api/v1/services/:serviceID", app.verifyTokeMiddleware, app.UpdateServiceController)
	app.Fiber.Delete("/api/v1/services/:serviceID", app.verifyTokeMiddleware, app.DeleteServiceController)
	app.Fiber.Post("/api/v1/volumes", app.verifyTokeMiddleware, app.CreateVolumeController)
	app.Fiber.Get("/api/v1/volumes/workspace/:workspaceID", app.verifyTokeMiddleware, app.GetVolumesByWorkspaceController)
	app.Fiber.Get("/api/v1/volumes/:volumeID", app.verifyTokeMiddleware, app.GetVolumeByIDController)
	app.Fiber.Post("/api/v1/volumes/attach/:serviceID", app.verifyTokeMiddleware, app.AttachVolumeToServiceController)
	app.Fiber.Post("/api/v1/volumes/detach/:serviceID", app.verifyTokeMiddleware, app.DeleteServiceController)
	app.Fiber.Delete("/api/v1/volumes/:volumeID", app.verifyTokeMiddleware, app.DeleteServiceController)
}
