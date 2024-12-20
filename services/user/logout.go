package user

import (
	"gin-auth-mongo/models/requests"
	"gin-auth-mongo/repositories"
)

func UserLogoutCurrentDevice(userID string, request *requests.LogoutRequest) error {
	return repositories.DeleteRefreshTokenByUserIDAndDevice(userID, request.Device)
}

func UserLogoutAllDevice(userID string) error {
	return repositories.DeleteRefreshTokenByUserID(userID)
}
