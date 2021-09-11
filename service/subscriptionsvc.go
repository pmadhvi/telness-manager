package service

import (
	"github.com/google/uuid"
	"github.com/pmadhvi/telness-manager/model"
)

type subscriptionRepo interface {
	CreateSubscription(sub *model.Subscription) error
	FindByIdSubscription(id uuid.UUID) (model.Subscription, error)
	UpdateSubscription(sub *model.Subscription) (model.Subscription, error)
}

type SubscriptionSvc struct {
	SubscriptionRepo subscriptionRepo
}

func (s SubscriptionSvc) CreateSubscription(sub *model.Subscription) error {
	err := s.SubscriptionRepo.CreateSubscription(sub)
	if err != nil {
		return err
	}
	return nil
}

func (s SubscriptionSvc) FindbyID(id uuid.UUID) (model.Subscription, error) {
	var sub model.Subscription
	sub, err := s.SubscriptionRepo.FindByIdSubscription(id)
	if err != nil {
		return model.Subscription{}, err
	}
	return sub, nil
}

func (s SubscriptionSvc) UpdateSubscription() (model.Subscription, error) {
	var sub model.Subscription
	sub, err := s.SubscriptionRepo.UpdateSubscription(sub)
	if err != nil {
		return model.Subscription{}, err
	}
	return sub, nil
}
