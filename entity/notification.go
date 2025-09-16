package model

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID             uuid.UUID          `json:"id" db:"id" binding:"required"`
	CreatedAt      time.Time          `json:"created_at" db:"created_at" binding:"required"`
	ModifiedAt     time.Time          `json:"modified_at" db:"modified_at"`
	ExpirationDate time.Time          `json:"expiration_date" db:"expiration_date" binding:"required"`
	Message        string             `json:"message" db:"message" binding:"required"`
	Error          string             `json:"error" db:"error"`
	UserUID        string             `json:"user_uid" db:"user_uid"`
	MessageType    string             `json:"message_type" db:"message_type"`
	Link           string             `json:"link" db:"link"`
	Status         NotificationStatus `json:"status" db:"status"`
	Subject        string             `json:"subject" db:"subject"`
	CreatedBy      string             `json:"created_by" db:"created_by"`
}
