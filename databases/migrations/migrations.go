package migrations

import (
	"context"
	// "gin-auth-mongo/models"
	// "gin-auth-mongo/utils/datetime"

	"go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddPremiumExpiredAtFieldToUser(db *mongo.Database) error {
	collection := db.Collection("user")
	// change the type of participants from array to map
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return err
	}

	defer cursor.Close(context.TODO())

	type OldUser struct {
		ID        string `bson:"_id,omitempty" json:"id"`
		Username  string `bson:"username" json:"username"`
		Email     string `bson:"email" json:"email"`
		Password  string `bson:"password" json:"-"`
		Nickname  string `bson:"nickname" json:"nickname"`
		Avatar    string `bson:"avatar" json:"avatar"`
		CreatedAt string `bson:"created_at" json:"createdAt"`
		UpdatedAt string `bson:"updated_at" json:"updatedAt"`
		Country   string `bson:"country" json:"country"`
		Settings  bson.M `bson:"settings" json:"settings"`
	}

	for cursor.Next(context.TODO()) {
		var user OldUser
		if err = cursor.Decode(&user); err != nil {
			return err
		}

		_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": user.ID}, bson.M{"$set": bson.M{
			"premium":            false,
			"premium_expired_at": "",
		}})
		if err != nil {
			return err
		}
	}

	return nil
}
