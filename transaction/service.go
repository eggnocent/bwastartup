package transaction

import (
	"bwastartup/campaign"
	"errors"
)

type service struct {
	repository         Repository
	campaignRepository campaign.Repository
}

type Service interface {
	GetTransactionByCampaignID(input GetCampaignTransactionsInput) ([]Transaction, error)
	GetTransactionByUserID(userID int) ([]Transaction, error)
	CreateTransaction(input CreateTransactionInput) (Transaction, error)
}

func NewService(repository Repository, campaignRepository campaign.Repository) *service {
	return &service{repository, campaignRepository}
}

func (s *service) GetTransactionByCampaignID(input GetCampaignTransactionsInput) ([]Transaction, error) {

	// Fetch the campaign by ID
	campaign, err := s.campaignRepository.FindByID(input.ID)
	if err != nil {
		return []Transaction{}, err
	}

	// Check if the current user is either the campaign owner or has made a transaction in the campaign
	if campaign.UserID != input.User.ID {
		// If not the owner, check if the user has made a transaction in the campaign
		userTransactions, err := s.repository.GetByCampaignID(input.ID)
		if err != nil {
			return []Transaction{}, err
		}
		userIsContributor := false
		for _, transaction := range userTransactions {
			if transaction.UserID == input.User.ID {
				userIsContributor = true
				break
			}
		}
		if !userIsContributor {
			return []Transaction{}, errors.New("Not authorized to view the transactions of this campaign")
		}
	}

	transactions, err := s.repository.GetByCampaignID(input.ID)
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (s *service) GetTransactionByUserID(userID int) ([]Transaction, error) {
	transactions, err := s.repository.GetByUserID(userID)
	if err != nil {
		return transactions, err
	}
	return transactions, nil
}

func (s *service) CreateTransaction(input CreateTransactionInput) (Transaction, error) {
	transaction := Transaction{}
	transaction.CampaignID = input.CampaignID
	transaction.Amount = input.Amount
	transaction.UserID = input.users.ID
	transaction.Status = "pending"

	newTransaction, err := s.repository.Save(transaction)
	if err != nil {
		return transaction, err
	}
	return newTransaction, nil
}
