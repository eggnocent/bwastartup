package transaction

import (
	"bwastartup/users"
	"time"
)

type Transaction struct {
	ID         int
	CampaignID int
	UserID     int
	Amount     int
	Status     string
	Code       string
	User       users.User
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
