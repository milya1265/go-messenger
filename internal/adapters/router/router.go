package router

import (
	"WSChats/internal/adapters/api/middleware"
	user "WSChats/internal/domain/user"
	"WSChats/internal/domain/ws"
	"github.com/gin-gonic/gin"
)

type Router struct {
	Router *gin.Engine
}

func NewRouter() *Router {
	return &Router{Router: gin.Default()}
}

func (r *Router) InitRoutes(auth auth.Middleware, userHandler user.Handler, wsHandler ws.Handler) {

	r.Router.Group("")

	{
		r.Router.POST("/register", userHandler.Register)
		r.Router.GET("/login", userHandler.Login, auth.Login)
		r.Router.GET("/ws/messenger", auth.Authorize, wsHandler.Messenger)
		r.Router.PUT("/user/:uuid", auth.Authorize, userHandler.UpdateUser)
		//authorize := r.Router.Group("")
		//authorize.Use(auth.Authorize)
		//{
		//
		//r.Router.PUT("/user/:uuid", userHandler.UpdateUser)
		//r.Router.DELETE("/user/:uuid")
		//}
	}
}

func (r *Router) Start(addr string) error {
	return r.Router.Run(addr)
}
