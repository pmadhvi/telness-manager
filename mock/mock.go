package mock

import (
	"github.com/google/uuid"
	"github.com/pmadhvi/telness-manager/model"
)

var (
	FindByID func(msidn uuid.UUID) (model.Subscription, error)
	Create   func(sub model.CreateSubscription) error
	Update   func(sub model.CreateSubscription) error
)

type DbMock struct{}

func (m DbMock) CreateSubscription(sub model.CreateSubscription) error {
	return Create(sub)
}
func (m DbMock) FindSubscriptionbyID(id uuid.UUID) (model.Subscription, error) {
	return FindByID(id)
}
func (m DbMock) UpdateSubscription(sub model.CreateSubscription) error {
	return Update(sub)
}
