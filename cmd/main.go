package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/sater-151/tt-auth/internal/config"
	"github.com/sater-151/tt-auth/internal/database"
	"github.com/sater-151/tt-auth/internal/handlers"
	logg "github.com/sater-151/tt-auth/internal/logger"
	"github.com/sater-151/tt-auth/internal/service"
	logger "github.com/sirupsen/logrus"
)

func main() {
	logg.Init()

	err := godotenv.Load()
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Info("getting configuration")
	serverConfig := config.GetServerConfig()
	dbConfig := config.GetDBConfig()

	logger.Info("database connecting")
	db, close, err := database.Open(dbConfig)
	if err != nil {
		logger.Error(err)
		return
	}
	defer close()
	logger.Info("database has been connected")

	logger.Info("starting migration")
	err = db.Migration()
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Info("migration done")

	service := service.New(db)

	r := chi.NewRouter()

	r.Get("/auth", handlers.GetTokens(service))
	r.Get("/refresh", handlers.RefreshTokens(service))

	logger.Info(fmt.Sprintf("server start at port: %s\n", serverConfig.Port))
	if err := http.ListenAndServe(":"+serverConfig.Port, r); err != nil {
		logger.Error(fmt.Sprintf("Server error: %s\n", err.Error()))
		return
	}
}
