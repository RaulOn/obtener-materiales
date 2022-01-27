package service

import (
	"github.com/nitotang/obtener-materiales/internal/soaphandler"
)

type Service struct{}

type Bank struct {
	ID         string
	Name       string
	Code       string
	Address    string
	PostalCode string
}

type BankService interface {
	GetBank(ID string) (Bank, error)
}

// NewService - returns a new comment service
func NewService() *Service {
	return &Service{}
}

// GetComment - retrieves comments by their ID from the database
func (s *Service) GetBank(ID string) (Bank, error) {
	var bank Bank
	bank.ID = ID

	soapRequest := soaphandler.Request{}
	soapRequest.Blz = bank.ID

	soapResponse, err := soaphandler.CallSOAPClientSteps(&soapRequest)

	if err != nil {
		return Bank{}, err
	}

	bank.Name = soapResponse.SoapBody.Resp.Response.Bezeichnung
	bank.Code = soapResponse.SoapBody.Resp.Response.Bic
	bank.Address = soapResponse.SoapBody.Resp.Response.Ort
	bank.PostalCode = soapResponse.SoapBody.Resp.Response.Plz

	return bank, nil
}
