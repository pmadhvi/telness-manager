package service

import (
	"errors"
	"os"

	"github.com/google/uuid"
	"github.com/pmadhvi/telness-manager/mock"
	"github.com/pmadhvi/telness-manager/model"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"testing"
	"time"
)

var (
	msidn = uuid.New()
	now   = time.Now().Format("2006-01-02")
)

func setupSubscriptionSvc() SubscriptionSvc {
	log := logrus.New()
	log.SetOutput(os.Stdout)

	return SubscriptionSvc{
		Log:              log,
		SubscriptionRepo: &mock.DbMock{},
	}
}
func TestSubscriptionSvc_FindbyID_Success(t *testing.T) {
	s := setupSubscriptionSvc()
	mock.FindByID = func(msidn uuid.UUID) (model.Subscription, error) {
		return model.Subscription{
			Msidn:      msidn,
			ActivateAt: now,
			SubType:    "cell",
			Status:     "pending",
			CreatedAt:  now,
			ModifiedAt: now,
		}, nil
	}
	got, err := s.FindbyID(msidn)

	assert.NotNil(t, got)
	assert.Nil(t, err)
	assert.EqualValues(t, msidn, got.Msidn)
	assert.EqualValues(t, now, got.ActivateAt)
	assert.EqualValues(t, "cell", got.SubType)
	assert.EqualValues(t, "pending", got.Status)
}

func TestSubscriptionSvc_FindbyID_NotFound(t *testing.T) {
	s := setupSubscriptionSvc()
	mock.FindByID = func(msidn uuid.UUID) (model.Subscription, error) {
		return model.Subscription{}, errors.New("subscription not found")
	}
	_, err := s.FindbyID(msidn)
	assert.NotNil(t, err)
}

func TestSubscriptionSvc_Update_Success(t *testing.T) {
	s := setupSubscriptionSvc()
	mock.Update = func(sub model.CreateSubscription) error {
		return nil
	}
	mock.FindByID = func(msidn uuid.UUID) (model.Subscription, error) {
		return model.Subscription{
			Msidn:      msidn,
			ActivateAt: now,
			SubType:    "pbx",
			Status:     "activated",
			CreatedAt:  now,
			ModifiedAt: now,
		}, nil
	}
	request := model.CreateSubscription{
		Msidn:      msidn,
		ActivateAt: now,
		SubType:    "pbx",
		Status:     "activated",
	}
	got, err := s.Update(request)

	assert.NotNil(t, got)
	assert.Nil(t, err)
	assert.EqualValues(t, msidn, got.Msidn)
	assert.EqualValues(t, now, got.ActivateAt)
	assert.EqualValues(t, "pbx", got.SubType)
	assert.EqualValues(t, "activated", got.Status)
}

func TestSubscriptionSvc_Update_Fail(t *testing.T) {
	s := setupSubscriptionSvc()
	mock.Update = func(sub model.CreateSubscription) error {
		return errors.New("cannot update this subscription")
	}
	request := model.CreateSubscription{
		Msidn:      msidn,
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
	mock.FindByID = func(msidn uuid.UUID) (model.Subscription, error) {
		return model.Subscription{
			Msidn:      msidn,
			ActivateAt: now,
			SubType:    "cell",
			Status:     "pending",
			CreatedAt:  now,
			ModifiedAt: now,
		}, nil
	}
	request := model.CreateSubscription{
		Msidn:      msidn,
		ActivateAt: now,
		SubType:    "cell",
		Status:     "pending",
	}
	got, err := s.Create(request)

	assert.NotNil(t, got)
	assert.Nil(t, err)
	assert.EqualValues(t, msidn, got.Msidn)
	assert.EqualValues(t, now, got.ActivateAt)
	assert.EqualValues(t, "cell", got.SubType)
	assert.EqualValues(t, "pending", got.Status)
}

func TestSubscriptionSvc_Create_Fail(t *testing.T) {
	s := setupSubscriptionSvc()
	mock.Create = func(sub model.CreateSubscription) error {
		return errors.New("cannot create this subscription")
	}
	request := model.CreateSubscription{
		Msidn:      msidn,
		ActivateAt: now,
		SubType:    "pbx",
		Status:     "activated",
	}
	_, err := s.Create(request)

	assert.NotNil(t, err)
}
