// +build integration

package integrationtest

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pmadhvi/telness-manager/mock"

	"github.com/gorilla/mux"
	"github.com/pmadhvi/telness-manager/model"
	"github.com/stretchr/testify/assert"
)

type response struct {
	Message string
}

func requestResponse(method string, url string, requestBody []byte) (*http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, "/api/subscription", bytes.NewBuffer(requestBody))
	rw := httptest.NewRecorder()
	return req, rw
}

func mockCreateSubscription(msisdn string, now, sub_type string, status model.SubStatus) {
	mock.Create = func(sub model.CreateSubscription) error {
		return nil
	}
	mock.FindByID = func(msisdn string) (model.Subscription, error) {
		return model.Subscription{
			Msisdn:     msisdn,
			ActivateAt: now,
			SubType:    sub_type,
			Status:     status,
			CreatedAt:  now,
			ModifiedAt: now,
		}, nil
	}
	mock.GetOperator = func(msisdn string) (model.PtsResponse, error) {
		return model.PtsResponse{
			D: model.OperatorDetails{Name: "Telness AB"},
		}, nil
	}
}

func mockFindSubscription(msisdn string, now, sub_type string, status model.SubStatus) {
	mock.FindByID = func(msisdn string) (model.Subscription, error) {
		return model.Subscription{
			Msisdn:     msisdn,
			ActivateAt: now,
			SubType:    sub_type,
			Status:     status,
			CreatedAt:  now,
			ModifiedAt: now,
		}, nil
	}
	mock.GetOperator = func(msisdn string) (model.PtsResponse, error) {
		return model.PtsResponse{
			D: model.OperatorDetails{Name: "Telness AB"},
		}, nil
	}
}

func mockFindNonExistingSubscription(msisdn string) {
	mock.FindByID = func(msisdn string) (model.Subscription, error) {
		return model.Subscription{}, errors.New("subscription not found")
	}
}

func mockUpdateSubscription(msisdn string, now, sub_type string, status model.SubStatus) {
	mock.FindByID = func(msisdn string) (model.Subscription, error) {
		return model.Subscription{
			Msisdn:     msisdn,
			ActivateAt: now,
			SubType:    sub_type,
			Status:     status,
			CreatedAt:  now,
			ModifiedAt: now,
		}, nil
	}
	mock.GetOperator = func(msisdn string) (model.PtsResponse, error) {
		return model.PtsResponse{
			D: model.OperatorDetails{Name: "Telness AB"},
		}, nil
	}
	mock.Update = func(sub model.CreateSubscription) error {
		return nil
	}

}

func TestCreateSubscription(t *testing.T) {
	var (
		msisdn = "+46107500500"
		now    = time.Now().Format("2006-01-02")
	)
	request := []byte(`{
		"msisdn": "+46107500500",
		"activate_at": "2021-10-17",
		"sub_type":    "pbx",
		"status":     "pending"}`)
	req, rw := requestResponse(http.MethodPost, "/api/subscription", request)

	mockCreateSubscription(msisdn, now, "pbx", "pending")
	handler := http.HandlerFunc(server.CreateHandler)
	handler.ServeHTTP(rw, req)
	if status := rw.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}
	var resp model.Subscription
	err := json.NewDecoder(rw.Body).Decode(&resp)
	if err != nil {
		t.Errorf("could not decode response: %v", err)
	}
	assert.EqualValues(t, msisdn, resp.Msisdn)
	assert.EqualValues(t, now, resp.ActivateAt)
	assert.EqualValues(t, "pbx", resp.SubType)
	assert.EqualValues(t, "pending", resp.Status)
	assert.EqualValues(t, "Telness AB", resp.Operator)
}

func TestCreateSubscriptionEmptyStatus(t *testing.T) {
	request := []byte(`{
		"msisdn": "+46107500500",
		"activate_at": "2021-10-17",
		"sub_type":    "pbx"}`)
	req, rw := requestResponse(http.MethodPost, "/api/subscription", request)

	handler := http.HandlerFunc(server.CreateHandler)
	handler.ServeHTTP(rw, req)
	if status := rw.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
	var resp model.ErrorMessage
	err := json.NewDecoder(rw.Body).Decode(&resp)
	if err != nil {
		t.Errorf("could not decode response: %v", err)
	}
	assert.EqualValues(t, "Create request body is not valid: status cannot be empty", resp.Message)
}

