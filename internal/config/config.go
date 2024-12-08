package config

import (
	"os"

	"github.com/sater-151/tt-auth/internal/models"
	logger "github.com/sirupsen/logrus"
)

func GetServerConfig() models.ServerConfig {
	var config models.ServerConfig
	config.Port = os.Getenv("SERVER_PORT")
	return config
}

func GetDBConfig() models.DBConfig {
	var dbConfig models.DBConfig
	dbConfig.User = os.Getenv("POSTGRES_USER")
	if dbConfig.User == "" {
		logger.Warn("postgres user is empty")
	}
	dbConfig.Pass = os.Getenv("POSTGRES_PASSWORD")
	if dbConfig.Pass == "" {
		logger.Warn("postgres password if empty")
	}
	dbConfig.Dbname = os.Getenv("POSTGRES_DB")
	if dbConfig.Dbname == "" {
		logger.Warn("postgres database name if empty")
	}
	dbConfig.Sslmode = os.Getenv("SSLMODE")
	if dbConfig.Sslmode == "" {
		logger.Warn("postgres sslmode is empty")
	}
	dbConfig.Port = os.Getenv("DB_PORT")
	if dbConfig.Port == "" {
		logger.Warn("postgres port is empty")
	}
	dbConfig.Host = os.Getenv("DB_HOST")
	if dbConfig.Host == "" {
		logger.Warn("postgres host is empty")
	}
	return dbConfig
}
