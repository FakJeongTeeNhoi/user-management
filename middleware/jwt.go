package middleware

import (
	"fmt"
	"github.com/FakJeongTeeNhoi/user-management/model/response"
	"github.com/FakJeongTeeNhoi/user-management/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"strings"
	"time"
)

const (
	headerBearerPrefix  = "Bearer "
	headerAuthorization = "authorization"
)

var (
	whiteList = []string{
		"/health",
		"/api/auth/login",
	}
)

func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		if isWhiteList(c.FullPath()) {
			c.Next()
			return
		}

		token := extractToken(c.GetHeader(headerAuthorization))
		if token == "" {
			response.Unauthorized("Missing token").AbortWithError(c)
			return
		}

		if err := verifyToken(token); err != nil {
			response.Unauthorized("Invalid token").AbortWithError(c)
			return
		}

		c.Next()
	}
}

func SetAccountInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c.GetHeader(headerAuthorization))
		accountInfo, err := service.GetInfoFromToken(token)
		if err != nil {
			response.Unauthorized("Invalid token").AbortWithError(c)
			return
		}
		c.Set("user", accountInfo)
		c.Request.Header.Del(headerAuthorization)

		c.Next()
	}
}

func ValidatePermission(allowList []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		info, exist := c.Get("user")
		if !exist {
			response.Unauthorized("Unauthorized").AbortWithError(c)
			return
		}
		userType := info.(jwt.MapClaims)["type"].(string)
		if !contains(allowList, userType) {
			response.Unauthorized("Unauthorized").AbortWithError(c)
			return
		}
		c.Next()
	}

}

func contains(list []string, target string) bool {
	for _, item := range list {
		if item == target {
			return true
		}
	}
	return false
}

func isWhiteList(path string) bool {
	for _, p := range whiteList {
		if p == path {
			return true
		}
	}
	return false
}

func extractToken(authorization string) string {
	if !strings.HasPrefix(authorization, headerBearerPrefix) {
		return ""
	}
	return strings.TrimPrefix(authorization, headerBearerPrefix)
}

func verifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	info, err := service.GetInfoFromToken(tokenString)
	if err != nil {
		return err
	}

	if info["exp"].(float64) < float64(time.Now().Unix()) {
		return fmt.Errorf("token is expired")
	}

	return nil
}
