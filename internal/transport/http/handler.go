package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nitotang/obtener-materiales/internal/service"
)

// Handler - stores pointer to bank service
type Handler struct {
	Router  *mux.Router
	Service service.Service
}

// Response - an object to store responses from our API
type Response struct {
	Message string
	Error   string
}

// NewHandler - returns a pointer to a Handler
func NewHandler(service *service.Service) *Handler {
	return &Handler{
		Service: *service,
	}
}

func (h *Handler) SetupRoutes() {
	fmt.Println("Setting Up Routes")
	h.Router = mux.NewRouter()

	h.Router.HandleFunc("/api/material/{id}/stock", h.GetStock).Methods("GET")

	h.Router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Alive!!")
	})
}

// GetStock - consulta stock de un material
func (h *Handler) GetStock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	vars := mux.Vars(r)
	id := vars["id"]

	response, err := h.Service.GetStock(id)
	if err != nil {
		sendErrorResponse(w, "Error", err)
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}

}

func sendErrorResponse(w http.ResponseWriter, message string, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	if err := json.NewEncoder(w).Encode(Response{Message: message, Error: err.Error()}); err != nil {
		panic(err)
	}
}
