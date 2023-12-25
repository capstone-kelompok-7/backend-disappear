package payment

import (
	"errors"
	"fmt"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

func CreateCoreAPIPaymentRequest(coreClient coreapi.Client, orderID string, totalAmountPaid int64, paymentType coreapi.CoreapiPaymentType, name string, email string) (*coreapi.ChargeResponse, error) {
	var paymentRequest *coreapi.ChargeReq

	switch paymentType {
	case coreapi.PaymentTypeQris, coreapi.PaymentTypeBankTransfer, coreapi.PaymentTypeGopay, coreapi.PaymentTypeShopeepay:
		paymentRequest = &coreapi.ChargeReq{
			PaymentType: paymentType,
			TransactionDetails: midtrans.TransactionDetails{
				OrderID:  orderID,
				GrossAmt: totalAmountPaid,
			},
		}
	default:
		return nil, errors.New("Jenis pembayaran tidak valid")
	}

	paymentRequest.CustomerDetails = &midtrans.CustomerDetails{
		FName: name,
		Email: email,
	}

	resp, err := coreClient.ChargeTransaction(paymentRequest)
	if err != nil {
		fmt.Println("Error creating payment request:", err.GetMessage())
		return nil, err
	}

	fmt.Println("Menyimpan data pembayaran: OrderID=", orderID, " Name=", name, " Email=", email)
	fmt.Println("Payment request created successfully:", resp)

	return resp, nil
}
