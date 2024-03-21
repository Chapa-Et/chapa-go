package chapa

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	PaymentRequest struct {
		Amount         float64                `json:"amount"`
		Currency       string                 `json:"currency"`
		Email          string                 `json:"email"`
		FirstName      string                 `json:"first_name"`
		LastName       string                 `json:"last_name"`
		Phone          string                 `json:"phone"`
		CallbackURL    string                 `json:"callback_url"`
		TransactionRef string                 `json:"tx_ref"`
		Customization  map[string]interface{} `json:"customization"`
	}

	PaymentResponse struct {
		Message string `json:"message"`
		Status  string `json:"status"`
		Data    struct {
			CheckoutURL string `json:"checkout_url"`
		}
	}

	VerifyResponse struct {
		Message string `json:"message"`
		Status  string `json:"status"`
		Data    struct {
			TransactionFee float64 `json:"charge"`
		}
	}
)

type (
	// BankTransfer is an object used in bank transfer.
	BankTransfer struct {
		// AccountName is the recipient Account Name matches on their bank account.
		AccountName string `json:"account_name"`
		// AccountNumber is the recipient Account Number.
		AccountNumber string `json:"account_number"`
		// Amount is the amount to be transferred to the recipient.
		Amount float64 `json:"amount"`
		// BeneficiaryName is the full name of the Transfer beneficiary (You may use it to match on your required).
		BeneficiaryName string `json:"beneficiary_name"`
		// Currency is the currency for the Transfer. Expected value is ETB.
		Currency string `json:"currency"`
		// Reference is merchant’s uniques reference for the transfer,
		// it can be used to query for the status of the transfer.
		Reference string `json:"reference"`
		// BankCode is the recipient bank code.
		// You can see a list of all the available banks and their codes from the get banks endpoint.
		BankCode string `json:"bank_code"`
	}

	BankTransferResponse struct {
		Message string `json:"message"`
		Status  string `json:"status"`
		Data    string `json:"data"`
	}
)

func (p PaymentRequest) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.TransactionRef, validation.Required.Error("transaction reference is required")),
		validation.Field(&p.Currency, validation.Required.Error("currency is required")),
		validation.Field(&p.Amount, validation.Required.Error("amount is required")),
	)
}

func (t BankTransfer) Validate() error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.AccountName, validation.Required.Error("account name is required")),
		validation.Field(&t.AccountNumber, validation.Required.Error("account number is required")),
		validation.Field(&t.Amount, validation.Required.Error("amount is required")),
		validation.Field(&t.Currency, validation.Required.Error("currency is required")),
		validation.Field(&t.Reference, validation.Required.Error("reference is required")),
		validation.Field(&t.BankCode, validation.Required.Error("bank code is required")),
	)
}
