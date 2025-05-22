package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/auth"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/auth/token"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
)

type TokenVerifier interface {
	VerifyByToken(tokenStr string) (*token.Payload, error)
}

func AuthMiddleware(authServ TokenVerifier, authZ auth.AuthZ, mandatory bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			if !mandatory {
				c.Next()
				return
			}
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
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		ctx := c.Request.Context()
		ctx = authZ.Authorize(ctx, *payload)
		c.Request = c.Request.WithContext(ctx)

		projLogger := ctx.Value(LoggerKey)
		projLogger.(MiddlewareLogger).Infow("Auth",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"ip", c.ClientIP(),
			"payload", payload,
			// "user-agent", c.Request.UserAgent(),
		)

		c.Next()
	}
}
