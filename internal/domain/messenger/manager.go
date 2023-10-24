package messenger

import (
	"WSChats/pkg/logger"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

var (

	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Maximum message size allowed from peer
	maxMessageSize int64 = 512

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	defaultBroadcastQueueSize = 10000
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
		// return r.Header.Get("Origin") != "http://"+r.Host
	},
}

type Manager struct {
	Errors    chan *error
	Broadcast chan *Message
	Quit      chan struct{}
	Service   Service
	Online    chan *Client
	Offline   chan *Client
	Clients   sync.Map
	logger    *logger.Logger
}

func NewManager(s *Service, l *logger.Logger) *Manager {
	return &Manager{
		Broadcast: make(chan *Message),
		Quit:      make(chan struct{}),
		Service:   *s,
		Online:    make(chan *Client),
		Offline:   make(chan *Client),
		Clients:   sync.Map{},
		logger:    l,
	}
}

func (m *Manager) Run() {
	for {
		select {

		case cl := <-m.Offline:
			m.Clients.Delete(cl.UUID)
			m.logger.Info("Delete connection with user ", cl.Username)

		case cl := <-m.Online:
			m.Clients.Store(cl.UUID, cl)
			m.logger.Info("Create connection with user ", cl.Username)

		case message := <-m.Broadcast:
			message, err := m.Service.NewMessage(message)
			if err != nil {
				m.logger.Error(err.Error())
			}

			receivers, err := m.Service.GetChatMembers(message.ChatID)
			if err != nil {
				sen, ok := m.Clients.Load(message.Sender)
				if ok {
					sen.(*Client).Errors <- &err
					continue
				}
			}

			for _, rec := range receivers {
				cl, ok := m.Clients.Load(rec)
				if ok {
					cl.(*Client).Messages <- message
				}
			}
		}
	}
}

func ping(ws *websocket.Conn) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		ws.Close()
	}()
	for {
		select {
		case <-ticker.C:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				// Don't use return, it will not trigger the defer function.
				break
			}
		}
	}
}
