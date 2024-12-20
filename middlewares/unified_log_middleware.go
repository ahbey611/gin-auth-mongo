package middlewares

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/url"
	"strings"
	"time"

	"gin-auth-mongo/utils/consts"

	"github.com/gin-gonic/gin"
)

// ResponseWriter is a custom response writer to capture response content
type ResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (rw *ResponseWriter) Write(b []byte) (int, error) {
	rw.body.Write(b) // capture response content
	return rw.ResponseWriter.Write(b)
}

// convert form data to JSON format
func formToJSON(form url.Values) string {
	jsonData := make(map[string]string)
	for key, values := range form {
		if len(values) > 0 {
			jsonData[key] = values[0]
		}
	}
	jsonBytes, _ := json.Marshal(jsonData)
	return string(jsonBytes)
}

// compress JSON string to single line format
func compressJSON(jsonStr string) string {
	var buf bytes.Buffer
	if err := json.Compact(&buf, []byte(jsonStr)); err != nil {
		return jsonStr // return original string if compression fails
	}
	return buf.String()
}

// UnifiedLogMiddleware records request and response logs
func UnifiedLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// create custom response writer
		rw := &ResponseWriter{ResponseWriter: c.Writer, body: bytes.NewBufferString("")}
		c.Writer = rw

		// get basic request information
		ip := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path

		// DO NOT LOG THE SENSITIVE PERSONAL DATA
		if strings.HasPrefix(path, "/api/auth/login") || strings.HasPrefix(path, "/api/auth/register") || strings.HasPrefix(path, "/api/auth/refresh") || strings.HasPrefix(path, "/api/auth/password-reset") {
			c.Next()
			return
		}

		var requestBody string
		contentType := c.Request.Header.Get("Content-Type")

		// handle different types of request bodies
		if contentType == "application/json" {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			requestBody = compressJSON(string(bodyBytes))
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // reassign to use it later
		} else if contentType == "application/x-www-form-urlencoded" {
			c.Request.ParseForm() // parse form data
			requestBody = formToJSON(c.Request.PostForm)
		} else if strings.HasPrefix(contentType, "multipart/form-data") {
			c.Request.ParseMultipartForm(consts.MAX_FILE_SIZE) // parse multipart/form-data, max 500 MB
			if c.Request.MultipartForm != nil {
				requestBody = formToJSON(c.Request.MultipartForm.Value)
			}
		}

		requestTime := time.Now()
		log.Printf("IP: %s | Method: %s | Path: %s | Request: %s\n", ip, method, path, requestBody)

		// continue processing request
		c.Next()

		// get response status code and captured response content
		statusCode := c.Writer.Status()
		responseBody := rw.body.String()

		duration := time.Since(requestTime)

		log.Printf("Path: %s | Status: %d | Duration: %v | Response: %s\n",
			path, statusCode, duration, responseBody)
	}
}
