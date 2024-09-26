package transaction

import "gorm.io/gorm"

type reposutory struct {
	db *gorm.DB
}

type Repository interface {
	GetByCampaignID(campaignID int) ([]Transaction, error)
}

func NewRepository(db *gorm.DB) reposutory {
	return reposutory{db}
}
