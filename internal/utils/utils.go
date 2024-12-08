package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrTypecastJWT = errors.New("failed to typecast jwt claims")

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

func GenerateTokens(host string) (string, string, error) {
	var aToken, rToken string

	tokenLink, err := CreateLink()
	if err != nil {
		return aToken, rToken, err
	}

	// access token generation
	atTimeExp, err := strconv.Atoi(os.Getenv("ATEXPIRES"))
	if err != nil {
		return aToken, rToken, err
	}
	atExp := time.Now().Add(time.Duration(atTimeExp) * time.Second)
	claims := &jwt.MapClaims{
		"ExpiresAt":  atExp.Unix(),
		"Host":       host,
		"LinkString": tokenLink,
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	aToken, err = accessToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return aToken, rToken, err
	}

	// refresh token generation
	rtSH := sha256.Sum256([]byte(tokenLink))
	rToken = fmt.Sprintf("%x%v", rtSH, aToken[len(aToken)-6:])
	return aToken, rToken, nil
}

func CheckHost(aToken, host string) (bool, error) {
	jwtToken, err := jwt.Parse(aToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return false, err
	}
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return false, ErrTypecastJWT
	}
	if claims["Host"].(string) == host {
		return true, nil
	}
	return false, nil
}

func SendMasseg(mail string) error {
	return nil
}
