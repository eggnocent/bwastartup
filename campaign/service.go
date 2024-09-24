package campaign

type Service interface {
	GetCampaigns(userID int) ([]Campaign, error)
	GetCampaignsByID(input GetCampaignsDetailInput) (Campaign, error)
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
