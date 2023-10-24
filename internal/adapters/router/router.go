package router

import (
	"WSChats/internal/adapters/api/middleware"
	"WSChats/internal/domain/messenger"
	user "WSChats/internal/domain/user"
	"github.com/gin-gonic/gin"
)

type Router struct {
	Router *gin.Engine
}

func NewRouter() *Router {
	return &Router{Router: gin.Default()}
}

func (r *Router) InitRoutes(auth auth.Middleware, userHandler user.Handler, wsHandler messenger.Handler) {

	r.Router.Group("")

	{
		r.Router.POST("/register", userHandler.Register)
		r.Router.GET("/login", userHandler.Login, auth.Login)
		r.Router.GET("/ws/messenger", auth.Authorize, wsHandler.NewClient)
		r.Router.PUT("/user/:uuid", auth.Authorize, userHandler.UpdateUser)
		r.Router.POST("/chat", auth.Authorize, wsHandler.NewChat)
		//authorize := r.Router.Group("")
		//authorize.Use(auth.Authorize)
		//{
		//
		//r.Rout	er.PUT("/user/:uuid", userHandler.UpdateUser)
		//r.Router.DELETE("/user/:uuid")
		//}
	}
}

func (r *Router) Start(addr string) error {
	return r.Router.Run(addr)
}
