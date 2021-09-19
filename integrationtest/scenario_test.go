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

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/pmadhvi/telness-manager/mock"

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

func mockCreateSubscription(msidn uuid.UUID, now, sub_type string, status model.SubStatus) {
	mock.Create = func(sub model.CreateSubscription) error {
		return nil
	}
	mock.FindByID = func(msidn uuid.UUID) (model.Subscription, error) {
		return model.Subscription{
			Msidn:      msidn,
			ActivateAt: now,
			SubType:    sub_type,
			Status:     status,
			CreatedAt:  now,
			ModifiedAt: now,
		}, nil
	}
}

func mockFindSubscription(msidn uuid.UUID, now, sub_type string, status model.SubStatus) {
	mock.FindByID = func(msidn uuid.UUID) (model.Subscription, error) {
		return model.Subscription{
			Msidn:      msidn,
			ActivateAt: now,
			SubType:    sub_type,
			Status:     status,
			CreatedAt:  now,
			ModifiedAt: now,
		}, nil
	}
}

func mockFindNonExistingSubscription(msidn uuid.UUID) {
	mock.FindByID = func(msidn uuid.UUID) (model.Subscription, error) {
		return model.Subscription{}, errors.New("subscription not found")
	}
}

func mockUpdateSubscription(msidn uuid.UUID, now, sub_type string, status model.SubStatus) {
	mock.FindByID = func(msidn uuid.UUID) (model.Subscription, error) {
		return model.Subscription{
			Msidn:      msidn,
			ActivateAt: now,
			SubType:    sub_type,
			Status:     status,
			CreatedAt:  now,
			ModifiedAt: now,
		}, nil
	}

	mock.Update = func(sub model.CreateSubscription) error {
		return nil
	}

}

func TestCreateSubscription(t *testing.T) {
	var (
		msidnUUID uuid.UUID
		now       = time.Now().Format("2006-01-02")
	)
	msidnUUID, _ = uuid.Parse("c019ecde-17cb-4ef8-8a7d-85937a9250ed")
	request := []byte(`{
		"msidn": "c019ecde-17cb-4ef8-8a7d-85937a9250ed",
		"activate_at": "2021-10-17",
		"sub_type":    "pbx",
		"status":     "pending"}`)
	req, rw := requestResponse(http.MethodPost, "/api/subscription", request)

	mockCreateSubscription(msidnUUID, now, "pbx", "pending")
	handler := http.HandlerFunc(server.CreateHandler)
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
	assert.EqualValues(t, msidnUUID, resp.Msidn)
	assert.EqualValues(t, now, resp.ActivateAt)
	assert.EqualValues(t, "pbx", resp.SubType)
	assert.EqualValues(t, "pending", resp.Status)
}

