package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	//"time"

	"github.com/gorilla/mux"
	"github.com/pmadhvi/telness-manager/model"
)

// CreateHandler is an httphandler to handle request to create an subscription
func (s Server) CreateHandler(rw http.ResponseWriter, req *http.Request) {
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		msg := fmt.Sprintf("Could not read request body: %v", err)
		s.Log.Error(msg)
		returnError(rw, msg, 400)
		return
	}
	var subreq model.CreateSubscription
	err = json.Unmarshal(reqBody, &subreq)
	if err != nil {
		msg := fmt.Sprintf("Could not unmarshalling request body into Subscription type: %v", err)
		s.Log.Error(msg)
		returnError(rw, msg, 400)
		return
	}
	err = validateRequest(subreq)
	if err != nil {
		msg := fmt.Sprintf("Create request body is not valid: %v", err)
		s.Log.Error(msg)
		returnError(rw, msg, 400)
		return
	}

	var sub model.Subscription
	sub, err = s.SubscriptionService.Create(subreq)
	if err != nil {
		msg := fmt.Sprintf("Could not create a new subscription, %v", err)
		s.Log.Error(msg)
		returnError(rw, msg, 400)
		return
	}
	respondSuccessJSON(rw, http.StatusCreated, sub)
}

// UpdateHandler is an httphandler to handle request to create an subscription
func (s Server) UpdateHandler(rw http.ResponseWriter, req *http.Request) {
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		msg := fmt.Sprintf("Could not read request body: %v", err)
		s.Log.Error(msg)
		returnError(rw, msg, 400)
		return
	}
	var subreq model.CreateSubscription
	err = json.Unmarshal(reqBody, &subreq)
	if err != nil {
		msg := fmt.Sprintf("Could not unmarshalling reqBody into Subscription type: %v", err)
		s.Log.Error(msg)
		returnError(rw, msg, 400)
		return
	}
	err = validateRequest(subreq)
	if err != nil {
		msg := fmt.Sprintf("Update request body is not valid: %v", err)
		s.Log.Error(msg)
		returnError(rw, msg, 400)
		return
	}

	var sub model.Subscription
	sub, err = s.SubscriptionService.Update(subreq)
	if err != nil {
		msg := fmt.Sprintf("Could not update subscription: %v", err)
		s.Log.Error(msg)
		returnError(rw, msg, 400)
		return
	}
	respondSuccessJSON(rw, http.StatusOK, sub)
}

// FindHandler is an httphandler to handle request to find an subscription
func (s Server) FindHandler(rw http.ResponseWriter, req *http.Request) {
	// feteching the quary parameters from request url
	vars := mux.Vars(req)
	msisdn := vars["msisdn"]
	if msisdn == "" {
		s.Log.Error("msisdn cannot be empty")
		returnError(rw, "msisdn cannot be empty", 400)
		return
	}
	var sub model.Subscription
	sub, err := s.SubscriptionService.FindbyID(msisdn)
	if err != nil {
		msg := fmt.Sprintf("Could not find subscription with msisdn %v, %v", msisdn, err)
		s.Log.Error(msg)
		returnError(rw, msg, 404)
		return
	}
	respondSuccessJSON(rw, http.StatusOK, sub)
}

// UpdateStatusHandler is an httphandler to handle request to find an subscription and updates it status
func (s Server) UpdateStatusHandler(rw http.ResponseWriter, req *http.Request) {
	// feteching the quary parameters from request url and validating it
	vars := mux.Vars(req)
	msisdn := vars["msisdn"]
	status := vars["status"]
	if msisdn == "" {
		msg := fmt.Sprint("msisdn cannot be empty")
		s.Log.Error(msg)
		returnError(rw, msg, 400)
		return
	} else if status == "" {
		msg := fmt.Sprint("status cannot be empty")
		s.Log.Error(msg)
		returnError(rw, msg, 400)
		return
	} else if !IsValidStatus(model.SubStatus(status)) {
		msg := fmt.Sprintf("Invalid status type %v", status)
		s.Log.Error(msg)
		returnError(rw, msg, 400)
		return
	}

	var sub model.Subscription
	foundSub, err := s.SubscriptionService.FindbyID(msisdn)
	if err != nil {
		msg := fmt.Sprintf("Could not find subscription with msisdn %v, %v", msisdn, err)
		s.Log.Error(msg)
		returnError(rw, msg, 404)
		return
	}
	updateSub := model.CreateSubscription{
		Msisdn:     foundSub.Msisdn,
		ActivateAt: foundSub.ActivateAt,
		SubType:    foundSub.SubType,
		Status:     model.SubStatus(status),
	}
	sub, err = s.SubscriptionService.Update(updateSub)
	if err != nil {
		msg := fmt.Sprintf("Could not update subscription: %v", err)
		s.Log.Error(msg)
		returnError(rw, msg, 400)
		return
	}
	respondSuccessJSON(rw, http.StatusOK, sub)
}

