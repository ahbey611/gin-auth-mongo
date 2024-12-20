package databases

import (
	"context"
	"fmt"
	"log"

	// "gin-auth-mongo/databases/migrations"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var MongoDatabase *mongo.Database

func InitMongoDB() error {

	// 设置上下文超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoUser := os.Getenv("MONGO_USER")
	mongoPassword := os.Getenv("MONGO_PASSWORD")
	mongoUrl := os.Getenv("MONGO_URL")
	mongoDatabase := os.Getenv("MONGO_DATABASE")
	authSource := os.Getenv("MONGO_AUTH_SOURCE")

	if mongoUser == "" || mongoPassword == "" || mongoUrl == "" || mongoDatabase == "" || authSource == "" {
		panic("MongoDB environment variables are not set")
	}

	uri := fmt.Sprintf("mongodb://%s:%s@%s/%s?authSource=%s", mongoUser, mongoPassword, mongoUrl, mongoDatabase, authSource)
	// log.Println("uri: ", uri)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	// 检查连接
	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}

	log.Println("Connected to MongoDB")

	MongoClient = client
	MongoDatabase = client.Database(mongoDatabase)

	return nil
}

func GetMongoContext() context.Context {
	return context.TODO()
}

func GetMongoDatabase(databaseName string) *mongo.Database {
	return MongoClient.Database(databaseName)
}

func GetMongoCollection(collectionName string) *mongo.Collection {
	return MongoDatabase.Collection(collectionName)
}

//	func GetMongoCollection(databaseName string, collectionName string) *mongo.Collection {
//		return MongoClient.Database(databaseName).Collection(collectionName)
//	}

func RunMigrations() error {
	mongoUser := os.Getenv("MONGO_USER")
	mongoPassword := os.Getenv("MONGO_PASSWORD")
	mongoUrl := os.Getenv("MONGO_URL")
	mongoDatabase := os.Getenv("MONGO_DATABASE")
	authSource := os.Getenv("MONGO_AUTH_SOURCE")

	if mongoUser == "" || mongoPassword == "" || mongoUrl == "" || mongoDatabase == "" || authSource == "" {
		return fmt.Errorf("MongoDB environment variables are not set")
	}

	uri := fmt.Sprintf("mongodb://%s:%s@%s/%s?authSource=%s", mongoUser, mongoPassword, mongoUrl, mongoDatabase, authSource)

	// log.Println("uri: ", uri)

	m, err := migrate.New(
		"file://migrations",
		uri,
	)
	if err != nil {
		log.Printf("Migration initialization error: %v", err)
		return err
	}

	// get version
	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		log.Printf("Error getting migration version: %v", err)
		return err
	}

	// if dirty, force to the previous version
	if dirty {
		log.Printf("Database is dirty, forcing version %d", version-1)
		err = m.Force(int(version - 1))
		if err != nil {
			log.Printf("Error forcing version: %v", err)
			return err
		}
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Printf("Migration up error: %v", err)
		return err
	}

	/* err = migrations.AddPremiumExpiredAtFieldToUser(MongoDatabase)
	if err != nil {
		log.Printf("Migration error: %v", err)
		return err
	} */

	// Check indexes
	/* collection := GetMongoCollection("user")
	cursor, err := collection.Indexes().List(context.TODO())
	if err != nil {
		log.Printf("Error checking indexes: %v", err)
		return err
	}

	var indexes []bson.M
	if err = cursor.All(context.TODO(), &indexes); err != nil {
		log.Printf("Error decoding indexes: %v", err)
		return err
	} */

	// log.Printf("Current indexes: %+v", indexes)
	return nil
}
