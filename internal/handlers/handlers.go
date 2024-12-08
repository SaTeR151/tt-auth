package handlers

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/sater-151/tt-auth/internal/database"
	"github.com/sater-151/tt-auth/internal/service"
	"github.com/sater-151/tt-auth/internal/utils"
	logger "github.com/sirupsen/logrus"
)

var ErrGUIDRequired = errors.New("guid required")

func GetTokens(s service.ServiceInterface) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		logger.Info("getting tokens")
		guid := req.FormValue("guid")
		if guid == "" {
			logger.Error(ErrGUIDRequired)
			http.Error(res, ErrGUIDRequired.Error(), http.StatusBadRequest)
			return
		}

		aToken, rToken, err := utils.GenerateTokens(req.Host)
		if err != nil {
			logger.Error(err)
			http.Error(res, "", http.StatusInternalServerError)
			return
		}
		// save refresh token
		err = s.InsertRT(guid, rToken)
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
		atTimeExp, err := strconv.Atoi(os.Getenv("ATEXPIRES"))
		if err != nil {
			logger.Error(err)
			http.Error(res, "", http.StatusInternalServerError)
			return
		}
		rtTimeExp, err := strconv.Atoi(os.Getenv("RTEXPIRES"))
		if err != nil {
			logger.Error(err)
			http.Error(res, "", http.StatusInternalServerError)
			return
		}
		atExp := time.Now().Add(time.Duration(atTimeExp) * time.Second)
		rtExp := time.Now().Add(time.Duration(rtTimeExp) * time.Second)
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

func RefreshTokens(s service.ServiceInterface) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		logger.Info("refreshing tokens")
		guid := req.FormValue("guid")
		if guid == "" {
			logger.Error(ErrGUIDRequired)
			http.Error(res, ErrGUIDRequired.Error(), http.StatusBadRequest)
			return
		}

		rtCookie, err := req.Cookie("rt")
		if err != nil {
			logger.Error(err)
			http.Error(res, "", http.StatusUnauthorized)
			return
		}
		gettingRTBase64, err := base64.StdEncoding.DecodeString(rtCookie.Value)
		if err != nil {
			logger.Error(err)
			http.Error(res, "", http.StatusInternalServerError)
			return
		}

		logger.Debug("starting generate tokens")
		aToken, rToken, err := utils.GenerateTokens(req.Host)
		if err != nil {
			logger.Error(err)
			http.Error(res, "", http.StatusInternalServerError)
			return
		}

		logger.Debug("comparing refresh tokens")
		comp, err := s.CompareRT(string(gettingRTBase64), guid)
		if err != nil {
			fmt.Println(3)
			logger.Error(err)
			http.Error(res, "", http.StatusInternalServerError)
			return
		}
		if !comp {
			logger.Error(database.ErrUnauthorized)
			http.Error(res, "", http.StatusUnauthorized)
			return
		}
		atCook, err := req.Cookie("at")
		if err != nil {
			logger.Error(err)
			http.Error(res, "", http.StatusUnauthorized)
			return
		}

		logger.Debug("checking host")
		ok, err := utils.CheckHost(atCook.Value, req.Host)
		if err != nil {
			logger.Error(err)
			http.Error(res, "", http.StatusInternalServerError)
			return
		}
		if !ok {
			logger.Warn("another ip")
			s.EmailWarning(guid)
		}

		// save refresh token
		err = s.InsertRT(guid, rToken)
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

		atTimeExp, err := strconv.Atoi(os.Getenv("ATEXPIRES"))
		if err != nil {
			logger.Error(err)
			http.Error(res, "", http.StatusInternalServerError)
			return
		}
		rtTimeExp, err := strconv.Atoi(os.Getenv("RTEXPIRES"))
		if err != nil {
			logger.Error(err)
			http.Error(res, "", http.StatusInternalServerError)
			return
		}
		atExp := time.Now().Add(time.Duration(atTimeExp) * time.Second)
		rtExp := time.Now().Add(time.Duration(rtTimeExp) * time.Second)

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
