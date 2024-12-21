package middlewares

import (
	// "log"
	// "encoding/json"
	// "log"
	"time"

	// "time"
	// "time"
	// "log"
	// "log"
	"strings"

	"gin-auth-mongo/utils/consts"
	"gin-auth-mongo/utils/jwkmanager"
	"gin-auth-mongo/utils/jwt"
	"gin-auth-mongo/utils/response"

	"github.com/gin-gonic/gin"
	"github.com/square/go-jose/v3"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		requestPath := c.Request.URL.Path

		// the auth routes are not need to verify token
		if strings.HasPrefix(requestPath, "/api/v1/auth") {
			c.Next()
		}

		// else need to verify token
		token, err := jwt.GetTokenFromHeader(c)
		if err != nil {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		parsedJWT, err := jwt.ParseToken(token)
		if err != nil {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		// check if the token has a header
		if len(parsedJWT.Headers) <= 0 {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		// get public jwks
		publicJWKs, err := jwkmanager.GetPublicJWKs()
		if err != nil {
			response.InternalServerError(c)
			c.Abort()
			return
		}

		// create jwks
		JWKs := jose.JSONWebKeySet{
			Keys: publicJWKs,
		}
		publicJWK := JWKs.Key(parsedJWT.Headers[0].KeyID)

		// check if the public jwk is valid
		if len(publicJWK) <= 0 {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		// extract the claims
		allClaims := make(map[string]interface{})
		if err := parsedJWT.Claims(publicJWK[0].Key, &allClaims); err != nil {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		exp := int64(allClaims["exp"].(float64))

		// check if the token is expired
		if time.Now().Unix() > exp {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		c.Set("userID", allClaims["sub"])
		c.Set("email", allClaims["email"])
		c.Set("expiredAtUnix", exp)
		// convert to time
		expiredAt := time.Unix(exp, 0).Format(consts.DATETIME_NANO_FORMAT)
		c.Set("expiredAt", expiredAt)

		c.Next()
	}
}
