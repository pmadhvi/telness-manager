package service

import (
	"github.com/pmadhvi/telness-manager/model"
	log "github.com/sirupsen/logrus"
)

type SubscriptionRepoInterface interface {
	CreateSubscription(sub model.CreateSubscription) error
	FindSubscriptionbyID(id string) (model.Subscription, error)
	UpdateSubscription(sub model.CreateSubscription) error
}

type PtsClientInterface interface {
	GetOperatorDetails(msisdn string) (model.PtsResponse, error)
}

type SubscriptionSvc struct {
	Log              *log.Logger
	SubscriptionRepo SubscriptionRepoInterface
	PtsClient        PtsClientInterface
}

func (s SubscriptionSvc) Create(subreq model.CreateSubscription) (model.Subscription, error) {
	err := s.SubscriptionRepo.CreateSubscription(subreq)
	if err != nil {
		s.Log.Errorf("Could not create subscription due to error: %v", err)
		return model.Subscription{}, err
	}
	sub, err := s.FindbyID(subreq.Msisdn)
	if err != nil {
		s.Log.Errorf("Could not find created subscription due to error: %v", err)
		return model.Subscription{}, err
	}
	return sub, nil
}

func (s SubscriptionSvc) FindbyID(msisdn string) (model.Subscription, error) {
	var sub model.Subscription
	sub, err := s.SubscriptionRepo.FindSubscriptionbyID(msisdn)
	if err != nil {
		s.Log.Errorf("Could not find subscription by id %v due to error: %v", msisdn, err)
		return model.Subscription{}, err
	}
	var ptsResponse model.PtsResponse
	ptsResponse, err = s.PtsClient.GetOperatorDetails(msisdn)
	if err != nil {
		s.Log.Errorf("Could not find operator details for subscription with msisdn %v due to error: %v", msisdn, err)
		return sub, err
	}
	sub.Operator = ptsResponse.D.Name
	return sub, nil
}

func (s SubscriptionSvc) Update(subreq model.CreateSubscription) (model.Subscription, error) {
	err := s.SubscriptionRepo.UpdateSubscription(subreq)
	if err != nil {
		s.Log.Errorf("Could not update subscription due to error: %v", err)
		return model.Subscription{}, err
	}
	sub, err := s.FindbyID(subreq.Msisdn)
	if err != nil {
		s.Log.Errorf("Could not find updated subscription due to error: %v", err)
		return model.Subscription{}, err
	}
	return sub, nil
}
