package payment

import (
	"bwastartup/users"
	"strconv"

	"github.com/veritrans/go-midtrans"
)

type service struct {
}

type Service interface {
	GetPaymentURL(transaction Transaction, user users.User) (string, error)
}

func Newservice() *service {
	return &service{}
}

func (s *service) GetPaymentURL(transaction Transaction, user users.User) (string, error) {
	midclient := midtrans.NewClient()
	midclient.ServerKey = "SB-Mid-server-dnT0rb4AQ4G1Yrss64G6OBpB"
	midclient.ClientKey = "SB-Mid-client-GRodqP_bd2veZmIa"
	midclient.APIEnvType = midtrans.Sandbox

	snapGateway := midtrans.SnapGateway{
		Client: midclient,
	}

	snapReq := &midtrans.SnapReq{
		CustomerDetail: &midtrans.CustDetail{
			Email: user.Email,
			FName: user.Name,
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(transaction.ID),
			GrossAmt: int64(transaction.Amount),
		},
	}

	snapTokenResp, err := snapGateway.GetToken(snapReq)

	if err != nil {
		return "", err
	}
	return snapTokenResp.RedirectURL, nil
}
