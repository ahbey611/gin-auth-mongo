package auth

import (
	"errors"
	"gin-auth-mongo/graph/model"
	"gin-auth-mongo/models"
	"gin-auth-mongo/models/requests"
	"gin-auth-mongo/repositories"

	"gin-auth-mongo/utils/crypto"
	"gin-auth-mongo/utils/jwt"
)

func UserEmailLoginWithPassword(request *requests.EmailLoginWithPasswordRequest) (*models.User, *model.Token, error) {

	user, err := repositories.GetUserByEmail(request.Email)

	// user not found
	if err != nil || user == nil {
		return nil, nil, errors.New("incorrect email or password")
	}

	// check if the password is correct
	match, err := crypto.VerifyPassword(request.Password, user.Password)
	if err != nil || !match {
		return nil, nil, errors.New("incorrect email or password")
	}

	token, err := jwt.HandleLogin(user, request.Device)
	if err != nil {
		return nil, nil, err
	}

	return user, token, nil
}

func UserUsernameLoginWithPassword(request *requests.UsernameLoginWithPasswordRequest) (*models.User, *model.Token, error) {
	user, err := repositories.GetUserByUsername(request.Username)

	// user not found
	if err != nil || user == nil {
		return nil, nil, errors.New("invalid username or password")
	}

	// check if the password is correct
	match, err := crypto.VerifyPassword(request.Password, user.Password)
	if err != nil || !match {
		return nil, nil, errors.New("invalid username or password")
	}

	token, err := jwt.HandleLogin(user, request.Device)
	if err != nil {
		return nil, nil, err
	}

	return user, token, nil
}
