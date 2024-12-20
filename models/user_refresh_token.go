package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// RefreshToken model for table `user_refresh_token`
type UserRefreshToken struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Token     string             `bson:"token" json:"token"`
	ExpiredAt string             `bson:"expired_at" json:"expired_at"`
	Device    string             `bson:"device" json:"device"`
}
