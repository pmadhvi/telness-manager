package postgress

import (
	"database/sql"
	"github.com/pmadhvi/telness-manager/model"
	log "github.com/sirupsen/logrus"
)

type subscriptionRepo struct {
	Log: log.logger,
	Db: *sql.DB,
}

func NewSubscriptionRepo(log log.logger, db *sql.DB) SubscriptionRepo {
	return &subscriptionRepo{
		Log: log,
		DB: db,
	}
}

func (sr subscriptionRepo)CreateSubscription() error {
	query := `INSERT INTO subscriptions(msidn, activate_at, type, status, created_at)
	VALUES($1, $2, $3, $4, $5);`
	err := sr.Db.Exec(query, msidn, activate_at, type, status, created_at)
	if err != nil {
		log.Errorf("could not insert the data in db: %v", err)
	}
	return nil
}

func (sr subscriptionRepo) FindbyID(msidn string) (model.Subscription, error) {
	query := `SELECT FROM subscriptions(
	WHERE msidn = $1;`
	var sub model.Subscription
	err := sr.Db.Exec(query, msidn).Scan(sub.msidn, sub.activate_at, sub.type, sub.status)
	if err != nil {
		log.Errorf("could not find the subscription with msidn: %v %v", msidn, err)
	}
	return model.Subscription{}, nil
}

// func (sr subscriptionRepo) UpdateSubscription() (model.Subscription, error) {
// 	query := `UPDATE subscriptions(
// 		activate_at = $1,
// 		type = $2, 
// 		status = $3)
// 	WHERE msidn = $4;`
// 	n := sr.Db.Exec(query, activate_at, type, status)
// 	if err != nil {
// 		log.Errorf("could not update the data in db: %v", err)
// 	}
// 	return nil
// 	return model.Subscription{}, nil
// }
