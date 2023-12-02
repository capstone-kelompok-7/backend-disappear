package payment

import (
	"errors"
	"fmt"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

func CreateCoreAPIPaymentRequest(coreClient coreapi.Client, orderID string, totalAmountPaid int64, paymentType coreapi.CoreapiPaymentType) (*coreapi.ChargeResponse, error) {
	var paymentRequest *coreapi.ChargeReq

	switch paymentType {
	case coreapi.PaymentTypeQris:
		paymentRequest = &coreapi.ChargeReq{
			PaymentType: paymentType,
			TransactionDetails: midtrans.TransactionDetails{
				OrderID:  orderID,
				GrossAmt: totalAmountPaid,
			},
		}
	case coreapi.PaymentTypeBankTransfer:
		paymentRequest = &coreapi.ChargeReq{
			PaymentType: paymentType,
			TransactionDetails: midtrans.TransactionDetails{
				OrderID:  orderID,
				GrossAmt: totalAmountPaid,
			},
		}
	case coreapi.PaymentTypeGopay:
		paymentRequest = &coreapi.ChargeReq{
			PaymentType: paymentType,
			TransactionDetails: midtrans.TransactionDetails{
				OrderID:  orderID,
				GrossAmt: totalAmountPaid,
			},
		}
	case coreapi.PaymentTypeShopeepay:
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

	resp, err := coreClient.ChargeTransaction(paymentRequest)
	if err != nil {
		fmt.Println("Error creating payment request:", err.GetMessage())
		return nil, err
	}

	fmt.Println("Payment request created successfully:", resp)
	return resp, nil
}
