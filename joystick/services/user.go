package services

import (
	"fmt"
	"joystick/shared"
	"os"
	"time"
	"unicode"

	"golang.org/x/crypto/bcrypt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

// UserService provides methods for user authentication and authorization.
type UserService struct {
	state     *shared.State
	jwtSecret []byte
}

// NewUserService creates a new instance of UserService with the provided state.
func NewUserService(state *shared.State) *UserService {

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		state.Logger.Error("JWT_SECRET environment variable is not set")
		os.Exit(1)
	}

	return &UserService{
		state:     state,
		jwtSecret: []byte(jwtSecret),
	}
}

// Signup creates a new user in the database if all the checks pass.
func (userService *UserService) Signup(username, password string) (bool, *shared.JoyStickError) {
	if len(username) < 3 {
		return false, shared.ErrUsernameTooShort
	}

	if !isValidPassword(password) {
		return false, shared.ErrInvalidPasswordFormat
	}

	var user shared.User

	err := userService.state.Database.Table("users").
		Where("username = ?", username).
		First(&user).
		Error
	if err == nil && user.Username != "" {
		return false, shared.ErrUserAlreadyExists
	}

	if err != nil && err.Error() != "record not found" {
		userService.state.Logger.Error("error checking for existing user" + err.Error())
		return false, shared.ErrInternalServerError
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		userService.state.Logger.Error(err.Error())
		return false, shared.ErrInternalServerError
	}

	err = userService.state.Database.Table("users").Create(&shared.User{
		Username: username,
		Password: string(hash),
	}).Error

	if err != nil {
		userService.state.Logger.Error("error creating user: " + err.Error())
		return false, shared.ErrInternalServerError
	}

	return true, nil
}

// Login creates a new JWT and returns it if there is no error.
func (userService *UserService) Login(username, password string) (*string, *shared.JoyStickError) {
	var user shared.User

	err := userService.state.Database.Table("users").
		Where("username = ?", username).
		First(&user).
		Error
	if err != nil {
		if err.Error() == "record not found" {
			return nil, shared.ErrUserNotFound
		}
		userService.state.Logger.Error("error fetching user: " + err.Error())
		return nil, shared.ErrInternalServerError
	}

	if user.Username == "" {
		return nil, shared.ErrUserNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, shared.ErrInvalidCredentials
	}

	return userService.createToken(user.ID)
}

func isValidPassword(password string) bool {

	const minLength = 8
	length := 0

	hasNumber := false
	hasUppercase := false
	hasLowercase := false
	hasSpecial := false

	for _, c := range password {

		switch {
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsUpper(c):
			hasUppercase = true
		case unicode.IsLower(c):
			hasLowercase = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		default:
			return false
		}

		length++
	}

	return length >= minLength && hasNumber && hasLowercase && hasUppercase && hasSpecial
}

func (userService *UserService) createToken(id int) (*string, *shared.JoyStickError) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(shared.TokenAge).Unix(),
	})

	tokenStr, err := token.SignedString(userService.jwtSecret)

	if err != nil {
		userService.state.Logger.Error(err.Error())
		return nil, shared.ErrInternalServerError
	}

	return &tokenStr, nil
}

// VerifyToken verifies a JWT and returns the user ID if it is valid.
func (userService *UserService) VerifyToken(token string) (int, *shared.JoyStickError) {

	parsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {

		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return userService.jwtSecret, nil
	})

	if err != nil {
		return 0, shared.ErrInvalidToken
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return 0, shared.ErrInvalidCredentials
	}

	return int(claims["id"].(float64)), nil
}

// GetAuthMiddleware returns a middleware function that checks for a valid JWT.
func GetAuthMiddleware(userService *UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return shared.SendStandardResponse(c, fiber.StatusUnauthorized, nil, "unauthorized")
		}

		token := authHeader[len("Bearer "):]

		userID, err := userService.VerifyToken(token)
		if err != nil {
			return shared.SendStandardResponse(c, fiber.StatusUnauthorized, nil, "unauthorized")
		}

		c.Locals("userID", userID)
		return c.Next()
	}
}
