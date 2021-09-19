package model

import (
	"github.com/google/uuid"
)

type SubStatus string

const (
	StatusPending   SubStatus = "pending"
	StatusPaused    SubStatus = "paused"
	StatusActivated SubStatus = "activated"
	StatusCancelled SubStatus = "cancelled"
)

// Subscription represents all data for a phone subscription
type Subscription struct {
	Msidn      uuid.UUID `json:"msidn"`
	ActivateAt string    `json:"activate_at"`
	SubType    string    `json:"sub_type"`
	Status     SubStatus `json:"status"`
	CreatedAt  string    `json:"created_at"`
	ModifiedAt string    `json:"modified_at"`
}

// CreateSubscription represents all data for a phone subscription create request
type CreateSubscription struct {
	Msidn      uuid.UUID `json:"msidn"`
	ActivateAt string    `json:"activate_at"`
	SubType    string    `json:"sub_type"`
	Status     SubStatus `json:"status"`
}

type ErrorMessage struct {
	Message string `json:"message"`
}
