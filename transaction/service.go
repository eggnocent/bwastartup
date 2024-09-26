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
