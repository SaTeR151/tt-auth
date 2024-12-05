package service

import (
	"crypto/rand"

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
