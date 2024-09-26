package transaction

type service struct {
	repository Repository
}

type Service interface {
	GetTransactionByCampaignID(campaignID int) ([]Transaction, error)
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetTransactionByCampaignID(campaignID int) ([]Transaction, error) {
	transactions, err := s.repository.GetByCampaignID(campaignID)
	if err != nil {
		return transactions, err
	}
	return transactions, nil
}
