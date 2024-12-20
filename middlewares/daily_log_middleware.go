package middlewares

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"gin-auth-mongo/utils/consts"

	"github.com/gin-gonic/gin"
)

var currentDate string
var logFile *os.File

// DailyLogMiddleware ensures a new log file is created each day
func DailyLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get current date
		today := time.Now().Format(consts.DATE_FORMAT)

		// if date changed, create a new log file
		if currentDate != today {
			// close old log file
			if logFile != nil {
				logFile.Close()
			}

			// create log directory
			logDir := "logs"

			// create new log file
			var err error
			logFile, err = os.OpenFile(filepath.Join(logDir, today+".log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
			if err != nil {
				log.Fatalf("Failed to open log file: %v", err)
			}

			// set log output to new file
			log.SetOutput(logFile)
			currentDate = today
		}

		// continue processing request
		c.Next()
	}
}
