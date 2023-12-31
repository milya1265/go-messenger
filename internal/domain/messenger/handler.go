package messenger

import (
	"WSChats/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type handler struct {
	service   Service
	logger    *logger.Logger
	messenger *Manager
}

func NewHandler(s *Service, l *logger.Logger, m *Manager) Handler { // , m *NewClient
	return &handler{
		service:   *s,
		logger:    l,
		messenger: m,
	}

}

type Handler interface {
	NewClient(c *gin.Context)
	//NewChat(c *gin.Context)
}

func (h *handler) NewClient(c *gin.Context) {
	h.logger.Logger.Info("Starting NewClient handler")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		conn.Close()
		return
	}

	clientID, ok := c.Keys["uuid"].(string)

	if !ok {
		h.logger.Error("uuid is not found")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "uuid is not found"})
		conn.Close()
		return
	}

	err = conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		return
	}
	conn.SetPongHandler(func(string) error {
		err := conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			return err
		}
		return nil
	})
	go ping(conn)

	h.logger.Info("Uuid new user: ", clientID)

	username, err := h.service.GetUsernameByUUID(clientID)
	if err != nil {
		h.logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		conn.Close()
		return
	}

	cl := newClient(username, clientID, conn, *h.logger, &h.service)

	h.messenger.Online <- cl
	go cl.SendRes()
	go cl.HandleReq(h.messenger)

}

/*func (h *handler) NewChat(c *gin.Context) {
	h.logger.Info("Starting new chat handler")

	var chatReq NewChatReq

	if err := c.BindJSON(&chatReq); err != nil {
		h.logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	chatReq.Creator = c.Keys["uuid"].(string)

	chatRes, err := h.service.NewChat(&chatReq)
	if err != nil {
		if errors.Is(err, errors.New(ErrorNotTwoMembersInDirect)) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": ErrorNotTwoMembersInDirect})
		} else if errors.Is(errors.New("chat has already created"), err) {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": ErrorNotTwoMembersInDirect})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	}

	c.JSON(http.StatusOK, chatRes)
}*/
