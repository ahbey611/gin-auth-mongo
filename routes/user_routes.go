package routes

import (
	userController "gin-auth-mongo/controllers/user"
	"gin-auth-mongo/middlewares"

	"github.com/gin-gonic/gin"
)

// /api/v1/user/*
func UserRoutes(r *gin.RouterGroup) {
	user := r.Group("/user")
	user.Use(middlewares.JWTAuthMiddleware())
	{
		user.GET("/", userController.GetUser)
		user.PUT("/nickname", userController.UpdateNickname)
		user.PUT("/avatar", userController.UpdateAvatar)
		user.PUT("/avatar/upload", userController.UploadAvatar)
		user.POST("/avatar/status", userController.GetAvatarStatus)

		user.DELETE("/", userController.DeleteUserAccount)

		logout := user.Group("/logout")
		logout.POST("", userController.UserLogoutCurrentDevice)
		logout.POST("/all", userController.UserLogoutAllDevice)

	}
}
