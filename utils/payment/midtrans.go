package payment

import (
	"github.com/capstone-kelompok-7/backend-disappear/config"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

func InitSnapMidtrans(config config.Config) coreapi.Client {
	var coreClient coreapi.Client
	coreClient.New(config.ServerKey, midtrans.Sandbox)
	return coreClient
}
