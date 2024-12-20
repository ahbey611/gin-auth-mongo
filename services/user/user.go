package user

import (
	"context"
	"errors"

	// "io"

	// "log"
	"gin-auth-mongo/databases"
	"gin-auth-mongo/models"
	"mime/multipart"

	// "gin-auth-mongo/models/requests"
	"gin-auth-mongo/repositories"
	"gin-auth-mongo/utils/consts"
	fileUtils "gin-auth-mongo/utils/file"

	// fileSe
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/minio/minio-go/v7"
)

func GetUserByID(userID string) (*models.User, error) {
	return repositories.GetUserByID(userID)
}

func UpdateNickname(userID string, nickname string) error {
	return repositories.UpdateNicknameByID(userID, nickname)
}

// Update avatar for rest api
func UpdateAvatarRest(userID string, avatar multipart.File, header *multipart.FileHeader) (string, error) {

	isImage := fileUtils.CheckFileIsImage(header.Header.Get("Content-Type"))
	if !isImage {
		return "", errors.New("invalid file type")
	}

	objectName := fileUtils.GenerateUniqueFileName("avatar", header.Filename)

	// set content type
	opts := minio.PutObjectOptions{ContentType: header.Header.Get("Content-Type")}

	// set timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	objectSize := header.Size
	if err := databases.PutObject(ctx, "", objectName, avatar, objectSize, opts); err != nil {
		return "", errors.New("failed to upload file")
	}

	// update avatar in database
	err := repositories.UpdateAvatarByID(userID, objectName)
	if err != nil {
		return "", errors.New("failed to update avatar in database")
	}

	return objectName, nil
}

func UpdateAvatarGraphQL(userID string, avatar *graphql.Upload) (string, error) {
	// Check if file is valid
	header := &multipart.FileHeader{
		Filename: avatar.Filename,
		Header:   make(map[string][]string),
	}
	header.Header.Set("Content-Type", avatar.ContentType)

	isImage := fileUtils.CheckFileIsImage(header.Header.Get("Content-Type"))
	if !isImage {
		return "", errors.New("invalid file type")
	}

	objectName := fileUtils.GenerateUniqueFileName("avatar", header.Filename)

	// set content type
	opts := minio.PutObjectOptions{ContentType: avatar.ContentType}

	// set timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// Upload file to MinIO
	if err := databases.PutObject(ctx, "", objectName, avatar.File, avatar.Size, opts); err != nil {
		return "", errors.New("failed to upload file")
	}

	// Update avatar in database
	err := repositories.UpdateAvatarByID(userID, objectName)
	if err != nil {
		return "", errors.New("failed to update avatar in database")
	}

	return objectName, nil
}

// Update avatar for rest or graphql
func UpdateAvatar(userID string, path string) error {

	exist := fileUtils.CheckFileIsExist(consts.MINIO_PUBLIC_BUCKET_NAME, path)
	if !exist {
		return errors.New("file not found")
	}

	// update avatar in database
	err := repositories.UpdateAvatarByID(userID, path)
	if err != nil {
		return errors.New("failed to update avatar in database")
	}

	return nil
}