func TestCreateSubscriptionEmptyStatus(t *testing.T) {
	request := []byte(`{
		"msidn": "c019ecde-17cb-4ef8-8a7d-85937a9250ed",
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
		msidnUUID uuid.UUID
		now       = time.Now().Format("2006-01-02")
	)
	msidnUUID, _ = uuid.Parse("c019ecde-17cb-4ef8-8a7d-85937a9250ed")
	request := []byte(`{
		"msidn": "c019ecde-17cb-4ef8-8a7d-85937a9250ed",
		"activate_at": "2021-10-17",
		"sub_type":    "cell",
		"status":     "activated"}`)
	req, rw := requestResponse(http.MethodPatch, "/api/subscription", request)

	mockUpdateSubscription(msidnUUID, now, "pbx", "pending")
	mockFindSubscription(msidnUUID, now, "cell", "activated")
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
	assert.EqualValues(t, msidnUUID, resp.Msidn)
	assert.EqualValues(t, now, resp.ActivateAt)
	assert.EqualValues(t, "cell", resp.SubType)
	assert.EqualValues(t, "activated", resp.Status)
}

func TestUpdateSubscriptionEmptyMsidn(t *testing.T) {
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
	assert.EqualValues(t, "Update request body is not valid: msidn cannot be nil", resp.Message)
}

func TestFindSubscription(t *testing.T) {
	var (
		msidnUUID uuid.UUID
		now       = time.Now().Format("2006-01-02")
	)
	msidnUUID, _ = uuid.Parse("c019ecde-17cb-4ef8-8a7d-85937a9250ed")
	req, rw := requestResponse(http.MethodGet, "/api/subscription/msidn/", nil)
	req = mux.SetURLVars(req, map[string]string{
		"msidn": "c019ecde-17cb-4ef8-8a7d-85937a9250ed",
	})
	mockFindSubscription(msidnUUID, now, "cell", "activated")
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
	assert.EqualValues(t, msidnUUID, resp.Msidn)
	assert.EqualValues(t, now, resp.ActivateAt)
	assert.EqualValues(t, "cell", resp.SubType)
	assert.EqualValues(t, "activated", resp.Status)
}

func TestFindSubscriptionWithEmptyMsidn(t *testing.T) {
	req, rw := requestResponse(http.MethodGet, "/api/subscription/msidn/", nil)
	req = mux.SetURLVars(req, map[string]string{
		"msidn": "",
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
	assert.EqualValues(t, "msidn cannot be empty", resp.Message)
}

func TestFindSubscriptionWithNonExistantMsidn(t *testing.T) {
	req, rw := requestResponse(http.MethodGet, "/api/subscription/msidn/", nil)
	req = mux.SetURLVars(req, map[string]string{
		"msidn": "85245804-a7a8-44f1-bfdf-a05896d81e5b",
	})
	msidnUUID, _ := uuid.Parse("85245804-a7a8-44f1-bfdf-a05896d81e5b")
	mockFindNonExistingSubscription(msidnUUID)
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
	assert.EqualValues(t, "Could not find subscription with msidn 85245804-a7a8-44f1-bfdf-a05896d81e5b, subscription not found", resp.Message)
}

func TestCancelSubscription(t *testing.T) {
	var (
		msidnUUID uuid.UUID
		now       = time.Now().Format("2006-01-02")
	)
	msidnUUID, _ = uuid.Parse("c019ecde-17cb-4ef8-8a7d-85937a9250ed")

	req, rw := requestResponse(http.MethodPatch, "/api/subscription/update-subscription/msidn/{msidn}/status/{status}", nil)
	req = mux.SetURLVars(req, map[string]string{
		"msidn":  "c019ecde-17cb-4ef8-8a7d-85937a9250ed",
		"status": "cancelled",
	})
	mockUpdateSubscription(msidnUUID, now, "cell", "activated")
	mockFindSubscription(msidnUUID, "2021-10-11", "cell", "cancelled")
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
	assert.EqualValues(t, msidnUUID, resp.Msidn)
	assert.EqualValues(t, "2021-10-11", resp.ActivateAt)
	assert.EqualValues(t, "cell", resp.SubType)
	assert.EqualValues(t, "cancelled", resp.Status)
}

func TestPauseSubscriptionWithEmptyMsidn(t *testing.T) {
	req, rw := requestResponse(http.MethodPost, "/api/subscription/update-subscription/msidn/{msidn}/status/{status}", nil)
	req = mux.SetURLVars(req, map[string]string{
		"msidn":  "",
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
	assert.EqualValues(t, "msidn cannot be empty", resp.Message)
}

func TestReactivateSubscriptionWithNonExistingMsidn(t *testing.T) {
	req, rw := requestResponse(http.MethodPost, "/api/subscription", nil)
	req = mux.SetURLVars(req, map[string]string{
		"msidn":  "85245804-a7a8-44f1-bfdf-a05896d81e5b",
		"status": "activated",
	})
	msidnUUID, _ := uuid.Parse("85245804-a7a8-44f1-bfdf-a05896d81e5b")
	mockFindNonExistingSubscription(msidnUUID)
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
	assert.EqualValues(t, "Could not find subscription with msidn 85245804-a7a8-44f1-bfdf-a05896d81e5b, subscription not found", resp.Message)
}

func TestUpdateActivationDate(t *testing.T) {
	var (
		msidnUUID uuid.UUID
		now       = time.Now().Format("2006-01-02")
	)
	msidnUUID, _ = uuid.Parse("c019ecde-17cb-4ef8-8a7d-85937a9250ed")

	req, rw := requestResponse(http.MethodPatch, "/api/subscription/update-activation-date/msidn/{msidn}/date/{date}", nil)
	req = mux.SetURLVars(req, map[string]string{
		"msidn": "c019ecde-17cb-4ef8-8a7d-85937a9250ed",
		"date":  "2021-10-11",
	})
	mockUpdateSubscription(msidnUUID, now, "cell", "pending")
	mockFindSubscription(msidnUUID, "2021-10-11", "cell", "pending")
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
	assert.EqualValues(t, msidnUUID, resp.Msidn)
	assert.EqualValues(t, "2021-10-11", resp.ActivateAt)
	assert.EqualValues(t, "cell", resp.SubType)
	assert.EqualValues(t, "pending", resp.Status)
}

func TestUpdateActivationDateWithWrongDate(t *testing.T) {
	req, rw := requestResponse(http.MethodPatch, "/api/subscription/update-activation-date/msidn/{msidn}/date/{date}", nil)
	req = mux.SetURLVars(req, map[string]string{
		"msidn": "c019ecde-17cb-4ef8-8a7d-85937a9250ed",
		"date":  "2021-09-11",
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
	req, rw := requestResponse(http.MethodPatch, "/api/subscription/update-activation-date/msidn/{msidn}/date/{date}", nil)
	req = mux.SetURLVars(req, map[string]string{
		"msidn": "c019ecde-17cb-4ef8-8a7d-85937a9250ed",
		"date":  "",
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
