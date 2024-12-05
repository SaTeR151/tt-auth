package handlers

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sater-151/tt-auth/internal/database"
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

func GetTokens(DB *database.DBStruct) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		logger.Info("start getting tokens")
		guid := req.FormValue("guid")
		if guid == "" {
			logger.Error("guid required")
			http.Error(res, "guid required", http.StatusBadRequest)
			return
		}

		tokenLink, err := service.CreateLink()
		if err != nil {
			logger.Error(err)
			http.Error(res, "", http.StatusInternalServerError)
			return
		}

		// access token generation
		atExp := time.Now().Add(30 * time.Minute)
		claims := &jwt.MapClaims{
			"Host":       req.Header.Get("Host"),
			"Expiration": atExp.Unix(),
			"LinkString": tokenLink,
		}

		accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
		atSign, err := accessToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
		if err != nil {
			logger.Error(err)
			http.Error(res, "", http.StatusInternalServerError)
			return
		}

		// refresh token generation
		rtExp := time.Now().Add(720 * time.Hour)
		rtSH := sha256.Sum256([]byte(tokenLink))
		rt := fmt.Sprintf("%x%v", rtSH, atSign[len(atSign)-6:])

		// save refresh token
		err = DB.InsertRT(guid, rt)
		if err != nil {
			if err == database.ErrUserNotFound {
				logger.Error(err)
				http.Error(res, "", http.StatusUnauthorized)
				return
			} else {
				logger.Error(err)
				http.Error(res, "", http.StatusInternalServerError)
				return
			}
		}

		rtB64 := base64.StdEncoding.EncodeToString([]byte(rt))

		http.SetCookie(res, &http.Cookie{
			Name:     "at",
			Value:    atSign,
			Expires:  atExp,
			HttpOnly: true,
		})
		http.SetCookie(res, &http.Cookie{
			Name:     "rt",
			Value:    string(rtB64),
			Expires:  rtExp,
			HttpOnly: true,
		})
		logger.Info("tokens have been sent")
	}

}

func RefreshTokens(s *service.ServiceStruct) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		logger.Info("start refresh token")
		guid := req.FormValue("guid")
		if guid == "" {
			logger.Error("guid required")
			http.Error(res, "guid required", http.StatusBadRequest)
			return
		}

		rtCookie, err := req.Cookie("rt")
		if err != nil {
			logger.Error(err)
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		rtb, err := base64.StdEncoding.DecodeString(rtCookie.Value)
		if err != nil {
			logger.Error(err)
			http.Error(res, "", http.StatusInternalServerError)
			return
		}

		comp, err := s.DB.CompareRT(string(rtb), guid)
		if err != nil {
			logger.Error(err)
			http.Error(res, "", http.StatusInternalServerError)
			return
		}
		if !comp {
			logger.Error(database.ErrUnauthorized)
			http.Error(res, "", http.StatusUnauthorized)
			return
		}

		tokenLink, err := service.CreateLink()
		if err != nil {
			logger.Error(err)
			http.Error(res, "", http.StatusInternalServerError)
			return
		}

		// access token generation
		atExp := time.Now().Add(30 * time.Minute)
		claims := &jwt.MapClaims{
			"Host":       req.Header.Get("Host"),
			"Expiration": atExp.Unix(),
			"LinkString": tokenLink,
		}

		accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
		atSign, err := accessToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
		if err != nil {
			logger.Error(err)
			http.Error(res, "", http.StatusInternalServerError)
			return
		}

		// refresh token generation
		rtExp := time.Now().Add(720 * time.Hour)
		rtSH := sha256.Sum256([]byte(tokenLink))
		rt := fmt.Sprintf("%x%v", rtSH, atSign[len(atSign)-6:])

		// save refresh token
		err = s.DB.InsertRT(guid, rt)
		if err != nil {
			if err == database.ErrUserNotFound {
				logger.Error(err)
				http.Error(res, "", http.StatusUnauthorized)
				return
			} else {
				logger.Error(err)
				http.Error(res, "", http.StatusInternalServerError)
				return
			}
		}

		rtB64 := base64.StdEncoding.EncodeToString([]byte(rt))

		http.SetCookie(res, &http.Cookie{
			Name:     "at",
			Value:    atSign,
			Expires:  atExp,
			HttpOnly: true,
		})
		http.SetCookie(res, &http.Cookie{
			Name:     "rt",
			Value:    string(rtB64),
			Expires:  rtExp,
			HttpOnly: true,
		})
		logger.Info("tokens have been refreshed")
	}
}
