package utilities

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	request "schoolPayment/dtos/request"
	"schoolPayment/dtos/response"
	models "schoolPayment/models"
	"schoolPayment/repositories"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

var (
	s             snap.Client
	c             coreapi.Client
	errorMidtrans *midtrans.Error
)

// Map a single payment method to Midtrans SnapPaymentType
func mapToSnapPaymentType(paymentMethod *models.PaymentMethod) snap.SnapPaymentType {
	switch paymentMethod.PaymentMethod {
	case "VA":
		switch paymentMethod.BankCode {
		case "002":
			return snap.SnapPaymentType(snap.PaymentTypeBRIVA)
		case "008":
			return snap.SnapPaymentType(snap.PaymentTypeEChannel)
		case "009":
			return snap.SnapPaymentType(snap.PaymentTypeBNIVA)
		case "014":
			return snap.SnapPaymentType(snap.PaymentTypeBCAVA)
		case "013":
			return snap.SnapPaymentType(snap.PaymentTypePermataVA)
		case "022":
			return snap.SnapPaymentType("cimb_va")
		}
	case "CC":
		return snap.SnapPaymentType(snap.PaymentTypeCreditCard)
	case "QR":
		return snap.SnapPaymentType("other_qris")
	}
	return "" // Return an empty string if no match is found
}

// calculateAdminFee calculates the admin fee based on the payment method type.
func CalculateAdminFee(billingAmount int64, paymentMethod *models.PaymentMethod) (int64, error) {
	var adminFee int64
	switch paymentMethod.PaymentMethod {
	case "VA":
		adminFee = int64(paymentMethod.AdminFee)
	case "CC", "QR":
		// Parse admin fee percentage
		adminFeePercentage, err := strconv.ParseFloat(paymentMethod.AdminFeePercentage, 64)
		if err != nil {
			return 0, fmt.Errorf("failed to parse AdminFeePercentage: %v", err)
		}
		adminFee = int64(float64(billingAmount) * adminFeePercentage / 100)

		// Add flat fee for credit card payments
		if paymentMethod.PaymentMethod == "CC" {
			adminFee += int64(paymentMethod.AdminFee)
		}
	default:
		return 0, fmt.Errorf("unsupported payment method for admin fee calculation")
	}
	return adminFee, nil
}

func SendRequestPayment(orderID string, studendID, billingAmount, paymentMethodId int, billingStudentIds []string, invoiceNumber string, listAccountNumber []string, listBillingId []int, bankName string, userId int) (*snap.Response, error) {
	s.New(os.Getenv("SERVER_KEY"), midtrans.Sandbox)

	paymentMethod, err := repositories.GetPaymentMethodByID(paymentMethodId)
	if err != nil {
		return nil, err
	}

	// Map the payment method to SnapPaymentType
	enabledPayment := mapToSnapPaymentType(paymentMethod)
	if enabledPayment == "" {
		return nil, fmt.Errorf("unsupported payment method or bank code")
	}

	// Calculate the admin fee using the helper function
	adminFee, err := CalculateAdminFee(int64(billingAmount), paymentMethod)
	if err != nil {
		return nil, err
	}

	totalAmount := int64(billingAmount) + adminFee

	// Request body
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: totalAmount,
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		EnabledPayments: []snap.SnapPaymentType{enabledPayment},
	}

	resp, errorMidtrans := s.CreateTransaction(req)
	if errorMidtrans != nil {
		apiError, err := ExtractErrorMessage(errorMidtrans.Message)
		if err != nil {
			return nil, err
		}
		err = errors.New(apiError.StatusMessage)
		return nil, err
	}

	userBytes, err := json.Marshal(resp)
	if err != nil {
		fmt.Println("Error:", err)
	}

	errPayment := repositories.SaveLogPaymentMidtrans(orderID, req, string(userBytes))
	if errPayment != nil {
		return nil, errPayment
	}

	errTransaction := repositories.CreateTransactionBilling(orderID, studendID, billingAmount, paymentMethodId, billingStudentIds, invoiceNumber, listAccountNumber, listBillingId, bankName, userId)
	if errTransaction != nil {
		return nil, errTransaction
	}

	return resp, nil
}

func CheckTransaction(orderID string) (*coreapi.TransactionStatusResponse, error) {
	c.New(os.Getenv("SERVER_KEY"), midtrans.Sandbox)
	res, errorMidtrans := c.CheckTransaction(orderID)
	if errorMidtrans != nil {
		apiError, err := ExtractErrorMessage(errorMidtrans.Message)
		if err != nil {
			return nil, err
		}
		err = errors.New(apiError.StatusMessage)
		return nil, err
	}
	fmt.Println("Response: ", res)

	err := repositories.SaveLogCheckPaymentMidtrans(orderID, res)
	if err != nil {
		return nil, err
	}

	// err = repositories.UpdateTransactionBilling(orderID, res.TransactionStatus)
	// if err != nil {
	// 	return nil, err
	// }

	return res, nil
}

func CancelTransaction(orderID string) (*coreapi.CancelResponse, error) {
	c.New(os.Getenv("SERVER_KEY"), midtrans.Sandbox)
	res, errorMidtrans := c.CancelTransaction(orderID)
	if errorMidtrans != nil {
		apiError, err := ExtractErrorMessage(errorMidtrans.Message)
		if err != nil {
			return nil, err
		}
		err = errors.New(apiError.StatusMessage)
		return nil, err
	}

	return res, nil
}

func RefundTransaction(orderID string) *coreapi.RefundResponse {
	refundRequest := &coreapi.RefundReq{
		Amount: 5000,
		Reason: "Item out of stock",
	}

	res, err := c.RefundTransaction(orderID, refundRequest)
	if err != nil {
		// do something on error handle
	}
	fmt.Println("Response: ", res)
	return nil
}

func ExtractErrorMessage(errorMessage string) (*response.MidtransExtractAPIError, error) {
	jsonStart := strings.Index(errorMessage, "{")
	if jsonStart == -1 {
		fmt.Println("No JSON found in the error response")
		return nil, fmt.Errorf("No JSON found in the error response")
	}

	jsonString := errorMessage[jsonStart:]

	// Unmarshal the JSON part into the MidtransAPIError struct
	var apiError response.MidtransExtractAPIError
	err := json.Unmarshal([]byte(jsonString), &apiError)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return nil, err
	}

	return &apiError, nil
}

func ChangeStatusTransactionFromWebhook(payload request.WebhookPayload) error {
	fmt.Println("Request from webhook: ", payload)

	err := repositories.SaveLogCheckPaymentMidtransFromWebhook(payload.OrderID, &payload)
	if err != nil {
		return err
	}

	err = repositories.UpdateTransactionBilling(payload) // Update the status
	if err != nil {
		return err
	}

	return nil
}

func GenerateSignature(orderID, statusCode, grossAmount, serverKey string) string {
	// Gabungkan data
	rawSignature := orderID + statusCode + grossAmount + serverKey

	// Hash dengan SHA-512
	hasher := sha512.New()
	hasher.Write([]byte(rawSignature))
	sha := hasher.Sum(nil)

	// Convert ke hex string
	return hex.EncodeToString(sha)
}

func ValidateSignature(orderID, statusCode, grossAmount, serverKey, midtransSignature string) error {
	// Generate the signature using the provided data
	localSignature := GenerateSignature(orderID, statusCode, grossAmount, serverKey)

	// Compare the generated signature with the one provided by Midtrans
	if localSignature != midtransSignature {
		return errors.New("signature mismatch: generated signature does not match Midtrans signature")
	}
	return nil
}
