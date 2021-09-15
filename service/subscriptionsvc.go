package service

import (
	"github.com/google/uuid"
	"github.com/pmadhvi/telness-manager/model"
	log "github.com/sirupsen/logrus"
)

type subscriptionRepo interface {
	CreateSubscription(sub model.CreateSubscription) error
	FindSubscriptionbyID(id uuid.UUID) (model.Subscription, error)
	UpdateSubscription(sub model.CreateSubscription) error
}

type SubscriptionSvc struct {
	Log              *log.Logger
	SubscriptionRepo subscriptionRepo
}

func (s SubscriptionSvc) Create(subreq model.CreateSubscription) (model.Subscription, error) {
	err := s.SubscriptionRepo.CreateSubscription(subreq)
	if err != nil {
		s.Log.Errorf("Could not create subscription due to error: %v", err)
		return model.Subscription{}, err
	}
	sub, err := s.SubscriptionRepo.FindSubscriptionbyID(subreq.Msidn)
	if err != nil {
		s.Log.Errorf("Could not find created subscription due to error: %v", err)
		return model.Subscription{}, err
	}
	return sub, nil
}

func (s SubscriptionSvc) FindbyID(id uuid.UUID) (model.Subscription, error) {
	var sub model.Subscription
	sub, err := s.SubscriptionRepo.FindSubscriptionbyID(id)
	if err != nil {
		s.Log.Errorf("Could not find subscription by id %v due to error: %v", id, err)
		return model.Subscription{}, err
	}
	return sub, nil
}

func (s SubscriptionSvc) Update(subreq model.CreateSubscription) (model.Subscription, error) {
	err := s.SubscriptionRepo.UpdateSubscription(subreq)
	if err != nil {
		s.Log.Errorf("Could not update subscription due to error: %v", err)
		return model.Subscription{}, err
	}
	sub, err := s.SubscriptionRepo.FindSubscriptionbyID(subreq.Msidn)
	if err != nil {
		s.Log.Errorf("Could not find updated subscription due to error: %v", err)
		return model.Subscription{}, err
	}
	return sub, nil
}
