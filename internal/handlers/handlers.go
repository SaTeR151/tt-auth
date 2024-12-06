package handlers

import (
	"database/sql"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/sater-151/tt-auth/internal/database"
	"github.com/sater-151/tt-auth/internal/service"
	logger "github.com/sirupsen/logrus"
)

func GetTokens(s *service.ServiceStruct) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		logger.Info("getting tokens")
		guid := req.FormValue("guid")
		if guid == "" {
			logger.Error("guid required")
			http.Error(res, "guid required", http.StatusBadRequest)
			return
		}

		aToken, rToken, err := s.GenerateTokens(req.Host)
		if err != nil {
			logger.Error(err)
			http.Error(res, "", http.StatusInternalServerError)
			return
		}
		// save refresh token
		err = s.DB.InsertRT(guid, rToken)
		if err != nil {
			if err == sql.ErrNoRows {
				logger.Error(err)
				http.Error(res, "", http.StatusUnauthorized)
				return
			} else {
				logger.Error(err)
				http.Error(res, "", http.StatusInternalServerError)
				return
			}
		}

		rtB64 := base64.StdEncoding.EncodeToString([]byte(rToken))

		atExp := time.Now().Add(30 * time.Second)
		rtExp := time.Now().Add(720 * time.Hour)
		http.SetCookie(res, &http.Cookie{
			Name:     "at",
			Value:    aToken,
			Expires:  atExp,
			HttpOnly: true,
		})
		http.SetCookie(res, &http.Cookie{
			Name:     "rt",
			Value:    rtB64,
			Expires:  rtExp,
			HttpOnly: true,
		})
		logger.Info("tokens have been sent")
	}

}

func RefreshTokens(s *service.ServiceStruct) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		logger.Info("refreshing tokens")
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

		//
		//
		//	Проверка хоста и отправка на mail
		//
		//

		// access token generation
		aToken, rToken, err := s.GenerateTokens(req.Header.Get("Host"))
		if err != nil {
			logger.Error(err)
			http.Error(res, "", http.StatusInternalServerError)
			return
		}

		// save refresh token
		err = s.DB.InsertRT(guid, rToken)
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

		rtB64 := base64.StdEncoding.EncodeToString([]byte(rToken))

		atExp := time.Now().Add(30 * time.Minute)
		rtExp := time.Now().Add(720 * time.Hour)
		http.SetCookie(res, &http.Cookie{
			Name:     "at",
			Value:    aToken,
			Expires:  atExp,
			HttpOnly: true,
		})
		http.SetCookie(res, &http.Cookie{
			Name:     "rt",
			Value:    rtB64,
			Expires:  rtExp,
			HttpOnly: true,
		})
		logger.Info("tokens have been refreshed")
	}
}
