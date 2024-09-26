package transaction

import "bwastartup/users"

type GetCampaignTransactionsInput struct {
	ID   int `uri:"id" binding:"required"`
	User users.User
}
