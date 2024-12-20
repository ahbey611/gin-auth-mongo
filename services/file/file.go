package file

import (
	"context"
	"errors"
	"gin-auth-mongo/databases"
	"gin-auth-mongo/utils/consts"
	fileUtils "gin-auth-mongo/utils/file"
	"mime/multipart"

	// "os"
	"time"

	"github.com/minio/minio-go/v7"
)

func UploadImage(bucketName string, folder string, file multipart.File, header *multipart.FileHeader, userID string) (string, string, error) {

	// set default bucket name and path
	if bucketName == "" || folder == "" {
		bucketName = consts.MINIO_PUBLIC_BUCKET_NAME
	}

	// check if the file is an image
	isImage := fileUtils.CheckFileIsImage(header.Header.Get("Content-Type"))
	if !isImage {
		return "", "", errors.New("invalid image type")
	}

	objectSize := header.Size

	// check if the file size is too large
	if !fileUtils.CheckImageSize(objectSize) {
		return "", "", errors.New("image size is too large")
	}

	rawName, _ := fileUtils.GetFileNameAndExtension(header.Filename)

	// generate unique file name
	objectName, err := fileUtils.GenerateEncryptedFileName(header.Filename)
	if err != nil {
		return "", "", err
	}

	// upload file to minio
	opts := minio.PutObjectOptions{
		ContentType: header.Header.Get("Content-Type"),
		UserTags: map[string]string{
			"userID": userID,
			"name":   rawName,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	path := folder + "/" + objectName

	if err := databases.PutObject(ctx, bucketName, path, file, objectSize, opts); err != nil {
		return "", "", err
	}

	return bucketName, path, nil
}
