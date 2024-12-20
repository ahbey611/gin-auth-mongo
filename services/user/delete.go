package user

import (
	"context"
	"gin-auth-mongo/databases"
	"gin-auth-mongo/repositories"

	"go.mongodb.org/mongo-driver/mongo"
)

func DeleteUserAccount(userID string) error {
	session, err := databases.MongoClient.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(context.Background())

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		// delete all the refresh tokens from the database
		err = repositories.DeleteRefreshTokenByUserID(userID)
		if err != nil {
			return nil, err
		}

		// delete the user from the database
		err = repositories.DeleteUserByID(userID)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}

	_, err = session.WithTransaction(context.Background(), callback)
	if err != nil {
		return err
	}
	return nil
}