// UpdateActivationDateHandler is an httphandler to handle request to update activation date of pending an subscription
func (s Server) UpdateActivationDateHandler(rw http.ResponseWriter, req *http.Request) {
	// feteching the quary parameters from request url and validate it
	vars := mux.Vars(req)
	msisdn := vars["msisdn"]
	date := vars["date"]

	if msisdn == "" {
		msg := fmt.Sprint("msisdn cannot be empty")
		s.Log.Error(msg)
		returnError(rw, msg, 400)
		return
	} else if date == "" {
		msg := fmt.Sprint("enter valid date for activation, date cannot be empty")
		s.Log.Error(msg)
		returnError(rw, msg, 400)
		return
	}

	// Check if activation date is future date
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		msg := fmt.Sprint("could not parse string date into time.Time format")
		s.Log.Error(msg)
		returnError(rw, msg, 400)
		return
	}

	if parsedDate.Before(time.Now()) {
		msg := fmt.Sprint("enter valid future date for activation")
		s.Log.Error(msg)
		returnError(rw, msg, 400)
		return
	}

	var sub model.Subscription
	foundSub, err := s.SubscriptionService.FindbyID(msisdn)
	if err != nil {
		msg := fmt.Sprintf("Could not find subscription with msisdn %v, %v", msisdn, err)
		s.Log.Error(msg)
		returnError(rw, msg, 404)
		return
	}
	if foundSub.Status != model.StatusPending {
		msg := fmt.Sprintf("Found subscription status %v is not pending to update activation date", foundSub.Status)
		s.Log.Error(msg)
		returnError(rw, msg, 400)
		return
	}
	updateSub := model.CreateSubscription{
		Msisdn:     foundSub.Msisdn,
		ActivateAt: date,
		SubType:    foundSub.SubType,
		Status:     foundSub.Status,
	}
	sub, err = s.SubscriptionService.Update(updateSub)
	if err != nil {
		msg := fmt.Sprintf("Could not update subscription: %v", err)
		s.Log.Error(msg)
		returnError(rw, msg, 400)
		return
	}
	respondSuccessJSON(rw, http.StatusOK, sub)
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
	rw.WriteHeader(statusCode)
	json.NewEncoder(rw).Encode(response)
}

func respondErrorJSON(rw http.ResponseWriter, errorCode int, errorMsg interface{}) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(errorCode)
	json.NewEncoder(rw).Encode(errorMsg)
}

func returnError(rw http.ResponseWriter, message string, statusCode int) {
	respMsg := model.ErrorMessage{
		Message: message,
	}
	respondErrorJSON(rw, statusCode, respMsg)
}

func validateRequest(sub model.CreateSubscription) error {
	activate_at, err := time.Parse("2006-01-02", sub.ActivateAt)
	if err != nil {
		return errors.New("could not parse string activate_at into time.Time format")
	}

	var regexp = regexp.MustCompile(`^\+46[1-9][0-9]{8}$`)
	if sub.Msisdn == "" {
		return errors.New("msisdn cannot be nil")
	} else if !regexp.MatchString(sub.Msisdn) {
		return errors.New("msisdn must be of format: +46 followed by 9 digits of phone number, example - [+46107500500]")
	} else if sub.ActivateAt == "" {
		return errors.New("activate_at cannot be empty")
	} else if activate_at.Before(time.Now()) {
		return errors.New("activate_at should be future date")
	} else if sub.SubType == "" {
		return errors.New("sub_type cannot be empty")
	} else if sub.Status == "" {
		return errors.New("status cannot be empty")
	} else if !IsValidStatus(sub.Status) {
		return errors.New("Invalid status type")
	}

	return nil
}

func IsValidStatus(status model.SubStatus) bool {
	switch status {
	case model.StatusPending, model.StatusPaused, model.StatusActivated, model.StatusCancelled:
		return true
	}
	return false
}
