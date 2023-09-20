package ws

import (
	"WSChats/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type handler struct {
	service   Service
	logger    *logger.Logger
	messenger *Messenger
}

func NewHandler(s *Service, l *logger.Logger, m *Messenger) Handler { // , m *Messenger
	return &handler{
		service:   *s,
		logger:    l,
		messenger: m,
	}

}

type Handler interface {
	Messenger(c *gin.Context)
}

func (h *handler) Messenger(c *gin.Context) {
	h.logger.Logger.Info("Starting Messenger handler")

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

	h.logger.Info("Uuid new user: ", clientID)

	username, err := h.service.GetUsernameByUUID(clientID)
	if err != nil {
		h.logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		conn.Close()
		return
	}

	cl := &Client{
		Username: username,
		UUID:     clientID,
		Conn:     *conn,
		Messages: make(chan *Message, 10),
	}

	h.messenger.Online <- cl
	go cl.SendWS()
	go cl.ListenWS(h.messenger)
}
