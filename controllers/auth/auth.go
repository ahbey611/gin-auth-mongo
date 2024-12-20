package auth

import (
	"gin-auth-mongo/models/requests"
	authService "gin-auth-mongo/services/auth"
	"gin-auth-mongo/utils/jwt"
	"gin-auth-mongo/utils/response"
	"gin-auth-mongo/utils/validation"

	"github.com/gin-gonic/gin"
)

// [POST] register
func UserEmailRegisterWithLink(c *gin.Context) {
	var request requests.EmailRegisterLinkRequest

	if err := validation.BindAndValidate(c, &request); err != nil {
		return
	}

	err := authService.UserEmailRegisterWithLink(&request)
	if err != nil {
		response.BadRequestWithMessage(c, err.Error())
		return
	}

	response.Success(c)
}

// [POST] register and send code to email
func UserEmailRegisterWithCode(c *gin.Context) {
	var request requests.EmailRegisterCodeRequest

	if err := validation.BindAndValidate(c, &request); err != nil {
		return
	}

	err := authService.UserEmailRegisterWithCode(&request)
	if err != nil {
		response.BadRequestWithMessage(c, err.Error())
		return
	}

	response.Success(c)
}

// [POST] verify email register
func UserEmailRegisterWithLinkVerify(c *gin.Context) {
	var request requests.EmailRegisterLinkVerifyRequest

	if err := validation.BindAndValidate(c, &request); err != nil {
		return
	}

	err := authService.UserEmailRegisterLinkVerify(&request)
	if err != nil {
		response.BadRequestWithMessage(c, err.Error())
		return
	}

	response.Success(c)
}

// [POST] verify registration code
func UserEmailRegisterWithCodeVerify(c *gin.Context) {
	var request requests.EmailRegisterCodeVerifyRequest

	if err := validation.BindAndValidate(c, &request); err != nil {
		return
	}

	err := authService.UserEmailRegisterCodeVerify(&request)
	if err != nil {
		response.BadRequestWithMessage(c, err.Error())
		return
	}

	response.Success(c)
}

// [GET] check registration link expired
func CheckUserEmailRegisterLinkExpired(c *gin.Context) {
	flowId := c.Query("flow_id")

	if flowId == "" {
		response.BadRequestWithMessage(c, "flow_id is required")
		return
	}

	info, err := authService.CheckUserEmailRegisterLinkExpired(flowId)
	if err != nil {
		response.BadRequestWithMessage(c, err.Error())
		return
	}

	response.SuccessWithData(c, info)

	// response.Success(c)
}

// [POST] login with email and password
func UserEmailLoginWithPassword(c *gin.Context) {
	var request requests.EmailLoginWithPasswordRequest

	if err := validation.BindAndValidate(c, &request); err != nil {
		return
	}

	user, token, err := authService.UserEmailLoginWithPassword(&request)
	if err != nil {
		response.BadRequestWithMessage(c, err.Error())
		return
	}

	response.SuccessWithData(c, map[string]interface{}{
		"user":  user,
		"token": token,
	})
}

// [POST] login with username and password
func UserUsernameLoginWithPassword(c *gin.Context) {
	var request requests.UsernameLoginWithPasswordRequest

	if err := validation.BindAndValidate(c, &request); err != nil {
		return
	}

	user, token, err := authService.UserUsernameLoginWithPassword(&request)
	if err != nil {
		response.BadRequestWithMessage(c, err.Error())
		return
	}

	response.SuccessWithData(c, map[string]interface{}{
		"user":  user,
		"token": token,
	})
}

// [POST] reset email password with link
func UserEmailResetPasswordWithLink(c *gin.Context) {
	var request requests.EmailPasswordResetLinkRequest

	if err := validation.BindAndValidate(c, &request); err != nil {
		return
	}

	err := authService.UserEmailResetPasswordWithLink(&request)
	if err != nil {
		response.BadRequestWithMessage(c, err.Error())
		return
	}

	response.Success(c)
}

// [POST] reset email password with code
func UserEmailResetPasswordWithCode(c *gin.Context) {
	var request requests.EmailPasswordResetCodeRequest

	if err := validation.BindAndValidate(c, &request); err != nil {
		return
	}

	err := authService.UserEmailResetPasswordWithCode(&request)
	if err != nil {
		response.BadRequestWithMessage(c, err.Error())
		return
	}

	response.Success(c)
}

// [POST] verify email password reset with link
func UserEmailResetPasswordWithLinkVerify(c *gin.Context) {
	var request requests.EmailPasswordResetLinkVerifyRequest

	if err := validation.BindAndValidate(c, &request); err != nil {
		return
	}

	err := authService.UserEmailResetPasswordLinkVerify(&request)
	if err != nil {
		response.BadRequestWithMessage(c, err.Error())
		return
	}

	response.Success(c)
}

// [POST] verify email password reset with code
func UserEmailResetPasswordWithCodeVerify(c *gin.Context) {
	var request requests.EmailPasswordResetCodeVerifyRequest

	if err := validation.BindAndValidate(c, &request); err != nil {
		return
	}

	err := authService.UserEmailResetPasswordCodeVerify(&request)
	if err != nil {
		response.BadRequestWithMessage(c, err.Error())
		return
	}

	response.Success(c)
}

// [GET] check reset password link expired
func CheckUserEmailResetPasswordLinkExpired(c *gin.Context) {
	flowId := c.Query("flow_id")

	if flowId == "" {
		response.BadRequestWithMessage(c, "flow_id is required")
		return
	}

	info, err := authService.CheckUserEmailResetPasswordLinkExpired(flowId)
	if err != nil {
		response.BadRequestWithMessage(c, err.Error())
		return
	}

	response.SuccessWithData(c, info)
}

// [GET] use refresh token to refresh access token
func RefreshToken(c *gin.Context) {
	token, err := jwt.GetTokenFromHeader(c)
	if err != nil {
		response.BadRequestWithMessage(c, err.Error())
		return
	}

	accessToken, err := authService.RefreshToken(token)
	if err != nil {
		response.BadRequestWithMessage(c, err.Error())
		return
	}
	response.SuccessWithData(c, accessToken)
}

// [GET] get token info
func GetTokenInfo(c *gin.Context) {
	token, err := jwt.GetTokenFromHeader(c)
	if err != nil {
		response.BadRequestWithMessage(c, err.Error())
		return
	}

	info, err := authService.GetTokenInfo(token)
	if err != nil {
		response.BadRequestWithMessage(c, err.Error())
		return
	}

	response.SuccessWithData(c, info)
}
