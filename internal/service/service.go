package service

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sater-151/tt-auth/internal/database"
)

type ServiceStruct struct {
	DB *database.DBStruct
}

func New(db *database.DBStruct) *ServiceStruct {
	service := &ServiceStruct{DB: db}
	return service
}

const str = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"

func CreateLink() (string, error) {
	tokenLink := make([]byte, 32)
	_, err := rand.Read(tokenLink)
	if err != nil {
		return "", err
	}
	for i, j := range tokenLink {
		tokenLink[i] = str[j%byte(len(str))]
	}
	return string(tokenLink), nil
}

func (s *ServiceStruct) EmailWarning(guid string) error {
	mail, err := s.DB.SelectMail(guid)
	if err != nil {
		return err
	}
	err = SendMasseg(mail)
	if err != nil {
		return err
	}
	return nil
}

func SendMasseg(mail string) error {
	return nil
}

func (s *ServiceStruct) GenerateTokens(host string) (string, string, error) {
	var aToken, rToken string

	tokenLink, err := CreateLink()
	if err != nil {
		return aToken, rToken, err
	}

	// access token generation
	atExp := time.Now().Add(30 * time.Minute)
	claims := &jwt.MapClaims{
		"ExpiresAt":  atExp.Unix(),
		"Host":       host,
		"LinkString": tokenLink,
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	atSign, err := accessToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return aToken, rToken, err
	}

	// refresh token generation
	rtSH := sha256.Sum256([]byte(tokenLink))
	rToken = fmt.Sprintf("%x%v", rtSH, atSign[len(atSign)-6:])

	return aToken, rToken, nil
}
