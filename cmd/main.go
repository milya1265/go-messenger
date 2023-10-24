package main

import (
	middleware1 "WSChats/internal/adapters/api/middleware"
	rout "WSChats/internal/adapters/router"
	"WSChats/internal/domain/auth"
	"WSChats/internal/domain/messenger"
	user2 "WSChats/internal/domain/user"
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
	jwtManager, err := auth.NewManager(JWTkey, time.Hour*10000, time.Hour*10000)
	if err != nil {
		logger.Fatalf(err.Error())
	}

	wsRepo := messenger.NewRepository(db.DB, &logger)
	wsService := messenger.NewService(&wsRepo, &logger)
	manager := messenger.NewManager(&wsService, &logger)
	go manager.Run()
	wsHandler := messenger.NewHandler(&wsService, &logger, manager)
	//wsHandler := messenger.NewHandler(&wsService, &logger)

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
