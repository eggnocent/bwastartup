package campaign

import (
	"fmt"
	"log"

	"github.com/gosimple/slug"
)

type Service interface {
	GetCampaigns(userID int) ([]Campaign, error)
	GetCampaignsByID(input GetCampaignsDetailInput) (Campaign, error)
	CreateCampaign(input CreateCampaignInput) (Campaign, error)
	UpdateCampaign(inputID GetCampaignsDetailInput, inputData CreateCampaignInput) (Campaign, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetCampaigns(userID int) ([]Campaign, error) {
	if userID != 0 {
		campaign, err := s.repository.FindByUserID(userID)
		if err != nil {
			return campaign, err
		}
		return campaign, nil
	}
	campaigns, err := s.repository.FindAll()
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (s *service) GetCampaignsByID(input GetCampaignsDetailInput) (Campaign, error) {
	campaign, err := s.repository.FindByID(input.ID)
	if err != nil {
		return campaign, err
	}
	return campaign, nil
}

func (s *service) CreateCampaign(input CreateCampaignInput) (Campaign, error) {
	campaign := Campaign{}
	campaign.Name = input.Name
	campaign.ShortDescription = input.ShortDescription
	campaign.Description = input.Description
	campaign.Perks = input.Perks
	campaign.GoalAmount = input.GoalAmount
	campaign.UserID = input.User.ID

	slugCandidate := fmt.Sprintf("%s %d", input.Name, input.User.ID)
	campaign.Slug = slug.Make(slugCandidate)

	newCampaign, err := s.repository.Save(campaign)
	if err != nil {
		return newCampaign, err
	}
	return newCampaign, nil
}

func (s *service) UpdateCampaign(inputID GetCampaignsDetailInput, inputData CreateCampaignInput) (Campaign, error) {
	log.Println("Service: Fetching campaign by ID:", inputID.ID)
	campaign, err := s.repository.FindByID(inputID.ID)
	if err != nil {
		log.Println("Service: Error fetching campaign by ID:", err)
		return campaign, err
	}

	log.Println("Service: Updating campaign fields")
	campaign.Name = inputData.Name
	campaign.ShortDescription = inputData.ShortDescription
	campaign.Description = inputData.Description
	campaign.Perks = inputData.Perks
	campaign.GoalAmount = inputData.GoalAmount

	log.Println("Service: Saving updated campaign")
	updateCampaign, err := s.repository.Update(campaign)
	if err != nil {
		log.Println("Service: Error while saving updated campaign:", err)
		return updateCampaign, err
	}
	log.Println("Service: Successfully updated campaign ID:", updateCampaign.ID)
	return updateCampaign, nil
}
