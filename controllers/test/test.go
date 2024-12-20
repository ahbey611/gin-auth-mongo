package test

// test mongodb

import (
	"context"
	"gin-auth-mongo/databases"
	"gin-auth-mongo/models"
	"gin-auth-mongo/repositories"
	"gin-auth-mongo/utils/consts"
	"gin-auth-mongo/utils/crypto"
	"gin-auth-mongo/utils/jwt"
	"gin-auth-mongo/utils/response"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InsertTestData(c *gin.Context) {

	// user2
	user2Collection := databases.MongoClient.Database(consts.MONGO_DATABASE).Collection("user2")
	time := time.Now().Format(consts.DATETIME_NANO_FORMAT)
	randomName, err := jwt.GenerateRefreshToken(5)
	if err != nil {
		log.Fatalf("Failed to generate random name: %v", err)
		response.InternalServerError(c)
		return
	}
	randomAge := rand.Intn(100)
	user := map[string]interface{}{"name": randomName, "age": randomAge, "created_at": time, "updated_at": time}
	_, err = user2Collection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Fatalf("Failed to insert user: %v", err)
		response.InternalServerError(c)
		return
	}
	response.Success(c)
}

func GetTestData(c *gin.Context) {
	user2Collection := databases.GetMongoCollection("user2")
	cursor, err := user2Collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatalf("Failed to get user: %v", err)
		response.InternalServerError(c)
		return
	}
	defer cursor.Close(context.TODO())

	var users []bson.M
	if err = cursor.All(context.TODO(), &users); err != nil {
		log.Fatalf("Failed to decode users: %v", err)
		response.InternalServerError(c)
		return
	}
	response.SuccessWithData(c, users)
}

func GetTestData2(c *gin.Context) {
	// get user by id
	id := "674732a2c5e731e6f4a9fd22"
	idObject, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatalf("Failed to convert id to object id: %v", err)
		response.InternalServerError(c)
		return
	}
	userCollection := databases.GetMongoCollection("user")
	user := User2{}
	// err = userCollection.FindOne(context.TODO(), bson.M{"_id": idObject}).Decode(&user)

	// hide password
	projection := bson.D{
		{Key: "password", Value: 0},
	}
	err = userCollection.FindOne(context.TODO(), bson.M{"_id": idObject}, options.FindOne().SetProjection(projection)).Decode(&user)
	// err := userCollection.FindOne(context.TODO(), bson.M{"username": "admin"}).Decode(&user)
	if err != nil {
		// log.Fatalf("Failed to get user: %v", err)
		log.Println("err: ", err)
		response.InternalServerError(c)
		return
	}
	response.SuccessWithData(c, user)
}

// get user by id without fix struct and hide password
func GetTestData3(c *gin.Context) {

	id := "674732a2c5e731e6f4a9fd22"
	idObject, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatalf("Failed to convert id to object id: %v", err)
		response.InternalServerError(c)
		return
	}
	userCollection := databases.GetMongoCollection("user")

	// hide password
	projection := bson.D{
		{Key: "password", Value: 0},
	}
	var user bson.M
	err = userCollection.FindOne(context.TODO(), bson.M{"_id": idObject}, options.FindOne().SetProjection(projection)).Decode(&user)

	if err != nil {
		// log.Fatalf("Failed to get user: %v", err)
		if err == mongo.ErrNoDocuments {
			response.BadRequestWithMessage(c, "User not found")
			return
		}
		log.Println("err: ", err)
		response.InternalServerError(c)
		return
	}
	response.SuccessWithData(c, user)
}

// Test find many without pagination
func GetTestData4(c *gin.Context) {
	var users []bson.M
	_, err := repositories.FindManyWithoutPagination(databases.GetMongoCollection("user2"), nil, nil, nil, &users)
	if err != nil {
		response.BadRequestWithMessage(c, err.Error())
		return
	}
	response.SuccessWithData(c, users)
}

// Test find many with pagination
func GetTestData5(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")
	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		response.BadRequestWithMessage(c, "Invalid page")
		return
	}
	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 64)
	if err != nil {
		response.BadRequestWithMessage(c, "Invalid page size")
		return
	}
	var users []bson.M
	_, err = repositories.FindMany(databases.GetMongoCollection("user2"), nil, nil, nil, page, pageSize, &users)
	if err != nil {
		response.BadRequestWithMessage(c, err.Error())
		return
	}
	response.SuccessWithData(c, users)
}

// Test find many with pagination and sort
func GetTestData6(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")
	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		response.BadRequestWithMessage(c, "Invalid page")
		return
	}
	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 64)
	if err != nil {
		response.BadRequestWithMessage(c, "Invalid page size")
		return
	}
	sortStr := c.DefaultQuery("sort", "created_at")
	sort := bson.D{{Key: sortStr, Value: -1}}
	var users []bson.M
	result, err := repositories.FindMany(databases.GetMongoCollection("user"), nil, sort, nil, page, pageSize, &users)
	if err != nil {
		response.BadRequestWithMessage(c, err.Error())
		return
	}
	// response.SuccessWithData(c, users)
	response.SuccessWithData(c, result)
}

