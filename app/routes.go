package app

import (
	"github.com/viveknathani/nattukaka/controllers"
)

// SetupRoutes will add all the endpoints to the fiber app
func (app App) SetupRoutes() {
	app.Fiber.Get("/", controllers.HomeController)
}
