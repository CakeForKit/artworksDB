package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/auth/token"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

type TokenVerifier interface {
	VerifyByToken(tokenStr string) (*token.Payload, error)
}

func AuthMiddleware(authServ TokenVerifier) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		accessToken := fields[1]
		payload, err := authServ.VerifyByToken(accessToken)
		fmt.Printf("payload %+v\n", payload)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.Set(authorizationPayloadKey, payload)
		c.Next()
	}
}

// func AuthUserMiddleware(authServ auth.AuthUser) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		authorizationHeader := c.GetHeader(authorizationHeaderKey)
// 		if len(authorizationHeader) == 0 {
// 			err := errors.New("authorization header is not provided")
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
// 			return
// 		}

// 		fields := strings.Fields(authorizationHeader)
// 		if len(fields) < 2 {
// 			err := errors.New("invalid authorization header format")
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
// 			return
// 		}

// 		authorizationType := strings.ToLower(fields[0])
// 		if authorizationType != authorizationTypeBearer {
// 			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
// 			return
// 		}

// 		accessToken := fields[1]
// 		payload, err := authServ.VerifyByToken(accessToken)
// 		if err != nil {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
// 			return
// 		}
// 		c.Set(authorizationPayloadKey, payload)
// 		c.Next()
// 	}
// }

// func AuthEmployeeMiddleware(authServ auth.AuthEmployee) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		authorizationHeader := c.GetHeader(authorizationHeaderKey)
// 		if len(authorizationHeader) == 0 {
// 			err := errors.New("authorization header is not provided")
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
// 			return
// 		}

// 		fields := strings.Fields(authorizationHeader)
// 		if len(fields) < 2 {
// 			err := errors.New("invalid authorization header format")
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
// 			return
// 		}

// 		authorizationType := strings.ToLower(fields[0])
// 		if authorizationType != authorizationTypeBearer {
// 			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
// 			return
// 		}

// 		accessToken := fields[1]
// 		payload, err := authServ.VerifyByToken(accessToken)
// 		if err != nil {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
// 			return
// 		}
// 		c.Set(authorizationPayloadKey, payload)
// 		c.Next()
// 	}
// }
