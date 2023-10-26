package auth

import (
	"WSChats/internal/domain/auth"
	"WSChats/pkg/logger"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type middleware struct {
	service auth.Service
	logger  *logger.Logger
}

func NewMiddleware(s *auth.Service, l *logger.Logger) Middleware {
	return &middleware{
		logger:  l,
		service: *s,
	}

}

type Middleware interface {
	Login(c *gin.Context)
	Authorize(c *gin.Context)
}

/*func (m *middleware) Authorize(c *gin.Context) {
	m.logger.Info("Start Authorize handler")

	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		m.logger.Error(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	token, err := jwt.ParseSubject(tokenString, func(token *jwt.Token) (interface{}, error) {
		return user2.JWTkey, nil
	})
	if err != nil {
		m.logger.Error(err.Error())
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			m.logger.Error("time is out")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		user := claims["sub"].(DTO.GetUserByEmailRes)

		c.Set("user", user)
		c.Next()
	} else {
		log.Println("Validation error")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}
*/

func (m *middleware) Login(c *gin.Context) {
	uuid, ok := c.Keys["uuid"].(string)
	if ok == false {
		return
	}

	access, err := m.service.NewSession(c.Request.Context(), uuid)
	if err != nil {
		m.logger.Error(err.Error())
		c.JSON(http.StatusForbidden, gin.H{"error": "try again later"})
		return
	}
	c.Header("jwt", access)
	m.logger.Info("Login successful for user, ", uuid)
}

func (m *middleware) Authorize(c *gin.Context) {
	m.logger.Info("Starting Authorize method")
	access := c.GetHeader("jwt")

	uuid, access, err := m.service.Authorize(c.Request.Context(), access)
	if err != nil {
		if errors.Is(err, auth.ErrorTokenTimeOut) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"messege": "login again"})
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}

	c.Header("jwt", access)
	c.Set("uuid", uuid)
	c.Next()
}
