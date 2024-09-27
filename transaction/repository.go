package transaction

import (
	"log"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

type Repository interface {
	GetByCampaignID(campaignID int) ([]Transaction, error)
	GetByUserID(userID int) ([]Transaction, error)
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetByCampaignID(campaignID int) ([]Transaction, error) {
	var transaction []Transaction
	err := r.db.Preload("User").Where("campaign_id = ?", campaignID).Order("id desc").Find(&transaction).Error
	if err != nil {
		return transaction, err
	}
	return transaction, nil
}

func (r *repository) GetByUserID(userID int) ([]Transaction, error) {
	var transaction []Transaction
	err := r.db.Preload("Campaign").Preload("User").Where("user_id = ?", userID).Order("id desc").Find(&transaction).Error
	if err != nil {
		log.Println("Repository: Error fetching transactions for user ID", userID, ":", err)
		return transaction, err
	}
	return transaction, nil
}
