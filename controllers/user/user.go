package controllers

import (
	// "log"
	// "log"
	"context"
	"log"

	// "log"
	"gin-auth-mongo/databases"
	"gin-auth-mongo/models/requests"
	fileServices "gin-auth-mongo/services/file"
	userServices "gin-auth-mongo/services/user"
	"gin-auth-mongo/utils/consts"

	// fileUtils "gin-auth-mongo/utils/file"
	"gin-auth-mongo/utils/response"
	"gin-auth-mongo/utils/validation"

	// "strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

// [GET] get user
func GetUser(c *gin.Context) {
	userID := c.GetString("userID")
	user, err := userServices.GetUserByID(userID)
	if err != nil {
		response.InternalServerError(c)
		return
	}

	response.SuccessWithData(c, user)
}

// [PUT] update nickname
func UpdateNickname(c *gin.Context) {

	// example: get claims from context
	/* userID := c.GetString("userID")
	email := c.GetString("email")
	expiredAt := c.GetString("expiredAt")
	expiredAtUnix := c.GetInt64("expiredAtUnix") */

	var request requests.UpdateNicknameRequest
	if err := validation.BindAndValidate(c, &request); err != nil {
		return
	}

	userID := c.GetString("userID")
	// log.Println("userID: ", userID)

	err := userServices.UpdateNickname(userID, request.Nickname)
	if err != nil {
		response.InternalServerError(c)
		return
	}

	response.Success(c)
}

// [PUT] update avatar
func UpdateAvatar(c *gin.Context) {

	var request requests.UpdateAvatarRequest
	if err := validation.BindAndValidate(c, &request); err != nil {
		return
	}
	userID := c.GetString("userID")

	err := userServices.UpdateAvatar(userID, request.Avatar)
	if err != nil {
		response.BadRequestWithMessage(c, err.Error())
		return
	}

	response.Success(c)
}

// [POST] upload avatar
func UploadAvatar(c *gin.Context) {

	file, header, err := c.Request.FormFile("avatar")
	if err != nil {
		response.BadRequestWithMessage(c, "Failed to get file")
		return
	}
	defer file.Close()

	userID := c.GetString("userID")
	bucketName, path, err := fileServices.UploadImage(consts.MINIO_PUBLIC_BUCKET_NAME, "avatars", file, header, userID)

	if err != nil {
		response.InternalServerError(c)
		return
	}
	response.SuccessWithData(c, gin.H{"bucket": bucketName, "path": path})
}

// [POST] get avatar status
func GetAvatarStatus(c *gin.Context) {

	fileName := c.PostForm("fileName")

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// check if file is exist in minio
	objectInfo, err := databases.StatObject(ctx, consts.MINIO_PUBLIC_BUCKET_NAME, fileName, minio.StatObjectOptions{})

	if err != nil {
		response.BadRequestWithMessage(c, "File not found")
		return
	}
	// log.Println("objectInfo: ", objectInfo)

	tag, err := databases.GetObjectTag(ctx, consts.MINIO_PUBLIC_BUCKET_NAME, fileName)

	if err != nil {
		response.BadRequestWithMessage(c, "File not found")
		return
	}

	response.SuccessWithData(c, gin.H{"objectInfo": objectInfo, "tag": tag})
}

// [POST] logout current device
func UserLogoutCurrentDevice(c *gin.Context) {
	var request requests.LogoutRequest

	if err := validation.BindAndValidate(c, &request); err != nil {
		return
	}

	userID := c.GetString("userID")

	err := userServices.UserLogoutCurrentDevice(userID, &request)
	if err != nil {
		response.BadRequestWithMessage(c, err.Error())
		return
	}

	response.Success(c)
}

// [POST] logout all device
func UserLogoutAllDevice(c *gin.Context) {
	userID := c.GetString("userID")

	err := userServices.UserLogoutAllDevice(userID)
	if err != nil {
		response.BadRequestWithMessage(c, err.Error())
		return
	}

	response.Success(c)
}

// [DELETE] delete user account
func DeleteUserAccount(c *gin.Context) {
	userID := c.GetString("userID")

	log.Println("userID: ", userID)

	// response.SuccessWithData(c, gin.H{"userID": userID})
	// return

	err := userServices.DeleteUserAccount(userID)
	if err != nil {
		response.InternalServerError(c)
		return
	}
	response.Success(c)
}
