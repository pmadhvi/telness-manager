package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pmadhvi/telness-manager/model"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	Log                 *log.Logger
	Port                string
	SubscriptionService SubscriptionService
}

type SubscriptionService interface {
	Create(sub model.CreateSubscription) (model.Subscription, error)
	FindbyID(id string) (model.Subscription, error)
	Update(sub model.CreateSubscription) (model.Subscription, error)
}

//  defines routes and their handlers and start the server
func (s Server) Start() error {
	log.Info("Telness server is starting up")
	// Initialize mux router
	router := mux.NewRouter()

	// define routes and call their handler function
	router.HandleFunc("/api/subscription/health", s.CheckHealthHandler)
	router.HandleFunc("/api/subscription/msisdn/{msisdn}", s.FindHandler).Methods("Get")
	router.HandleFunc("/api/subscription", s.CreateHandler).Methods("Post")
	router.HandleFunc("/api/subscription", s.UpdateHandler).Methods("Patch")
	router.HandleFunc("/api/subscription/update-subscription/msisdn/{msisdn}/status/{status}", s.UpdateStatusHandler).Methods("Patch")
	router.HandleFunc("/api/subscription/update-activation-date/msisdn/{msisdn}/date/{date}", s.UpdateActivationDateHandler).Methods("Patch")

	// start the server on specified port
	err := http.ListenAndServe(fmt.Sprintf(":%s", s.Port), router)
	log.Errorf("error starting server: %v", err)
	return err
}
