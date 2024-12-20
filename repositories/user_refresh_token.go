package repositories

import (
	"gin-auth-mongo/databases"
	"gin-auth-mongo/models"
	"gin-auth-mongo/utils/consts"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var userRefreshTokenTable = "user_refresh_token"

type RefreshTokenRepository interface {
	CreateRefreshToken(userID string, token string) error
	GetRefreshTokenByToken(token string) (*models.UserRefreshToken, error)
}

func CreateRefreshToken(userID string, token string, expiredAt time.Time, device string) error {
	userIDObject, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	refreshToken := &models.UserRefreshToken{
		UserID:    userIDObject,
		Token:     token,
		ExpiredAt: expiredAt.Format(consts.DATETIME_NANO_FORMAT),
		Device:    device,
	}
	return InsertOne(databases.GetMongoCollection(userRefreshTokenTable), refreshToken)
}

func GetRefreshTokenByToken(token string) (*models.UserRefreshToken, error) {
	var refreshToken models.UserRefreshToken
	return FindOne(databases.GetMongoCollection(userRefreshTokenTable), bson.M{"token": token}, nil, &refreshToken)
}

func DeleteRefreshTokenByUserIDAndDevice(userID string, device string) error {
	idObject, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	return DeleteMany(databases.GetMongoCollection(userRefreshTokenTable), bson.M{"user_id": idObject, "device": device})
}

func DeleteRefreshTokenByUserID(userID string) error {
	idObject, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	return DeleteMany(databases.GetMongoCollection(userRefreshTokenTable), bson.M{"user_id": idObject})
}

func DeleteRefreshTokenByToken(token string) error {
	return DeleteOne(databases.GetMongoCollection(userRefreshTokenTable), bson.M{"token": token})
}
