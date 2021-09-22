package mock

import (
	"github.com/pmadhvi/telness-manager/model"
)

var (
	FindByID    func(msisdn string) (model.Subscription, error)
	Create      func(sub model.CreateSubscription) error
	Update      func(sub model.CreateSubscription) error
	GetOperator func(msisdn string) (model.PtsResponse, error)
)

type DbMock struct{}

func (m DbMock) CreateSubscription(sub model.CreateSubscription) error {
	return Create(sub)
}
func (m DbMock) FindSubscriptionbyID(msisdn string) (model.Subscription, error) {
	return FindByID(msisdn)
}
func (m DbMock) UpdateSubscription(sub model.CreateSubscription) error {
	return Update(sub)
}

type ClientMock struct{}

func (c ClientMock) GetOperatorDetails(msisdn string) (model.PtsResponse, error) {
	return GetOperator(msisdn)
}
