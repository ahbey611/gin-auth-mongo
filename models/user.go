package models

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User model for table `user`
type User struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username         string             `bson:"username" json:"username"`
	Email            string             `bson:"email" json:"email"`
	Password         string             `bson:"password" json:"-"`
	Nickname         string             `bson:"nickname" json:"nickname"`
	Avatar           string             `bson:"avatar" json:"avatar"`
	CreatedAt        string             `bson:"created_at" json:"createdAt"`
	UpdatedAt        string             `bson:"updated_at" json:"updatedAt"`
	Country          string             `bson:"country" json:"country"`
	Settings         bson.M             `bson:"settings" json:"settings"`
	Premium          bool               `bson:"premium" json:"premium"`
	PremiumExpiredAt string             `bson:"premium_expired_at" json:"premiumExpiredAt"`
}
