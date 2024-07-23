package controllers

import (
	"github.com/gofiber/fiber/v2"
)

// HomeController ...
func HomeController(c *fiber.Ctx) error {
	return c.SendFile("./index.html")
}
