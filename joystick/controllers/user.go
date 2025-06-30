package controllers

import (
	"joystick/services"
	"joystick/shared"

	"github.com/gofiber/fiber/v2"
)

// UserController provides methods for user management.
type UserController struct {
	userService *services.UserService
}

// NewUserController creates a new instance of UserController with the provided state.
func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// Signup creates a new user in the database if all the checks pass.
func (userController *UserController) Signup(c *fiber.Ctx) error {
	var request shared.UserSignupRequest

	if err := shared.ParseBody(c, &request); err != nil {
		return shared.SendStandardResponse(c, fiber.StatusBadRequest, nil, err.Error())
	}

	ok, err := userController.userService.Signup(request.Username, request.Password)
	if err != nil {
		return shared.SendStandardResponse(c, err.ErrorCode(), nil, err.Error())
	}

	if !ok {
		return shared.SendStandardResponse(
			c,
			fiber.StatusConflict,
			nil,
			"user with this username already exists",
		)
	}

	return shared.SendStandardResponse(c, fiber.StatusCreated, nil, "user created")
}

// Login creates a new JWT and returns it if there is no error.
func (userController *UserController) Login(c *fiber.Ctx) error {
	var request shared.UserLoginRequest

	if err := shared.ParseBody(c, &request); err != nil {
		return shared.SendStandardResponse(c, fiber.StatusBadRequest, nil, err.Error())
	}

	token, err := userController.userService.Login(request.Username, request.Password)
	if err != nil {
		return shared.SendStandardResponse(c, err.ErrorCode(), nil, err.Error())
	}

	return shared.SendStandardResponse(
		c,
		fiber.StatusOK,
		&map[string]interface{}{"token": *token},
		"login successful",
	)
}
