package app

import "github.com/gofiber/fiber/v2"

// HomeController ...
func (app *App) HomeController(c *fiber.Ctx) error {
	return c.SendFile("./index.html")
}
