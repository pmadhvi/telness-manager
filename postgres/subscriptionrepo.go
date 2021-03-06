package postgres

import (
	"database/sql"
	"time"

	"github.com/pmadhvi/telness-manager/model"
	log "github.com/sirupsen/logrus"
)

type subscriptionRepo struct {
	db  *sql.DB
	log *log.Logger
}

func NewSubscriptionRepo(db *sql.DB, log *log.Logger) *subscriptionRepo {
	return &subscriptionRepo{
		db:  db,
		log: log,
	}
}

func (sr subscriptionRepo) CreateSubscription(sub model.CreateSubscription) error {
	query := `INSERT INTO subscription(msisdn, activate_at, sub_type, status, created_at, modified_at)
	VALUES($1, $2, $3, $4, $5, $6)`
	_, err := sr.db.Exec(query, sub.Msisdn, sub.ActivateAt, sub.SubType, sub.Status, time.Now(), time.Now())
	if err != nil {
		sr.log.Errorf("could not insert the data in db: %v", err)
		return err
	}
	return nil
}

func (sr subscriptionRepo) FindSubscriptionbyID(msisdn string) (model.Subscription, error) {
	query := `SELECT * FROM subscription
	WHERE msisdn = $1`
	var sub model.Subscription
	row := sr.db.QueryRow(query, msisdn)
	err := row.Scan(&sub.Msisdn, &sub.ActivateAt, &sub.SubType, &sub.Status, &sub.CreatedAt, &sub.ModifiedAt)
	if err != nil || err == sql.ErrNoRows {
		sr.log.Errorf("No rows were returned! %v", err)
		return model.Subscription{}, err
	}
	return sub, nil
}

func (sr subscriptionRepo) UpdateSubscription(sub model.CreateSubscription) error {
	query := `UPDATE subscription
		SET 
		(activate_at, sub_type, status, modified_at) = ($1, $2, $3, $4)
		WHERE msisdn = $5`
	_, err := sr.db.Exec(query, sub.ActivateAt, sub.SubType, sub.Status, time.Now(), sub.Msisdn)
	if err != nil {
		sr.log.Errorf("could not update the data in db: %v", err)
		return err
	}

	return nil
}
