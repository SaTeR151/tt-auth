package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/jmoiron/sqlx"
	"github.com/sater-151/tt-auth/internal/models"
	logger "github.com/sirupsen/logrus"
)

var ErrUserNotFound = errors.New("user not found")
var ErrUnauthorized = errors.New("unauthorized user")

type DBInterface interface {
	Migration() error
	UpdateRT(guid, rt string) error
	SelectMail(guid string) (string, error)
	GetBcrypt(rToken string) (string, error)
	GetToken(guid string) (string, error)
}

type DBStruct struct {
	db *sql.DB
}

func Open(config models.DBConfig) (*DBStruct, func() error, error) {
	var err error
	connInfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Host,
		config.User,
		config.Pass,
		config.Dbname,
		config.Port,
		config.Sslmode)
	db, err := sql.Open("pgx", connInfo)
	if err != nil {
		return nil, nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, nil, err
	}
	DB := &DBStruct{db: db}

	return DB, db.Close, nil
}

func (db *DBStruct) Migration() error {
	driver, err := postgres.WithInstance(db.db, &postgres.Config{})
	if err != nil {
		return err
	}
	migrator, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		return err
	}
	err = migrator.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}

func (db *DBStruct) UpdateRT(guid, rt string) error {
	logger.Debug("set refresh token")
	res, err := db.db.Exec("UPDATE users_auth SET rt=crypt($1, 'nothing') WHERE user_id=$2 RETURNING rt", rt, guid)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	logger.Debug("refresh token has beeb set")
	return nil
}

func (db *DBStruct) SelectMail(guid string) (string, error) {
	return "example@mail.ru", nil
}

func (db *DBStruct) GetBcrypt(rToken string) (string, error) {
	var rTokenBcrypt sql.NullString
	err := db.db.QueryRow("SELECT crypt($1, 'nothing')", rToken).Scan(&rTokenBcrypt)
	if err != nil {
		return "", err
	}
	return rTokenBcrypt.String, nil
}

func (db *DBStruct) GetToken(guid string) (string, error) {
	var rtDB sql.NullString
	err := db.db.QueryRow("SELECT rt FROM users_auth WHERE user_id=$1", guid).Scan(&rtDB)
	if err != nil {
		return "", err
	}
	if !rtDB.Valid {
		return "", ErrUserNotFound
	}
	return rtDB.String, nil
}
