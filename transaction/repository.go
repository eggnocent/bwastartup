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
	var transactions []Transaction

	// Preload untuk Campaign dan CampaignImages dengan log yang lebih detail
	err := r.db.Preload("Campaign.CampaignImages", "is_primary = 1").
		Preload("User").Where("user_id = ?", userID).
		Order("id desc").Find(&transactions).Error

	if err != nil {
		log.Println("Repository: Error fetching transactions for user ID", userID, ":", err)
		return transactions, err
	}

	// Log setelah preload untuk memastikan gambar ter-load
	for _, t := range transactions {
		log.Printf("Repository: Transaction ID: %d, Campaign ID: %d, Number of Images Preloaded: %d",
			t.ID, t.CampaignID, len(t.Campaign.CampaignImages))
	}

	return transactions, nil
}
