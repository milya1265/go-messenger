package user

import (
	"WSChats/internal/adapters/api/DTO"

	"WSChats/pkg/logger"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type handler struct {
	service Service
	logger  logger.Logger
}

func NewHandler(s *Service, l *logger.Logger) Handler {

	return &handler{
		service: *s,
		logger:  *l,
	}

}

type Handler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
}

func (h *handler) Register(c *gin.Context) {
	h.logger.Info("Start CreateUser handler")

	var reqU DTO.CreateUserReq
	if err := c.BindJSON(&reqU); err != nil {
		h.logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	resU := &DTO.CreateUserRes{}

	resU, err := h.service.CreateUser(c.Request.Context(), &reqU)
	if err != nil {
		h.logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("CreateUser successful for user, ", resU.Username)
	c.JSON(http.StatusCreated, gin.H{"message": "user created", "user": resU})
}

func (h *handler) Login(c *gin.Context) {
	h.logger.Info("Start Login handler")

	var reqU DTO.GetUserByEmailReq
	if err := c.BindJSON(&reqU); err != nil {
		h.logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	resU, err := h.service.Login(c.Request.Context(), &reqU)
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			h.logger.Error(err.Error())
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Incorrect username or password"})
			return
		}

		h.logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.Set("uuid", resU.UUID)
	c.Next()
	//c.SetSameSite(http.SameSiteLaxMode)
	//c.SetCookie("Authorization", resU.JWTtoken, 3600*24*30, "", "", false, true)

	//c.JSON(http.StatusOK, gin.H{
	//	"user": *resU,
	//})
}

func (h *handler) UpdateUser(c *gin.Context) {
}

func (h *handler) DeleteUser(c *gin.Context) {

}
