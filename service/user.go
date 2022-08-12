package service

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"time"
	"unicode"

	"github.com/golang-jwt/jwt"
	"github.com/viveknathani/nattukaka/cache"
	"github.com/viveknathani/nattukaka/entity"
	"golang.org/x/crypto/bcrypt"
)

const (
	ageOfToken = time.Hour * 24 * 2
)

// Signup creates a new user in the database if all the checks pass.
func (service *Service) Signup(ctx context.Context, u *entity.User) error {

	if u == nil {
		return ErrNilUser
	}

	_, err := mail.ParseAddress(u.Email)
	if err != nil {
		service.Logger.Error(err.Error(), zapReqID(ctx))
		return ErrInvalidEmailFormat
	}

	service.Logger.Info("database: email check start.", zapReqID(ctx))
	user, err := service.Repo.GetUser(u.Email)
	service.Logger.Info("database: email check complete.", zapReqID(ctx))
	if err != nil {
		service.Logger.Error(err.Error(), zapReqID(ctx))
		return err
	}

	if user != nil {
		return ErrEmailExists
	}

	if !isValidPassword(string(u.Password)) {
		return ErrInvalidPasswordFormat
	}

	hash, err := bcrypt.GenerateFromPassword(u.Password, bcrypt.DefaultCost)
	if err != nil {
		service.Logger.Error(err.Error(), zapReqID(ctx))
		return errors.New("bcrypt error, check logs")
	}

	u.Password = hash

	service.Logger.Info("database: insert user start.", zapReqID(ctx))
	err = service.Repo.CreateUser(u)
	service.Logger.Info("database: insert user complete.", zapReqID(ctx))
	if err != nil {
		service.Logger.Error(err.Error(), zapReqID(ctx))
		return err
	}
	return nil
}

// Login creates a new JWT and returns it if there is no error.
func (service *Service) Login(ctx context.Context, u *entity.User) (string, error) {

	if u == nil {
		return "", ErrNilUser
	}

	service.Logger.Info("database: email check start.", zapReqID(ctx))
	user, err := service.Repo.GetUser(u.Email)
	service.Logger.Info("database: email check complete.", zapReqID(ctx))
	if err != nil {
		service.Logger.Error(err.Error(), zapReqID(ctx))
		return "", err
	}
	if user == nil {
		return "", ErrInvalidEmailPassword
	}

	err = bcrypt.CompareHashAndPassword(user.Password, u.Password)
	if err != nil {
		return "", ErrInvalidEmailPassword
	}

	return service.createToken(user.Id)
}

// createToken will create a new JWT with id as payload and an expiry time
func (service *Service) createToken(id string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(ageOfToken).Unix(),
	})

	return token.SignedString(service.JwtSecret)
}

// VerifyAndDecodeToken will get the payload we need if the token is valid.
func (service *Service) VerifyAndDecodeToken(ctx context.Context, token string) (string, error) {

	service.Logger.Info("cache: token search start.", zapReqID(ctx))
	if service.isBlacklistedToken(ctx, token) {
		return "", ErrInvalidToken
	}
	service.Logger.Info("cache: token search complete.", zapReqID(ctx))

	parsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {

		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return service.JwtSecret, nil
	})

	if claims, ok := parsed.Claims.(jwt.MapClaims); ok && parsed.Valid {
		return claims["id"].(string), nil
	}

	service.Logger.Error(err.Error(), zapReqID(ctx))
	return "", ErrInvalidToken
}

// isValidPassword does a linear time check for password format.
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

// Logout will put the JWT in the cache which acts as a blacklist.
func (service *Service) Logout(ctx context.Context, token string) error {

	service.Logger.Info("cache: blacklist token start.", zapReqID(ctx))
	err := service.blacklistToken(token)
	if err != nil {
		service.Logger.Error(err.Error(), zapReqID(ctx))
		return err
	}
	service.Logger.Info("cache: blacklist token complete.", zapReqID(ctx))
	return nil
}

func (service *Service) blacklistToken(token string) error {
	_, err := cache.Set(service.Conn, token, []byte("true"))
	return err
}

func (service *Service) isBlacklistedToken(ctx context.Context, token string) bool {
	res, err := cache.Get(service.Conn, token)
	if err != nil {
		service.Logger.Error(err.Error(), zapReqID(ctx))
		return false
	}
	return string(res) == "true"
}
