package transaction

import (
	"bwastartup/users"
)

type GetCampaignTransactionsInput struct {
	ID   int `uri:"id" binding:"required"`
	User users.User
}

type CreateTransactionInput struct {
	Amount     int `json:"amount" binding:"required"`
	CampaignID int `json:"campaign_id" binding:"required"`
	users      users.User
}
