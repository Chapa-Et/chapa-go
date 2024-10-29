package chapa

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/shopspring/decimal"
)

type (
	PaymentRequest struct {
		Amount         decimal.Decimal        `json:"amount"`
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
			Charge float64 `json:"charge"`
		}
	}

	// BankTransfer is an object used in bank transfer.
	BankTransfer struct {
		// AccountName is the recipient Account Name matches on their bank account.
		AccountName string `json:"account_name"`
		// AccountNumber is the recipient Account Number.
		AccountNumber string `json:"account_number"`
		// Amount is the amount to be transferred to the recipient.
		Amount float64 `json:"amount"`
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

	Transaction struct {
		Status        TransactionStatus `json:"status"`
		RefID         string            `json:"ref_id"`
		Type          string            `json:"type"`
		CreatedAt     string            `json:"created_at"`
		Currency      string            `json:"currency"`
		Amount        decimal.Decimal   `json:"amount"`
		Charge        string            `json:"charge"`
		TransID       string            `json:"trans_id"`
		PaymentMethod string            `json:"payment_method"`
		Customer      Customer          `json:"customer"`
	}

	Customer struct {
		ID        int64  `json:"id"`
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Mobile    string `json:"mobile"`
	}

	Pagination struct {
		PerPage      int    `json:"per_page"`
		CurrentPage  int    `json:"current_page"`
		FirstPageURL string `json:"first_page_url"`
		NextPageURL  string `json:"next_page_url"`
		PrevPageURL  string `json:"prev_page_url"`
	}

	TransactionList struct {
		Transactions []Transaction `json:"transactions"`
		Pagination   Pagination    `json:"pagination"`
	}

	TransactionsResponse struct {
		Message string          `json:"message"`
		Status  string          `json:"status"`
		Data    TransactionList `json:"data"`
	}
	CheckoutForm struct {
		Amount   decimal.Decimal `json:"amount"`
		Currency string          `json:"currency"`
	}

	TransactionStatus string

	Currency string

	Bank struct {
		ID            int64    `json:"id"`
		Swift         string   `json:"swift"`
		Name          string   `json:"name"`
		AcctLength    int64    `json:"acct_length"`
		CountryID     int64    `json:"country_id"`
		CreatedAt     string   `json:"created_at"`
		UpdatedAt     string   `json:"updated_at"`
		IsRTGS        int64    `json:"is_rtgs"`
		IsMobileMoney int64    `json:"is_mobilemoney"`
		Currency      Currency `json:"currency"`
	}

	BanksResponse struct {
		Message string `json:"message"`
		Data    []Bank `json:"data"`
	}

	BulkData struct {
		AccountName   string `json:"account_name"`
		AccountNumber string `json:"account_number"`
		Amount        int64  `json:"amount"`
		Reference     string `json:"reference"`
		BankCode      string `json:"bank_code"`
	}

	BulkTransferRequest struct {
		Title    string     `json:"title"`
		Currency string     `json:"currency"`
		BulkData []BulkData `json:"bulk_data"`
	}

	BulkTransferResponseData struct {
		ID        int    `json:"id"`
		CreatedAt string `json:"created_at"`
	}

	BulkTransferResponse struct {
		Message string                   `json:"message"`
		Status  string                   `json:"status"`
		Data    BulkTransferResponseData `json:"data"`
	}
)

const (
	FailedTransactionStatus  TransactionStatus = "failed"
	PendingTransactionStatus TransactionStatus = "pending"
	SuccessTransactionStatus TransactionStatus = "success"
	ETB                      Currency          = "ETB"
	USD                      Currency          = "USD"
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

func (t BulkTransferRequest) Validate() error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.Title, validation.Required.Error("title of the bulk transfer is required")),
		validation.Field(&t.Currency, validation.Required.Error("currency is required")),
		validation.Field(&t.BulkData, validation.NilOrNotEmpty.Error("at least one account is required")),
	)
}
