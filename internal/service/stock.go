package service

import (
	"fmt"

	"obtener-materiales/internal/soapHandler"
)

type Service struct{}

type Stock struct {
	ID         string
	Name       string
	Code       string
	Address    string
	PostalCode string
}

type StockService interface {
	GetStock(ID string) (Stock, error)
}

// NewService - returns a new comment service
func NewService() *Service {
	return &Service{}
}

// GetComment - retrieves comments by their ID from the database
func (s *Service) GetStock(ID string) (*soapHandler.Response, error) {

	fmt.Println("ID: " + ID)

	soapRequest := soapHandler.Request{}
	soapRequest.CodigoProducto = ID

	response, err := soapHandler.CallSOAPClientSteps(&soapRequest)

	if err != nil {
		return response, err
	}

	return response, nil
}
