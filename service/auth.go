package service

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/viveknathani/nattukaka/cache"
	"github.com/viveknathani/nattukaka/types"
	"github.com/viveknathani/nattukaka/utils"
)

const (
	ageOfToken = time.Hour * 24 * 2
)

// Signup creates a user in the database
func (srv *Service) Signup(req *types.SignupRequest) (int, string) {
	if req == nil {
		return fiber.StatusInternalServerError, "nil check fail"
	}

	existingUser, err := srv.Db.GetUserByEmail(req.Email)
	if err != nil {
		srv.Logger.Error(err.Error())
		return fiber.StatusInternalServerError, "something went wrong"
	}

	if existingUser != nil {
		return fiber.StatusBadRequest, "user exists"
	}

	if !utils.IsValidEmail(req.Email) {
		return fiber.StatusBadRequest, "invalid email"
	}

	publicID, err := utils.GeneratePublicId("user")
	if err != nil {
		srv.Logger.Error(err.Error())
		return fiber.StatusInternalServerError, "something went wrong"
	}

	err = srv.Db.InsertUser(&types.User{
		Name:     req.Name,
		Email:    req.Email,
		PublicID: publicID,
	})
	if err != nil {
		srv.Logger.Error(err.Error())
		return fiber.StatusInternalServerError, "something went wrong!"
	}
	return fiber.StatusCreated, "created!"
}

// SendOTP sends an OTP on your email!
func (srv *Service) SendOTP(request *types.SendOTPRequest) (int, string) {
	if request == nil {
		return fiber.StatusInternalServerError, "nil check fail"
	}

	// Pull existing user
	existingUser, err := srv.Db.GetUserByEmail(request.Email)
	if err != nil {
		srv.Logger.Error(err.Error())
		return fiber.StatusInternalServerError, "something went wrong"
	}
	if existingUser == nil {
		return fiber.StatusBadRequest, "user does not exist"
	}

	// Generate and store OTP in Redis
	otp := strconv.Itoa(utils.GenerateRandomNumber(1000, 9999))
	_, err = cache.Set(srv.Cache, prefixRedisKeyUserOTP+strconv.Itoa((existingUser.ID)), []byte(otp))
	if err != nil {
		srv.Logger.Error(err.Error())
		return fiber.StatusInternalServerError, "something went wrong"
	}

	// Send OTP via email
	err = utils.SendEmail(existingUser.Email, "Your OTP is here!", "Your OTP for logging into nattukaka is: "+otp)
	if err != nil {
		srv.Logger.Error(err.Error())
		return fiber.StatusInternalServerError, "something went wrong"
	}

	return fiber.StatusOK, "OTP sent"
}

// createToken will create a new JWT with id as payload and an expiry time
func (srv *Service) createToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(ageOfToken).Unix(),
	})

	return token.SignedString(srv.JwtSecret)
}

// VerifyOTP verifies your OTP!
func (srv *Service) VerifyOTP(request *types.VerifyOTPRequest) (int, string, *types.VerifyOTPResponse) {
	if request == nil {
		return fiber.StatusInternalServerError, "nil check fail", nil
	}

	// Pull existing user
	existingUser, err := srv.Db.GetUserByEmail(request.Email)
	if err != nil {
		srv.Logger.Error(err.Error())
		return fiber.StatusInternalServerError, "something went wrong", nil
	}
	if existingUser == nil {
		return fiber.StatusBadRequest, "user does not exist", nil
	}

	correctOTPAsBytes, err := cache.Get(srv.Cache, prefixRedisKeyUserOTP+strconv.Itoa((existingUser.ID)))
	if err != nil {
		srv.Logger.Error(err.Error())
		return fiber.StatusInternalServerError, "something went wrong", nil
	}

	if string(correctOTPAsBytes) != request.OTP {
		return fiber.StatusBadRequest, "incorrect OTP", nil
	}

	token, err := srv.createToken(existingUser.Email)
	if err != nil {
		srv.Logger.Error(err.Error())
		return fiber.StatusInternalServerError, "something went wrong", nil
	}

	return fiber.StatusOK, "OTP verified", &types.VerifyOTPResponse{
		User:      *existingUser,
		AuthToken: token,
	}
}

// GetProfile returns your profile data by email
func (srv *Service) GetProfile(email string) (int, string, *types.User) {
	// Pull existing user
	existingUser, err := srv.Db.GetUserByEmail(email)
	if err != nil {
		srv.Logger.Error(err.Error())
		return fiber.StatusInternalServerError, "something went wrong", nil
	}

	if existingUser == nil {
		return fiber.StatusBadRequest, "user does not exist", nil
	}

	return fiber.StatusOK, "", existingUser
}
