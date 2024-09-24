package campaign

type GetCampaignsDetailInput struct {
	ID int `uri:"id" binding:"required"`
}
