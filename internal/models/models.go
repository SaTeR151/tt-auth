package models

import (
	"time"
)

type ServerConfig struct {
	Port string
}

type DBConfig struct {
	User    string
	Pass    string
	Dbname  string
	Sslmode string
	Port    string
	Host    string
}

type Claims struct {
	Host string
	Exp  time.Time
}

type RTStruct struct {
	rt string
}
