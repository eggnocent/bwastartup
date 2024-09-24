package campaign

import (
	"log"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]Campaign, error)
	FindByUserID(userID int) ([]Campaign, error)
	FindByID(userID int) (Campaign, error)
	Save(campaign Campaign) (Campaign, error)
	Update(campaign Campaign) (Campaign, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindAll() ([]Campaign, error) {
	var campaigns []Campaign
	err := r.db.Preload("CampaignImages").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}
	return campaigns, nil
}

func (r *repository) FindByUserID(userID int) ([]Campaign, error) {
	var campaigns []Campaign
	err := r.db.Where("user_id = ?", userID).
		Preload("CampaignImages", "is_primary = ?", 1).Find(&campaigns).Error

	if err != nil {
		return campaigns, err
	}
	return campaigns, nil
}

func (r *repository) FindByID(userID int) (Campaign, error) {
	var campaign Campaign
	err := r.db.Preload("User").Preload("CampaignImages").Where("id = ?", userID).First(&campaign).Error // Gunakan First untuk satu record
	if err != nil {
		return campaign, err
	}
	return campaign, nil
}

func (r *repository) Save(campaign Campaign) (Campaign, error) {
	err := r.db.Create(&campaign).Error
	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (r *repository) Update(campaign Campaign) (Campaign, error) {
	log.Println("Repository: Starting Update process for campaign ID:", campaign.ID)
	err := r.db.Save(&campaign).Error
	if err != nil {
		log.Println("Repository: Error while updating campaign:", err)
		return campaign, err
	}
	log.Println("Repository: Successfully updated campaign ID:", campaign.ID)
	return campaign, nil
}
