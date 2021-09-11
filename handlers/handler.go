package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pmadhvi/telness-manager/model"
	log "github.com/sirupsen/logrus"
)

// CreateHandler is an httphandler to handle request to create an subscription
func (s Server) CreateHandler(rw http.ResponseWriter, req *http.Request) {
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Errorf("error occured when reading request body: %v", err)
	}
	var sub model.Subscription
	err = json.Unmarshal(reqBody, &sub)
	if err != nil {
		log.Errorf("error occured when unmarshalling reqBody into Subscription type: %v", err)
	}
	s.SubscriptionService.Find()
	respMsg := struct{}{} // TODO: define this later
	log.Infof("Request Body: %v", sub)
	respondSuccessJSON(rw, http.StatusCreated, respMsg)
}

// FindHandler is an httphandler to handle request to find an subscription
func FindHandler(rw http.ResponseWriter, req *http.Request) {
	// feteching the quary parameters from request url
	vars := mux.Vars(req)
	iban := vars["id"]

	respMsg := struct{}{} // TODO: define this later
	log.Infof("Request Body ID: %v", iban)
	respondSuccessJSON(rw, http.StatusOK, respMsg)
}

// UpdateHandler is an httphandler to handle request to update an subscription
func UpdateHandler(rw http.ResponseWriter, req *http.Request) {

	// feteching the quary parameters from request url
	vars := mux.Vars(req)
	iban := vars["id"]

	respMsg := struct{}{} // TODO: define this later
	log.Infof("Request Body ID: %v", iban)
	respondSuccessJSON(rw, http.StatusOK, respMsg)
}

// CheckHealthHandler is an httphandler to handle request to check application health
func CheckHealthHandler(rw http.ResponseWriter, req *http.Request) {
	respMsg := struct {
		Message string `json:"message"`
	}{
		Message: "Application is alive.",
	}
	log.Infof("Health check response: %v", respMsg)
	respondSuccessJSON(rw, 200, respMsg)
}

func respondSuccessJSON(rw http.ResponseWriter, statusCode int, response interface{}) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(response)
}

func respondErrorJSON(rw http.ResponseWriter, errorCode int, errorMsg interface{}) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(errorCode)
	json.NewEncoder(rw).Encode(errorMsg)
}
