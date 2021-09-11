// +build integration

package integrationtest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/pmadhvi/iban-validator/handlers"
)

type response struct {
	Message string
}

func TestValidIban(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/iban/validate/", nil)
	req = mux.SetURLVars(req, map[string]string{
		"iban": "SE4550000000058398257466",
	})
	rw := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.ValidateIbanHandler)
	handler.ServeHTTP(rw, req)
	if status := rw.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := "Iban is valid."
	var resp response
	err := json.NewDecoder(rw.Body).Decode(&resp)
	if err != nil {
		t.Errorf("could not decode response: %v", err)
	}
	if resp.Message != expected {
		t.Errorf("handler returned unexpected body: got: %v, want: %v",
			resp.Message, expected)
	}
}

func TestInvalidIban(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/iban/validate/", nil)
	req = mux.SetURLVars(req, map[string]string{
		"iban": "BR9700360305000010009795493",
	})
	rw := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.ValidateIbanHandler)
	handler.ServeHTTP(rw, req)
	if status := rw.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := "Invalid IBAN, IBAN format for country does not match."
	var resp response
	err := json.NewDecoder(rw.Body).Decode(&resp)
	if err != nil {
		t.Errorf("could not decode response: %v", err)
	}
	if resp.Message != expected {
		t.Errorf("handler returned unexpected body: got: %v, want: %v",
			resp.Message, expected)
	}
}

func TestEmptyIban(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/iban/validate/", nil)
	req = mux.SetURLVars(req, map[string]string{
		"iban": "",
	})
	rw := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.ValidateIbanHandler)
	handler.ServeHTTP(rw, req)
	if status := rw.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := "Invalid IBAN, iban cannot be empty."
	var resp response
	err := json.NewDecoder(rw.Body).Decode(&resp)
	if err != nil {
		t.Errorf("could not decode response: %v", err)
	}
	if resp.Message != expected {
		t.Errorf("handler returned unexpected body: got: %v, want: %v",
			resp.Message, expected)
	}
}

func TestIbanWithMinimumLength(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/iban/validate/", nil)
	req = mux.SetURLVars(req, map[string]string{
		"iban": "HU40",
	})
	rw := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.ValidateIbanHandler)
	handler.ServeHTTP(rw, req)
	if status := rw.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := "Invalid IBAN, minimum length for IBAN should be 5."
	var resp response
	err := json.NewDecoder(rw.Body).Decode(&resp)
	if err != nil {
		t.Errorf("could not decode response: %v", err)
	}
	if resp.Message != expected {
		t.Errorf("handler returned unexpected body: got: %v, want: %v",
			resp.Message, expected)
	}
}

func TestIbanWithWrongChecksum(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/iban/validate/", nil)
	req = mux.SetURLVars(req, map[string]string{
		"iban": "HU40117730161111101800000000",
	})
	rw := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.ValidateIbanHandler)
	handler.ServeHTTP(rw, req)
	if status := rw.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := "Invalid IBAN, wrong checksum for IBAN."
	var resp response
	err := json.NewDecoder(rw.Body).Decode(&resp)
	if err != nil {
		t.Errorf("could not decode response: %v", err)
	}
	if resp.Message != expected {
		t.Errorf("handler returned unexpected body: got: %v, want: %v",
			resp.Message, expected)
	}
}

func TestIbanHealth(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/iban/validate/health", nil)
	rw := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.CheckHealthHandler)
	handler.ServeHTTP(rw, req)
	if status := rw.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := "Iban Validator application is alive."
	var resp response
	err := json.NewDecoder(rw.Body).Decode(&resp)
	if err != nil {
		t.Errorf("could not decode response: %v", err)
	}
	if resp.Message != expected {
		t.Errorf("handler returned unexpected body: got: %v, want: %v",
			resp.Message, expected)
	}
}
