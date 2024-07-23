package service

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/viveknathani/nattukaka/cache"
	"github.com/viveknathani/nattukaka/types"
	"github.com/viveknathani/nattukaka/utils"
)

const (
	ageOfToken = time.Hour * 24 * 2
)

// Signup creates a user in the database
func (srv *Service) Signup(user *types.User) (int, string) {
	if user == nil {
		return httpCodeInternalServerError, "nil check fail!"
	}

	existingUser, err := srv.Db.GetUserByEmail(user.Email)
	if err != nil {
		srv.Logger.Error(err.Error())
		return httpCodeInternalServerError, "something went wrong!"
	}

	if existingUser != nil {
		return httpCodeBadRequest, "user exists!"
	}

	if !utils.IsValidEmail(user.Email) {
		return httpCodeBadRequest, "invalid email!"
	}

	err = srv.Db.InsertUser(user)
	if err != nil {
		srv.Logger.Error(err.Error())
		return httpCodeInternalServerError, "something went wrong!"
	}
	return httpCodeCreated, "created!"
}

// SendOTP sends an OTP on your email!
func (srv *Service) SendOTP(request *types.SendOTPRequest) (int, string) {
	if request == nil {
		return httpCodeInternalServerError, "nil check fail!"
	}

	// Pull existing user
	existingUser, err := srv.Db.GetUserByEmail(request.Email)
	if err != nil {
		srv.Logger.Error(err.Error())
		return httpCodeInternalServerError, "something went wrong!"
	}
	if existingUser == nil {
		return httpCodeBadRequest, "user does not exist!"
	}

	// Generate and store OTP in Redis
	otp := strconv.Itoa(utils.GenerateRandomNumber(1000, 9999))
	_, err = cache.Set(srv.Cache, prefixRedisKeyUserOTP+strconv.Itoa((existingUser.ID)), []byte(otp))
	if err != nil {
		srv.Logger.Error(err.Error())
		return httpCodeInternalServerError, "something went wrong!"
	}

	// Send OTP via email
	err = utils.SendEmail(existingUser.Email, "Your OTP is here!", "Your OTP for logging into nattukaka is: "+otp)
	if err != nil {
		srv.Logger.Error(err.Error())
		return httpCodeInternalServerError, "something went wrong!"
	}

	return httpCodeOk, "OTP sent!"
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
		return httpCodeInternalServerError, "nil check fail!", nil
	}

	// Pull existing user
	existingUser, err := srv.Db.GetUserByEmail(request.Email)
	if err != nil {
		srv.Logger.Error(err.Error())
		return httpCodeInternalServerError, "something went wrong!", nil
	}
	if existingUser == nil {
		return httpCodeBadRequest, "user does not exist!", nil
	}

	correctOTPAsBytes, err := cache.Get(srv.Cache, prefixRedisKeyUserOTP+strconv.Itoa((existingUser.ID)))
	if err != nil {
		srv.Logger.Error(err.Error())
		return httpCodeInternalServerError, "something went wrong!", nil
	}

	if string(correctOTPAsBytes) != request.OTP {
		return httpCodeBadRequest, "incorrect OTP!", nil
	}

	token, err := srv.createToken(existingUser.Email)
	if err != nil {
		srv.Logger.Error(err.Error())
		return httpCodeInternalServerError, "something went wrong!", nil
	}

	return httpCodeOk, "OTP verified!", &types.VerifyOTPResponse{
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
		return httpCodeInternalServerError, "something went wrong!", nil
	}

	if existingUser == nil {
		return httpCodeBadRequest, "user does not exist!", nil
	}

	return httpCodeOk, "", existingUser
}
