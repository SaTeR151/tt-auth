package service

import (
	"github.com/sater-151/tt-auth/internal/database"
	"github.com/sater-151/tt-auth/internal/utils"
)

type ServiceInterface interface {
	EmailWarning(guid string) error
	InsertRT(guid, rToken string) error
	CompareRT(rtb, guid string) (bool, error)
}

type ServiceStruct struct {
	DB database.DBInterface
}

func New(db database.DBInterface) *ServiceStruct {
	service := &ServiceStruct{DB: db}
	return service
}
func (s *ServiceStruct) EmailWarning(guid string) error {
	mail, err := s.DB.SelectMail(guid)
	if err != nil {
		return err
	}
	err = utils.SendMasseg(mail)
	if err != nil {
		return err
	}
	return nil
}

func (s *ServiceStruct) InsertRT(guid, rToken string) error {
	err := s.DB.UpdateRT(guid, rToken)
	return err
}

func (s *ServiceStruct) CompareRT(rtB64, guid string) (bool, error) {
	rTokenBcryptNew, err := s.DB.GetBcrypt(rtB64)
	if err != nil {
		return false, err
	}
	rTokenOld, err := s.DB.GetToken(guid)
	if err != nil {
		return false, err
	}
	if rTokenBcryptNew != rTokenOld {
		return false, nil
	}
	return true, nil
}
