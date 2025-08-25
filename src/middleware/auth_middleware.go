package middleware

import (
	"github.com/gin-gonic/gin"
	"st-novel-go/src/utils"
	"strings"
)

const (
	AuthorizationHeader = "Authorization"
	BearerPrefix        = "Bearer "
	UserClaimsKey       = "userClaims"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			utils.FailWithUnauthorized(c, "Authorization header is required")
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, BearerPrefix) {
			utils.FailWithUnauthorized(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, BearerPrefix)
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			utils.FailWithUnauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		c.Set(UserClaimsKey, claims)
		c.Next()
	}
}
