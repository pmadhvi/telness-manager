package handlers

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/pmadhvi/telness-manager/model"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	Log                 *log.Logger
	Port                string
	SubscriptionService subscriptionService
}

type subscriptionService interface {
	Create(sub *model.Subscription) error
	FindbyID(id uuid.UUID) (model.Subscription, error)
	Update(sub *model.Subscription) (model.Subscription, error)
}

//  defines routes and their handlers and start the server
func (s Server) Start() error {
	// Initialize mux router
	router := mux.NewRouter()

	// define routes and call their handler function
	router.HandleFunc("/api/subscription/health", CheckHealthHandler)
	router.HandleFunc("/api/subscription/{id}", FindHandler)
	//router.HandleFunc("/api/subscription/", ValidateIbanHandler).Methods(http.Post)
	//router.HandleFunc("/api/subscription/", ValidateIbanHandler).Method(Patch)

	// start the server on specified port
	err := http.ListenAndServe(fmt.Sprintf(":%s", s.Port), router)
	log.Errorf("error starting server: %v", err)
	return err
}
