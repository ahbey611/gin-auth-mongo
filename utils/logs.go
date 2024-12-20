package utils

import (
	"os"
)

func InitLoggerDir() {
	// create log directory
	logDir := "logs"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		os.Mkdir(logDir, os.ModePerm)
	}
}
