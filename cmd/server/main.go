package main

import (
	"fmt"
	"net/http"

	"github.com/nitotang/obtener-materiales/internal/service"
	transportHTTP "github.com/nitotang/obtener-materiales/internal/transport/http"
)

// App - the struct which contains things like
// pointers to database connections
type App struct{}

func (app *App) Run() error {
	fmt.Println("Setting Up Our App")

	bankService := service.NewService()
	handler := transportHTTP.NewHandler(bankService)
	handler.SetupRoutes()

	if err := http.ListenAndServe(":8080", handler.Router); err != nil {
		fmt.Println("Failed to set up server")
		return err
	}

	return nil
}

func main() {
	fmt.Println("SOAP Call test")

	app := App{}
	if err := app.Run(); err != nil {
		fmt.Println("Error Starting Up")
		fmt.Println(err)
	}
}
