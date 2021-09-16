package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/pmadhvi/telness-manager/model"
)

// CreateHandler is an httphandler to handle request to create an subscription
func (s Server) CreateHandler(rw http.ResponseWriter, req *http.Request) {
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		s.Log.Errorf("error occured when reading request body: %v", err)
	}
	var subreq model.CreateSubscription
	err = json.Unmarshal(reqBody, &subreq)
	if err != nil {
		s.Log.Errorf("error occured when unmarshalling reqBody into Subscription type: %v", err)
	}
	var sub model.Subscription
	sub, err = s.SubscriptionService.Create(subreq)
	if err != nil {
		s.Log.Errorf("Could not create a new subscription, %v", err)
		respMsg := model.ErrorMessage{
			Message: "Could not create a new subscription",
		}
		respondErrorJSON(rw, 400, respMsg)
		return
	}
	respondSuccessJSON(rw, http.StatusCreated, sub)
}

// UpdateHandler is an httphandler to handle request to create an subscription
func (s Server) UpdateHandler(rw http.ResponseWriter, req *http.Request) {
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		s.Log.Errorf("error occured when reading request body: %v", err)
	}
	var subreq model.CreateSubscription
	err = json.Unmarshal(reqBody, &subreq)
	if err != nil {
		s.Log.Errorf("error occured when unmarshalling reqBody into Subscription type: %v", err)
	}
	var sub model.Subscription
	sub, err = s.SubscriptionService.Update(subreq)
	if err != nil {
		s.Log.Errorf("Could not update subscription with msidn %v, %v", err)
		respMsg := model.ErrorMessage{
			Message: "Could not update subscription",
		}
		respondErrorJSON(rw, 400, respMsg)
		return
	}
	respondSuccessJSON(rw, http.StatusOK, sub)
}

// FindHandler is an httphandler to handle request to find an subscription
func (s Server) FindHandler(rw http.ResponseWriter, req *http.Request) {
	// feteching the quary parameters from request url
	vars := mux.Vars(req)
	msidn := vars["msidn"]
	uuidMsidn := uuid.MustParse(msidn)

	var sub model.Subscription
	sub, err := s.SubscriptionService.FindbyID(uuidMsidn)
	if err != nil {
		s.Log.Errorf("Could not find subscription with msidn %v, %v", err)
		respMsg := model.ErrorMessage{
			Message: "Could not find subscription",
		}
		respondErrorJSON(rw, 400, respMsg)
		return
	}
	respondSuccessJSON(rw, http.StatusOK, sub)
}

// CancelHandler is an httphandler to handle request to find an subscription and cancel it
func (s Server) CancelHandler(rw http.ResponseWriter, req *http.Request) {
	// feteching the quary parameters from request url
	vars := mux.Vars(req)
	msidn := vars["msidn"]
	uuidMsidn := uuid.MustParse(msidn)

	var sub model.Subscription
	foundSub, err := s.SubscriptionService.FindbyID(uuidMsidn)
	if err != nil {
		s.Log.Errorf("Could not find subscription with msidn for cancelling it: %v, %v", err)
		respMsg := model.ErrorMessage{
			Message: "Could not find subscription for cancelling it",
		}
		respondErrorJSON(rw, 400, respMsg)
		return
	}
	updateSub := model.CreateSubscription{
		Msidn:      foundSub.Msidn,
		ActivateAt: foundSub.ActivateAt,
		SubType:    foundSub.SubType,
		Status:     model.StatusCancelled,
	}
	sub, err = s.SubscriptionService.Update(updateSub)
	if err != nil {
		s.Log.Errorf("Could not cancel subscription with msidn %v, %v", err)
		respMsg := model.ErrorMessage{
			Message: "Could not cancel subscription",
		}
		respondErrorJSON(rw, 400, respMsg)
		return
	}
	respondSuccessJSON(rw, http.StatusOK, sub)
}

