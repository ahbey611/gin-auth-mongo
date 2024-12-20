package middlewares

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS middleware configuration
func CORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"*", "http://localhost:3000"},                                                                           // set allowed origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},                                                     // set allowed HTTP methods
		AllowHeaders:     []string{"Authorization", "Origin", "X-Requested-With", "Content-Type", "Accept", "Access-Control-Allow-Origin"}, // set allowed headers
		ExposeHeaders:    []string{"Content-Length"},                                                                                       // set exposed headers
		AllowCredentials: true,                                                                                                             // allow credentials
		MaxAge:           12 * time.Hour,                                                                                                   // set cache time for preflight requests
	})
}
