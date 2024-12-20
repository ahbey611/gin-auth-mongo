package routes

import (
	"gin-auth-mongo/middlewares"
	"os"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	// flow limit middleware
	r.Use(middlewares.CORSMiddleware(), middlewares.FlowLimitMiddleware())

	api := r.Group("/api")
	{
		logEnable := os.Getenv("LOG_ENABLE")
		if logEnable == "true" {
			api.Use(middlewares.DailyLogMiddleware(), middlewares.UnifiedLogMiddleware())
		}

		v1 := api.Group("/v1")
		{
			UserRoutes(v1)
			AuthRoutes(v1)
			FileRoutes(v1)
		}

	}

}
