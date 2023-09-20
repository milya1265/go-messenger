package main

import (
	middleware1 "WSChats/internal/adapters/api/middleware"
	rout "WSChats/internal/adapters/router"
	"WSChats/internal/domain/auth"
	user2 "WSChats/internal/domain/user"
	"WSChats/internal/domain/ws"
	"WSChats/pkg/PostgreSQL"
	"WSChats/pkg/logger"
	"time"
)

var JWTkey = "qwerty"

func main() {
	logger := logger.GetLogger("debug")
	logger.Logger.Info("Logger init success")
	db := PostgreSQL.DataBase{}
	err := db.Open()
	if err != nil {
		logger.Logger.Error(err.Error())
		return
	}
	jwtManager, err := auth.NewManager(JWTkey, time.Hour*24, time.Minute*600)
	if err != nil {
		logger.Fatalf(err.Error())
	}

	wsRepo := ws.NewRepository(db.DB, &logger)
	wsService := ws.NewService(&wsRepo, &logger)
	messenger := ws.NewMessenger(&wsService, &logger)
	go messenger.Run()
	wsHandler := ws.NewHandler(&wsService, &logger, messenger)
	//wsHandler := ws.NewHandler(&wsService, &logger)

	authRepo := auth.NewRepository(db.DB, &logger)
	authService := auth.NewService(&authRepo, &logger, jwtManager)
	authHandler := middleware1.NewMiddleware(&authService, &logger)

	userRepo := user2.NewRepository(db.DB, &logger)
	userService := user2.NewService(&userRepo, &logger)
	userHandler := user2.NewHandler(&userService, &logger)

	router := rout.NewRouter()
	router.InitRoutes(authHandler, userHandler, wsHandler)
	err = router.Start("localhost:8080")
	if err != nil {
		return
	}

}
