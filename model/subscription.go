package model

import "github.com/google/uuid"

type SubStatus string

const (
	StatusPending   SubStatus = "pending"
	StatusPaused    SubStatus = "paused"
	StatusActivated SubStatus = "activated"
	StatusCancelled SubStatus = "cancelled"
)

// Subscription represnts all data for a phone subscription
type Subscription struct {
	Msidn      uuid.UUID `json:"msidn"`
	ActivateAt string    `json:"activate_at"`
	Type       string    `json:"type"`
	Status     SubStatus `json:"status"`
}