// Test find many with pagination, sort and projection
func GetTestData7(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")
	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		response.BadRequestWithMessage(c, "Invalid page")
		return
	}
	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 64)
	if err != nil {
		response.BadRequestWithMessage(c, "Invalid page size")
		return
	}
	projectionStr := c.DefaultQuery("projection", "")
	projection := bson.D{}
	if projectionStr != "" {
		projection = bson.D{{Key: projectionStr, Value: 1}}
	}
	var users []bson.M
	_, err = repositories.FindMany(databases.GetMongoCollection("user2"), bson.D{{Key: "age", Value: 91}}, nil, projection, page, pageSize, &users)
	if err != nil {
		response.BadRequestWithMessage(c, err.Error())
		return
	}
	response.SuccessWithData(c, users)
}

// Test find many fixed struct
func GetTestData8(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")
	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		response.BadRequestWithMessage(c, "Invalid page")
		return
	}
	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 64)
	if err != nil {
		response.BadRequestWithMessage(c, "Invalid page size")
		return
	}
	var users []models.User
	_, err = repositories.FindMany(databases.GetMongoCollection("user"), nil, nil, nil, page, pageSize, &users)
	if err != nil {
		response.BadRequestWithMessage(c, err.Error())
		return
	}
	response.SuccessWithData(c, users)
}

type User2 struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username  string             `bson:"username" json:"username"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"-"`
	Nickname  string             `bson:"nickname" json:"nickname"`
	Avatar    string             `bson:"avatar" json:"avatar"`
	CreatedAt string             `bson:"created_at" json:"created_at"`
	UpdatedAt string             `bson:"updated_at" json:"updated_at"`
}

func Register(c *gin.Context) {
	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")

	// 添加基本验证
	if username == "" || email == "" || password == "" {
		response.BadRequestWithMessage(c, "Username, email and password are required")
		return
	}

	// 检查用户名是否已存在
	userCollection := databases.GetMongoCollection("user")
	// count, err := userCollection.CountDocuments(context.TODO(), bson.M{"username": username})
	// if err != nil {
	// 	response.InternalServerError(c)
	// 	return
	// }
	// if count > 0 {
	// 	response.BadRequestWithMessage(c, "Username already exists")
	// 	return
	// }

	// // 检查邮箱是否已存在
	// count, err = userCollection.CountDocuments(context.TODO(), bson.M{"email": email})
	// if err != nil {
	// 	response.InternalServerError(c)
	// 	return
	// }
	// if count > 0 {
	// 	response.BadRequestWithMessage(c, "Email already exists")
	// 	return
	// }

	passwordHash, err := crypto.HashPassword(password)
	if err != nil {
		response.InternalServerError(c)
		return
	}

	user := User2{
		Username:  username,
		Email:     email,
		Password:  passwordHash,
		Nickname:  username,
		Avatar:    "",
		CreatedAt: time.Now().Format(consts.DATETIME_NANO_FORMAT),
		UpdatedAt: time.Now().Format(consts.DATETIME_NANO_FORMAT),
	}
	log.Println(user)
	_, err = userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Printf("Error inserting user: %v", err)
		response.BadRequestWithMessage(c, "Failed to register")
		return
	}

	response.Success(c)
}

func Register2(c *gin.Context) {
	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")

	// 添加基本验证
	if username == "" || email == "" || password == "" {
		response.BadRequestWithMessage(c, "Username, email and password are required")
		return
	}

	// 检查用户名是否已存在
	userCollection := databases.GetMongoCollection("user")

	passwordHash, err := crypto.HashPassword(password)
	if err != nil {
		response.InternalServerError(c)
		return
	}

	user := bson.M{
		"username":   username,
		"email":      email,
		"password":   passwordHash,
		"nickname":   username,
		"avatar":     "",
		"created_at": time.Now().Format(consts.DATETIME_NANO_FORMAT),
		"updated_at": time.Now().Format(consts.DATETIME_NANO_FORMAT),
		"settings": bson.M{
			"language": "en",
		},
	}
	log.Println(user)
	_, err = userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Printf("Error inserting user: %v", err)
		response.BadRequestWithMessage(c, "Failed to register")
		return
	}

	response.Success(c)
}

func CheckIndexes(c *gin.Context) {
	collection := databases.GetMongoCollection("user")

	cursor, err := collection.Indexes().List(context.TODO())
	if err != nil {
		log.Printf("Error getting indexes: %v", err)
		response.InternalServerError(c)
		return
	}

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Printf("Error decoding indexes: %v", err)
		response.InternalServerError(c)
		return
	}

	// 打印所有索引
	for _, result := range results {
		log.Printf("Index: %+v", result)
	}

	response.SuccessWithData(c, results)
}
