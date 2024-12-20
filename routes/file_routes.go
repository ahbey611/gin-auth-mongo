package routes

import (
	fileController "gin-auth-mongo/controllers/file"

	"github.com/gin-gonic/gin"
)

func FileRoutes(r *gin.RouterGroup) {
	file := r.Group("/file")
	{
		file.POST("/upload", fileController.UploadFile)
		file.POST("/decrypt", fileController.DecryptFileName)
	}
}