func TestUpdateSubscription(t *testing.T) {
	var (
		msisdn = "+46107500501"
		now    = time.Now().Format("2006-01-02")
	)
	request := []byte(`{
		"msisdn": "+46107500501",
		"activate_at": "2021-10-17",
		"sub_type":    "cell",
		"status":     "activated"}`)
	req, rw := requestResponse(http.MethodPatch, "/api/subscription", request)

	mockUpdateSubscription(msisdn, now, "pbx", "pending")
	mockFindSubscription(msisdn, now, "cell", "activated")
	handler := http.HandlerFunc(server.UpdateHandler)
	handler.ServeHTTP(rw, req)
	if status := rw.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}
	var resp model.Subscription
	err := json.NewDecoder(rw.Body).Decode(&resp)
	if err != nil {
		t.Errorf("could not decode response: %v", err)
	}
	assert.EqualValues(t, msisdn, resp.Msisdn)
	assert.EqualValues(t, now, resp.ActivateAt)
	assert.EqualValues(t, "cell", resp.SubType)
	assert.EqualValues(t, "activated", resp.Status)
	assert.EqualValues(t, "Telness AB", resp.Operator)
}

func TestUpdateSubscriptionEmptymsisdn(t *testing.T) {
	request := []byte(`{
		"activate_at": "2021-09-17",
		"sub_type":    "pbx",
		"status": "pending"}`)
	req, rw := requestResponse(http.MethodPatch, "/api/subscription", request)

	handler := http.HandlerFunc(server.UpdateHandler)
	handler.ServeHTTP(rw, req)
	if status := rw.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
	var resp model.ErrorMessage
	err := json.NewDecoder(rw.Body).Decode(&resp)
	if err != nil {
		t.Errorf("could not decode response: %v", err)
	}
	assert.EqualValues(t, "Update request body is not valid: msisdn cannot be nil", resp.Message)
}

func TestFindSubscription(t *testing.T) {
	var (
		msisdn = "+46107500500"
		now    = time.Now().Format("2006-01-02")
	)
	req, rw := requestResponse(http.MethodGet, "/api/subscription/msisdn/", nil)
	req = mux.SetURLVars(req, map[string]string{
		"msisdn": "+46107500500",
	})
	mockFindSubscription(msisdn, now, "cell", "activated")
	handler := http.HandlerFunc(server.FindHandler)
	handler.ServeHTTP(rw, req)
	if status := rw.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}
	var resp model.Subscription
	err := json.NewDecoder(rw.Body).Decode(&resp)
	if err != nil {
		t.Errorf("could not decode response: %v", err)
	}
	assert.EqualValues(t, msisdn, resp.Msisdn)
	assert.EqualValues(t, now, resp.ActivateAt)
	assert.EqualValues(t, "cell", resp.SubType)
	assert.EqualValues(t, "activated", resp.Status)
	assert.EqualValues(t, "Telness AB", resp.Operator)
}

func TestFindSubscriptionWithEmptymsisdn(t *testing.T) {
	req, rw := requestResponse(http.MethodGet, "/api/subscription/msisdn/", nil)
	req = mux.SetURLVars(req, map[string]string{
		"msisdn": "",
	})
	handler := http.HandlerFunc(server.FindHandler)
	handler.ServeHTTP(rw, req)
	if status := rw.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
	var resp model.ErrorMessage
	err := json.NewDecoder(rw.Body).Decode(&resp)
	if err != nil {
		t.Errorf("could not decode response: %v", err)
	}
	assert.EqualValues(t, "msisdn cannot be empty", resp.Message)
}

func TestFindSubscriptionWithNonExistantmsisdn(t *testing.T) {
	req, rw := requestResponse(http.MethodGet, "/api/subscription/msisdn/", nil)
	req = mux.SetURLVars(req, map[string]string{
		"msisdn": "+46107500578",
	})
	msisdn := "+46107500578"
	mockFindNonExistingSubscription(msisdn)
	handler := http.HandlerFunc(server.FindHandler)
	handler.ServeHTTP(rw, req)
	if status := rw.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
	var resp model.ErrorMessage
	err := json.NewDecoder(rw.Body).Decode(&resp)
	if err != nil {
		t.Errorf("could not decode response: %v", err)
	}
	assert.EqualValues(t, "Could not find subscription with msisdn +46107500578, subscription not found", resp.Message)
}

func TestCancelSubscription(t *testing.T) {
	var (
		msisdn = "+46107500500"
		now    = time.Now().Format("2006-01-02")
	)

	req, rw := requestResponse(http.MethodPatch, "/api/subscription/update-subscription/msisdn/{msisdn}/status/{status}", nil)
	req = mux.SetURLVars(req, map[string]string{
		"msisdn": "+46107500500",
		"status": "cancelled",
	})
	mockUpdateSubscription(msisdn, now, "cell", "activated")
	mockFindSubscription(msisdn, "2021-10-11", "cell", "cancelled")
	handler := http.HandlerFunc(server.UpdateStatusHandler)
	handler.ServeHTTP(rw, req)
	if status := rw.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}
	var resp model.Subscription
	err := json.NewDecoder(rw.Body).Decode(&resp)
	if err != nil {
		t.Errorf("could not decode response: %v", err)
	}
	assert.EqualValues(t, msisdn, resp.Msisdn)
	assert.EqualValues(t, "2021-10-11", resp.ActivateAt)
	assert.EqualValues(t, "cell", resp.SubType)
	assert.EqualValues(t, "cancelled", resp.Status)
	assert.EqualValues(t, "Telness AB", resp.Operator)
}