// PauseHandler is an httphandler to handle request to find an subscription and pause it
func (s Server) PauseHandler(rw http.ResponseWriter, req *http.Request) {
	// feteching the quary parameters from request url
	vars := mux.Vars(req)
	msidn := vars["msidn"]
	uuidMsidn := uuid.MustParse(msidn)

	var sub model.Subscription
	foundSub, err := s.SubscriptionService.FindbyID(uuidMsidn)
	if err != nil {
		s.Log.Errorf("Could not find subscription with msidn for pausing it: %v, %v", err)
		respMsg := model.ErrorMessage{
			Message: "Could not find subscription for pausing",
		}
		respondErrorJSON(rw, 400, respMsg)
		return
	}
	updateSub := model.CreateSubscription{
		Msidn:      foundSub.Msidn,
		ActivateAt: foundSub.ActivateAt,
		SubType:    foundSub.SubType,
		Status:     model.StatusPaused,
	}
	sub, err = s.SubscriptionService.Update(updateSub)
	if err != nil {
		s.Log.Errorf("Could not pause subscription with msidn %v, %v", err)
		respMsg := model.ErrorMessage{
			Message: "Could not pause subscription",
		}
		respondErrorJSON(rw, 400, respMsg)
		return
	}
	respondSuccessJSON(rw, http.StatusOK, sub)
}

// ReactivateHandler is an httphandler to handle request to find an subscription and reactivate it
func (s Server) ReactivateHandler(rw http.ResponseWriter, req *http.Request) {
	// feteching the quary parameters from request url
	vars := mux.Vars(req)
	msidn := vars["msidn"]
	uuidMsidn := uuid.MustParse(msidn)

	var sub model.Subscription
	foundSub, err := s.SubscriptionService.FindbyID(uuidMsidn)
	if err != nil {
		s.Log.Errorf("Could not find subscription with msidn for reactivating it: %v, %v", err)
		respMsg := model.ErrorMessage{
			Message: "Could not find subscription for reactivating",
		}
		respondErrorJSON(rw, 400, respMsg)
		return
	}
	updateSub := model.CreateSubscription{
		Msidn:      foundSub.Msidn,
		ActivateAt: foundSub.ActivateAt,
		SubType:    foundSub.SubType,
		Status:     model.StatusActivated,
	}
	sub, err = s.SubscriptionService.Update(updateSub)
	if err != nil {
		s.Log.Errorf("Could not reactivate subscription with msidn %v, %v", err)
		respMsg := model.ErrorMessage{
			Message: "Could not reactivate subscription",
		}
		respondErrorJSON(rw, 400, respMsg)
		return
	}
	respondSuccessJSON(rw, http.StatusOK, sub)
}

// UpdateActivationDateHandler is an httphandler to handle request to update activation date of pending an subscription
func (s Server) UpdateActivationDateHandler(rw http.ResponseWriter, req *http.Request) {
	// feteching the quary parameters from request url
	vars := mux.Vars(req)
	msidn := vars["msidn"]
	uuidMsidn := uuid.MustParse(msidn)

	date := vars["date"]
	if date == "" {
		s.Log.Error("enter valid date for activation")
		respMsg := model.ErrorMessage{
			Message: "Not a valid date for activation",
		}
		respondErrorJSON(rw, 400, respMsg)
		return
	}
	var sub model.Subscription
	foundSub, err := s.SubscriptionService.FindbyID(uuidMsidn)
	if err != nil {
		s.Log.Errorf("Could not find subscription with msidn %v for updating the activation date: %v", err)
		respMsg := model.ErrorMessage{
			Message: "Could not find subscription for updating activation date",
		}
		respondErrorJSON(rw, 400, respMsg)
		return
	}
	if foundSub.Status == model.StatusPending {
		updateSub := model.CreateSubscription{
			Msidn:      foundSub.Msidn,
			ActivateAt: date,
			SubType:    foundSub.SubType,
			Status:     model.StatusPending,
		}
		sub, err = s.SubscriptionService.Update(updateSub)
		if err != nil {
			s.Log.Errorf("Could not reactivate subscription with msidn %v, %v", err)
			respMsg := model.ErrorMessage{
				Message: "Could not reactivate subscription",
			}
			respondErrorJSON(rw, 400, respMsg)
			return
		}
		respondSuccessJSON(rw, http.StatusOK, sub)
		return
	}
	respMsg := model.ErrorMessage{
		Message: "Found subscription status is not pending to update activation date",
	}
	respondErrorJSON(rw, 400, respMsg)
}

// CheckHealthHandler is an httphandler to handle request to check application health
func (s Server) CheckHealthHandler(rw http.ResponseWriter, req *http.Request) {
	respMsg := struct {
		Message string `json:"message"`
	}{
		Message: "Application is alive.",
	}
	s.Log.Infof("Health check response: %v", respMsg)
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
