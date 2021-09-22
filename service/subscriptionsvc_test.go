package service

import (
	"errors"
	"os"

	"github.com/pmadhvi/telness-manager/mock"
	"github.com/pmadhvi/telness-manager/model"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"testing"
	"time"
)

var (
	msisdn = "+46107500500"
	now    = time.Now().Format("2006-01-02")
)

func setupSubscriptionSvc() SubscriptionSvc {
	log := logrus.New()
	log.SetOutput(os.Stdout)

	return SubscriptionSvc{
		Log:              log,
		SubscriptionRepo: &mock.DbMock{},
		PtsClient:        &mock.ClientMock{},
	}
}
func TestSubscriptionSvc_FindbyID_Success(t *testing.T) {
	s := setupSubscriptionSvc()
	mock.FindByID = func(msisdn string) (model.Subscription, error) {
		return model.Subscription{
			Msisdn:     msisdn,
			ActivateAt: now,
			SubType:    "cell",
			Status:     "pending",
			CreatedAt:  now,
			ModifiedAt: now,
		}, nil
	}
	mock.GetOperator = func(msisdn string) (model.PtsResponse, error) {
		return model.PtsResponse{
			D: model.OperatorDetails{Name: "Telness AB"},
		}, nil
	}
	got, err := s.FindbyID(msisdn)
	assert.NotNil(t, got)
	assert.Nil(t, err)
	assert.EqualValues(t, msisdn, got.Msisdn)
	assert.EqualValues(t, now, got.ActivateAt)
	assert.EqualValues(t, "cell", got.SubType)
	assert.EqualValues(t, "pending", got.Status)
	assert.EqualValues(t, "Telness AB", got.Operator)
}

func TestSubscriptionSvc_FindbyID_NotFound(t *testing.T) {
	s := setupSubscriptionSvc()
	mock.FindByID = func(msisdn string) (model.Subscription, error) {
		return model.Subscription{}, errors.New("subscription not found")
	}
	_, err := s.FindbyID(msisdn)
	assert.NotNil(t, err)
}

func TestSubscriptionSvc_Update_Success(t *testing.T) {
	s := setupSubscriptionSvc()
	mock.Update = func(sub model.CreateSubscription) error {
		return nil
	}
	mock.FindByID = func(msisdn string) (model.Subscription, error) {
		return model.Subscription{
			Msisdn:     msisdn,
			ActivateAt: now,
			SubType:    "pbx",
			Status:     "activated",
			Operator:   "Telness AB",
			CreatedAt:  now,
			ModifiedAt: now,
		}, nil
	}
	mock.GetOperator = func(msisdn string) (model.PtsResponse, error) {
		return model.PtsResponse{
			D: model.OperatorDetails{Name: "Telness AB"},
		}, nil
	}
	request := model.CreateSubscription{
		Msisdn:     msisdn,
		ActivateAt: now,
		SubType:    "pbx",
		Status:     "activated",
	}
	got, err := s.Update(request)

	assert.NotNil(t, got)
	assert.Nil(t, err)
	assert.EqualValues(t, msisdn, got.Msisdn)
	assert.EqualValues(t, now, got.ActivateAt)
	assert.EqualValues(t, "pbx", got.SubType)
	assert.EqualValues(t, "activated", got.Status)
	assert.EqualValues(t, "Telness AB", got.Operator)
}

func TestSubscriptionSvc_Update_Fail(t *testing.T) {
	s := setupSubscriptionSvc()
	mock.Update = func(sub model.CreateSubscription) error {
		return errors.New("cannot update this subscription")
	}
	request := model.CreateSubscription{
		Msisdn:     msisdn,
		ActivateAt: now,
		SubType:    "cell",
		Status:     "activated",
	}
	_, err := s.Update(request)

	assert.NotNil(t, err)
}

func TestSubscriptionSvc_Create_Success(t *testing.T) {
	s := setupSubscriptionSvc()
	mock.Create = func(sub model.CreateSubscription) error {
		return nil
	}
	mock.FindByID = func(msisdn string) (model.Subscription, error) {
		return model.Subscription{
			Msisdn:     msisdn,
			ActivateAt: now,
			SubType:    "cell",
			Status:     "pending",
			Operator:   "Telness AB",
			CreatedAt:  now,
			ModifiedAt: now,
		}, nil
	}
	mock.GetOperator = func(msisdn string) (model.PtsResponse, error) {
		return model.PtsResponse{
			D: model.OperatorDetails{Name: "Telness AB"},
		}, nil
	}
	request := model.CreateSubscription{
		Msisdn:     msisdn,
		ActivateAt: now,
		SubType:    "cell",
		Status:     "pending",
	}
	got, err := s.Create(request)

	assert.NotNil(t, got)
	assert.Nil(t, err)
	assert.EqualValues(t, msisdn, got.Msisdn)
	assert.EqualValues(t, now, got.ActivateAt)
	assert.EqualValues(t, "cell", got.SubType)
	assert.EqualValues(t, "pending", got.Status)
	assert.EqualValues(t, "Telness AB", got.Operator)
}

func TestSubscriptionSvc_Create_Fail(t *testing.T) {
	s := setupSubscriptionSvc()
	mock.Create = func(sub model.CreateSubscription) error {
		return errors.New("cannot create this subscription")
	}
	request := model.CreateSubscription{
		Msisdn:     msisdn,
		ActivateAt: now,
		SubType:    "pbx",
		Status:     "activated",
	}
	_, err := s.Create(request)

	assert.NotNil(t, err)
}
