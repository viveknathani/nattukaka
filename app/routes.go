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
}
