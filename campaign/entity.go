package campaign

import (
	"bwastartup/users"
	"time"
)

type Campaign struct {
	ID               int
	UserID           int
	Name             string
	ShortDescription string
	Description      string
	Perks            string
	BackerCount      int
	GoalAmount       int
	CurrentAmount    int
	Slug             string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	CampaignImages   []CampaignImage `gorm:"foreignKey:CampaignID"`
	User             users.User
}

type CampaignImage struct {
	ID         int       `gorm:"primaryKey"`
	CampaignID int       `gorm:"column:campaign_id"`
	FileName   string    `gorm:"column:file_name"`
	IsPrimary  int       `gorm:"column:is_primary"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:update_at"`
}

func (CampaignImage) TableName() string {
	return "campaigns_images"
}
