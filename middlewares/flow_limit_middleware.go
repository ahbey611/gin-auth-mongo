package middlewares

import (
	"net/http"

	"gin-auth-mongo/utils/consts"
	"gin-auth-mongo/utils/flow"
	"gin-auth-mongo/utils/response"

	"github.com/gin-gonic/gin"
)

func FlowLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		// check if IP is blocked
		isBlocked := flow.CheckIPBlocked(ip)
		if isBlocked {
			// IP is blocked due to too many requests
			response.Failure(c, http.StatusTooManyRequests, "Your IP is blocked due to too many requests.")
			c.Abort()
			return
		}

		// increase request count using atomic operation
		count, err := flow.IncreaseIPRequestCountV3(ip)
		if err != nil {
			response.InternalServerError(c)
			c.Abort()
			return
		}

		// if it's the first request, set the expiry time for the counter
		if count == 1 {
			flow.SetIPRequestCountExpiry(ip)
		}

		// check if it exceeds the limit
		if count > consts.FLOW_LIMIT_MAX {
			// block IP and set block time
			flow.SetIPBlocked(ip)
			response.Failure(c, http.StatusTooManyRequests, "Too many requests. Your IP is temporarily blocked.")
			c.Abort()
			return
		}

		c.Next()
	}
}
