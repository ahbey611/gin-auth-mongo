package routes

import (
	authController "gin-auth-mongo/controllers/auth"

	"github.com/gin-gonic/gin"
)

// /api/v1/auth/*
func AuthRoutes(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		auth.POST("/register/email/link", authController.UserEmailRegisterWithLink)
		auth.POST("/register/email/link/verify", authController.UserEmailRegisterWithLinkVerify)
		auth.GET("/register/email/link/check", authController.CheckUserEmailRegisterLinkExpired)
		auth.POST("/register/email/code", authController.UserEmailRegisterWithCode)
		auth.POST("/register/email/code/verify", authController.UserEmailRegisterWithCodeVerify)

		auth.POST("/login/email", authController.UserEmailLoginWithPassword)
		auth.POST("/login/username", authController.UserUsernameLoginWithPassword)

		auth.POST("/password-reset/email/link", authController.UserEmailResetPasswordWithLink)
		auth.POST("/password-reset/email/link/verify", authController.UserEmailResetPasswordWithLinkVerify)
		auth.POST("/password-reset/email/code", authController.UserEmailResetPasswordWithCode)
		auth.POST("/password-reset/email/code/verify", authController.UserEmailResetPasswordWithCodeVerify)
		auth.GET("/password-reset/email/link/check", authController.CheckUserEmailResetPasswordLinkExpired)

		auth.GET("/token/info", authController.GetTokenInfo)
		auth.POST("/token/refresh", authController.RefreshToken)

	}
}
