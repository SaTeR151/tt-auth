package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/sater-151/tt-auth/internal/config"
	"github.com/sater-151/tt-auth/internal/database"
	"github.com/sater-151/tt-auth/internal/handlers"
	"github.com/sater-151/tt-auth/internal/service"
	logger "github.com/sirupsen/logrus"
)

// для настройки logrus
func init() {
	logger.SetFormatter(&logger.TextFormatter{FullTimestamp: true})
	lvl, ok := os.LookupEnv("LOG_LEVEL")

	if !ok {
		lvl = "debug"
	}

	ll, err := logger.ParseLevel(lvl)
	if err != nil {
		ll = logger.DebugLevel
	}

	logger.SetLevel(ll)
}

func main() {
	logger.Info(".env download")
	err := godotenv.Load()
	if err != nil {
		logger.Error(err)
		return
	}

	logger.Info("configuration generation")
	serverConfig := config.GetServerConfig()
	dbConfig := config.GetDBConfig()

	logger.Info("database connection")
	db, err := database.Open(dbConfig)
	if err != nil {
		logger.Error(fmt.Sprintf("Database connect error: %s\n", err.Error()))
		return
	}
	defer db.Close()
	logger.Info("database connected")

	logger.Info("start migration")
	err = db.Migration()
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Info("migration done")

	service := service.New(db)

	r := chi.NewRouter()

	r.Get("/secret", handlers.GetTokens(db))
	r.Get("/refresh", handlers.RefreshTokens(service))

	logger.Info(fmt.Sprintf("server start at port: %s\n", serverConfig.Port))
	if err := http.ListenAndServe(":"+serverConfig.Port, r); err != nil {
		logger.Error(fmt.Sprintf("Server startup error: %s\n", err.Error()))
		return
	}
}