func TestPauseSubscriptionWithEmptymsisdn(t *testing.T) {
	req, rw := requestResponse(http.MethodPost, "/api/subscription/update-subscription/msisdn/{msisdn}/status/{status}", nil)
	req = mux.SetURLVars(req, map[string]string{
		"msisdn": "",
		"status": "paused",
	})
	handler := http.HandlerFunc(server.UpdateStatusHandler)
	handler.ServeHTTP(rw, req)
	if status := rw.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
	var resp model.ErrorMessage
	err := json.NewDecoder(rw.Body).Decode(&resp)
	if err != nil {
		t.Errorf("could not decode response: %v", err)
	}
	assert.EqualValues(t, "msisdn cannot be empty", resp.Message)
}

func TestReactivateSubscriptionWithNonExistingmsisdn(t *testing.T) {
	req, rw := requestResponse(http.MethodPost, "/api/subscription", nil)
	req = mux.SetURLVars(req, map[string]string{
		"msisdn": "+46107500500",
		"status": "activated",
	})
	msisdn := "+46107500500"
	mockFindNonExistingSubscription(msisdn)
	handler := http.HandlerFunc(server.UpdateStatusHandler)
	handler.ServeHTTP(rw, req)
	if status := rw.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
	var resp model.ErrorMessage
	err := json.NewDecoder(rw.Body).Decode(&resp)
	if err != nil {
		t.Errorf("could not decode response: %v", err)
	}
	assert.EqualValues(t, "Could not find subscription with msisdn +46107500500, subscription not found", resp.Message)
}

func TestUpdateActivationDate(t *testing.T) {
	var (
		msisdn = "+46107500500"
		now    = time.Now().Format("2006-01-02")
	)

	req, rw := requestResponse(http.MethodPatch, "/api/subscription/update-activation-date/msisdn/{msisdn}/date/{date}", nil)
	req = mux.SetURLVars(req, map[string]string{
		"msisdn": "+46107500500",
		"date":   "2021-10-11",
	})
	mockUpdateSubscription(msisdn, now, "cell", "pending")
	mockFindSubscription(msisdn, "2021-10-11", "cell", "pending")
	handler := http.HandlerFunc(server.UpdateActivationDateHandler)
	handler.ServeHTTP(rw, req)
	if status := rw.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}
	var resp model.Subscription
	err := json.NewDecoder(rw.Body).Decode(&resp)
	if err != nil {
		t.Errorf("could not decode response: %v", err)
	}
	assert.EqualValues(t, msisdn, resp.Msisdn)
	assert.EqualValues(t, "2021-10-11", resp.ActivateAt)
	assert.EqualValues(t, "cell", resp.SubType)
	assert.EqualValues(t, "pending", resp.Status)
	assert.EqualValues(t, "Telness AB", resp.Operator)
}

func TestUpdateActivationDateWithWrongDate(t *testing.T) {
	req, rw := requestResponse(http.MethodPatch, "/api/subscription/update-activation-date/msisdn/{msisdn}/date/{date}", nil)
	req = mux.SetURLVars(req, map[string]string{
		"msisdn": "+46107500500",
		"date":   "2021-09-11",
	})
	handler := http.HandlerFunc(server.UpdateActivationDateHandler)
	handler.ServeHTTP(rw, req)
	if status := rw.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
	var resp model.ErrorMessage
	err := json.NewDecoder(rw.Body).Decode(&resp)
	if err != nil {
		t.Errorf("could not decode response: %v", err)
	}
	assert.EqualValues(t, "enter valid future date for activation", resp.Message)
}

func TestUpdateActivationDateWithEmptyDate(t *testing.T) {
	req, rw := requestResponse(http.MethodPatch, "/api/subscription/update-activation-date/msisdn/{msisdn}/date/{date}", nil)
	req = mux.SetURLVars(req, map[string]string{
		"msisdn": "+46107500500",
		"date":   "",
	})
	handler := http.HandlerFunc(server.UpdateActivationDateHandler)
	handler.ServeHTTP(rw, req)
	if status := rw.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
	var resp model.ErrorMessage
	err := json.NewDecoder(rw.Body).Decode(&resp)
	if err != nil {
		t.Errorf("could not decode response: %v", err)
	}
	assert.EqualValues(t, "enter valid date for activation, date cannot be empty", resp.Message)
}
