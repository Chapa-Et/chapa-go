package chapa

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type (
	PaymentRequest struct {
		Amount         float64                `json:"amount"`
		Currency       string                 `json:"currency"`
		Email          string                 `json:"email"`
		FirstName      string                 `json:"first_name"`
		LastName       string                 `json:"last_name"`
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

func (c PaymentRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Email, validation.Required.Error("email is required"), is.Email),
		validation.Field(&c.FirstName, validation.Required.Error("first name is required")),
		validation.Field(&c.LastName, validation.Required.Error("last name is required")),
		validation.Field(&c.TransactionRef, validation.Required.Error("transaction reference is required")),
		validation.Field(&c.Currency, validation.Required.Error("currency is required")),
		validation.Field(&c.Amount, validation.Required.Error("amount is required")),
		validation.Field(&c.CallbackURL, validation.Required.Error("callback url is required"), is.URL),
	)
}
