package service

import (
	"fmt"

	"github.com/nitotang/obtener-materiales/internal/soapHandler"
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

	fmt.Println("ID: " + ID)

	soapRequest := soapHandler.Request{}
	soapRequest.Codigo = bank.ID

	soapResponse, err := soapHandler.CallSOAPClientSteps(&soapRequest)

	if err != nil {
		return Bank{}, err
	}

	bank.Name = soapResponse.SoapBody.Resp.Result.ResultItem.ResultMessage

	return bank, nil
}
