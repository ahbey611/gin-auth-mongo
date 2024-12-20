package repositories

import (
	"gin-auth-mongo/databases"
	"gin-auth-mongo/models"
	"gin-auth-mongo/utils/consts"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var userTable = "user"

// UserRepository data access interface
type UserRepository interface {
	GetAllUsers() ([]models.User, error)
	GetUserByID(id string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	GetUserByEmailOrUsername(email, username string) (*models.User, error)
}

func GetAllUsers() ([]models.User, error) {
	var users []models.User
	return FindManyWithoutPagination(databases.GetMongoCollection(userTable), nil, nil, nil, &users)
}

func GetUserByID(id string) (*models.User, error) {
	idObject, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var user models.User
	return FindOne(databases.GetMongoCollection(userTable), bson.M{"_id": idObject}, nil, &user)
}

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	// user := models.User{}
	return FindOne(databases.GetMongoCollection(userTable), bson.M{"email": email}, nil, &user)
}

func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	return FindOne(databases.GetMongoCollection(userTable), bson.M{"username": username}, nil, &user)
}

func GetUserByEmailOrUsername(email, username string) (*models.User, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"email": email},
			{"username": username},
		},
	}
	var user models.User
	return FindOne(databases.GetMongoCollection(userTable), filter, nil, &user)
}

func CreateUser(email, username, password, nickname string) error {
	user := models.User{
		Username:         username,
		Email:            email,
		Password:         password,
		Nickname:         nickname,
		Avatar:           consts.DEFAULT_AVATAR,
		CreatedAt:        time.Now().Format(consts.DATETIME_NANO_FORMAT),
		UpdatedAt:        time.Now().Format(consts.DATETIME_NANO_FORMAT),
		Premium:          false,
		PremiumExpiredAt: "",
	}
	return InsertOne(databases.GetMongoCollection(userTable), &user)
}

// login
func LoginWithEmailPassword(email, password string) (*models.User, error) {
	var user models.User
	return FindOne(databases.GetMongoCollection(userTable), bson.M{"email": email, "password": password}, nil, &user)
}

func LoginWithUsernamePassword(username, password string) (*models.User, error) {
	var user models.User
	return FindOne(databases.GetMongoCollection(userTable), bson.M{"username": username, "password": password}, nil, &user)
}

func UpdateUserPasswordByID(userID string, password string) error {
	idObject, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	return UpdateOne(databases.GetMongoCollection(userTable), bson.M{"_id": idObject}, bson.M{"$set": bson.M{"password": password}})
}

func UpdateNicknameByID(userID string, nickname string) error {
	idObject, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	return UpdateOne(databases.GetMongoCollection(userTable), bson.M{"_id": idObject}, bson.M{"$set": bson.M{"nickname": nickname}})
}

func UpdateAvatarByID(userID string, avatar string) error {
	idObject, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	return UpdateOne(databases.GetMongoCollection(userTable), bson.M{"_id": idObject}, bson.M{"$set": bson.M{"avatar": avatar}})
}

func DeleteUserByID(userID string) error {
	idObject, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	return DeleteOne(databases.GetMongoCollection(userTable), bson.M{"_id": idObject})
}
