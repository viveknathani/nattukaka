package app

// SetupRoutes will add all the endpoints to the fiber app
func (app App) SetupRoutes() {
	app.Fiber.Get("/", app.HomeController)
}
