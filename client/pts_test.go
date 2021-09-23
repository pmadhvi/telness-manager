package client

import (
	"testing"

	"github.com/pmadhvi/telness-manager/mock"
	"github.com/pmadhvi/telness-manager/model"
	"github.com/stretchr/testify/assert"
)

var client = &mock.ClientMock{}

func GetOperatorSuccess() {
	mock.GetOperator = func(msisdn string) (model.PtsResponse, error) {
		return model.PtsResponse{
			D: model.OperatorDetails{Name: "Telness AB"},
		}, nil
	}
}

func GetOperatorFail() {
	mock.GetOperator = func(msisdn string) (model.PtsResponse, error) {
		return model.PtsResponse{
			D: model.OperatorDetails{Name: "Operatör saknas"},
		}, nil
	}
}

func TestGetOperatorDetailsSuccessWithoutCountryCode(t *testing.T) {
	msisdn := "0107500500"
	GetOperatorSuccess()
	resp, err := client.GetOperatorDetails(msisdn)
	assert.NotNil(t, resp)
	assert.Nil(t, err)
	assert.EqualValues(t, "Telness AB", resp.D.Name)
}

func TestGetOperatorDetailsSuccessWithCountryCode(t *testing.T) {
	msisdn := "+46107500500"
	GetOperatorSuccess()
	resp, err := client.GetOperatorDetails(msisdn)
	assert.NotNil(t, resp)
	assert.Nil(t, err)
	assert.EqualValues(t, "Telness AB", resp.D.Name)
}

func TestGetOperatorDetailsWithWrongFormat(t *testing.T) {
	msisdn := "0107500500000"
	GetOperatorFail()
	resp, err := client.GetOperatorDetails(msisdn)
	assert.NotNil(t, resp)
	assert.Nil(t, err)
	assert.EqualValues(t, "Operatör saknas", resp.D.Name)
}
