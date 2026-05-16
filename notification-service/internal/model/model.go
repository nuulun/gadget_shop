package model

import "time"

type Notification struct {
	ID        uint64    `json:"id"`
	Recipient string    `json:"recipient"`
	Type      string    `json:"type"`
	Message   string    `json:"message"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type SendNotificationInput struct {
	Recipient string `json:"recipient"`
	Type      string `json:"type"`
	Message   string `json:"message"`
}
