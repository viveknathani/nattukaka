package app

// SetupRoutes will add all the endpoints to the fiber app
func (app *App) SetupRoutes() {
	app.Fiber.Get("/", app.HomeController)
	app.Fiber.Post("/api/v1/auth/signup", app.SignupController)
	app.Fiber.Post("/api/v1/auth/sendOTP", app.SendOTPController)
	app.Fiber.Post("/api/v1/auth/verifyOTP", app.VerifyOTPController)
	app.Fiber.Get("/api/v1/auth/me", app.verifyTokeMiddleware, app.GetProfileController)
}
