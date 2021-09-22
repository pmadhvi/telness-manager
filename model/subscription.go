package model

type SubStatus string

const (
	StatusPending   SubStatus = "pending"
	StatusPaused    SubStatus = "paused"
	StatusActivated SubStatus = "activated"
	StatusCancelled SubStatus = "cancelled"
)

// Subscription represents all data for a phone subscription
type Subscription struct {
	Msisdn     string    `json:"msisdn"`
	ActivateAt string    `json:"activate_at"`
	SubType    string    `json:"sub_type"`
	Status     SubStatus `json:"status"`
	Operator   string    `json:"operator"`
	CreatedAt  string    `json:"created_at"`
	ModifiedAt string    `json:"modified_at"`
}

// CreateSubscription represents all data for a phone subscription create request
type CreateSubscription struct {
	Msisdn     string    `json:"msisdn"`
	ActivateAt string    `json:"activate_at"`
	SubType    string    `json:"sub_type"`
	Status     SubStatus `json:"status"`
}

type ErrorMessage struct {
	Message string `json:"message"`
}

type PtsResponse struct {
	D OperatorDetails `json:"d"`
}

type OperatorDetails struct {
	Type   string `json: "__type"`
	Name   string `json:"Name"`
	Number string `json:"Number"`
}

// var validSubStatusValues = map[SubStatus]struct{}{
// 	StatusPending:   struct{}{},
// 	StatusPaused:    struct{}{},
// 	StatusActivated: struct{}{},
// 	StatusCancelled: struct{}{},
// }

// func (v SubStatus) Valid() bool {
// 	_, ok := validSubStatusValues[v]
// 	return ok
// }
