package app

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/viveknathani/nattukaka/types"
)

// SignupController ...
func (app *App) SignupController(c *fiber.Ctx) error {

	// Parse body
	var req types.SignupRequest
	if err := c.BodyParser(&req); err != nil {
		return sendResponse(c, fiber.StatusBadRequest, "cannot parse json", nil)
	}

	code, message := app.Service.Signup(&req)

	return sendResponse(c, code, message, nil)
}

// SendOTPController ...
func (app *App) SendOTPController(c *fiber.Ctx) error {

	// Parse body
	var req types.SendOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return sendResponse(c, fiber.StatusBadRequest, "cannot parse json", nil)
	}

	code, message := app.Service.SendOTP(&req)

	return sendResponse(c, code, message, nil)
}

// VerifyOTPController ...
func (app *App) VerifyOTPController(c *fiber.Ctx) error {

	// Parse body
	var req types.VerifyOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return sendResponse(c, fiber.StatusBadRequest, "cannot parse json", nil)
	}

	code, message, data := app.Service.VerifyOTP(&req)

	return sendResponse(c, code, message, data)
}

func (app *App) verifyTokeMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return sendResponse(c, fiber.StatusUnauthorized, "unauthorized", nil)
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return sendResponse(c, fiber.StatusUnauthorized, "unauthorized", nil)
	}

	tokenString := parts[1]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.ErrUnauthorized
		}
		return app.Service.JwtSecret, nil
	})

	if err != nil || !token.Valid {
		return sendResponse(c, fiber.StatusUnauthorized, "unauthorized", nil)
	}

	// Store user information in context
	c.Locals("user", token.Claims)

	return c.Next()
}

// GetProfileController ...
func (app *App) GetProfileController(c *fiber.Ctx) error {
	email := c.Locals("user").(jwt.MapClaims)["email"].(string)
	code, message, data := app.Service.GetProfile(email)

	return sendResponse(c, code, message, data)
}
