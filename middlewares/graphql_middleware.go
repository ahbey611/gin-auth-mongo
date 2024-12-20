package middlewares

import (
	"context"
	"errors"
	"gin-auth-mongo/utils/consts"
	"gin-auth-mongo/utils/jwt"
	"gin-auth-mongo/utils/response"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Define a custom type for context keys
type contextKey string

// Define constants for context keys
const (
	TokenContextKey contextKey = "token"
)

// https://github.com/graphql-go/graphql/issues/378#issuecomment-568980284
// this middleware is used to get the token from the header and add it to the context
func GraphQLMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		requestPath := c.Request.URL.Path
		log.Println(requestPath)

		var token string
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader != "" {
			var err error
			token, err = jwt.GetTokenFromHeader(c)
			if err != nil {
				response.Unauthorized(c)
				return
			}
			c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "token", token))

			// get the claims
			claims, err := jwt.ParseJWTClaims(token)
			if err != nil {
				response.Unauthorized(c)
				return
			}

			exp := int64(claims["exp"].(float64))
			expiredAt := time.Unix(exp, 0).Format(consts.DATETIME_NANO_FORMAT)

			// check if the token is expired
			if time.Now().Unix() > exp {
				response.Unauthorized(c)
				return
			}

			jwtClaims := map[string]interface{}{
				"userID":        claims["sub"],
				"email":         claims["email"],
				"expiredAt":     expiredAt,
				"expiredAtUnix": exp,
			}
			c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "claims", jwtClaims))
		}

		// if the api is protected, the resolver SHOULD get the claims from the context to authorize the user
		c.Next()
	}
}

// if no claims found in the context, return unauthorized
func GetClaimsFromContext(ctx context.Context) (map[string]interface{}, error) {
	c := ctx.Value("claims")
	if c == nil {
		return nil, errors.New("Unauthorized")
	}
	claims, ok := c.(map[string]interface{})
	if !ok {
		return nil, errors.New("Unauthorized")
	}
	return claims, nil
}
