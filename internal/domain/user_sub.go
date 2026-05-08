package domain

import (
	"time"

	"github.com/google/uuid"
)

type Subscriptions struct {
	ID             uuid.UUID `json:"id"`
	UserID         uuid.UUID `json:"user_id"`
	ServiceName    string    `json:"service_name"`
	Price          int       `json:"price"`
	DateCreated    time.Time `json:"start_date"`
	DateConclusion time.Time `json:"conclusion_date"`
}
