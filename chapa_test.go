package chapa

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChapa(t *testing.T) {
	var request *PaymentRequest

	t.Run("chap API", func(t *testing.T) {
		InitConfig()
		paymentProvider := New()

		t.Run("can prompt payment from users", func(t *testing.T) {
			request = &PaymentRequest{
				Amount:         10,
				Currency:       "ETB",
				FirstName:      "chap",
				LastName:       "ET",
				Email:          "chap@et.io",
				CallbackURL:    "https://webhook.site/077164d6-29cb-40df-ba29-8a00e59a7e60",
				TransactionRef: RandomString(20),
				Customization: map[string]interface{}{
					"title":       "title",
					"description": "description",
					"logo":        "https://company.com/logo",
				},
			}

			response, err := paymentProvider.PaymentRequest(request)
			assert.NoError(t, err)

			assert.Equal(t, "success", response.Status)
			assert.Equal(t, "Hosted Link", response.Message)
			assert.Contains(t, response.Data.CheckoutURL, "https://checkout.chapa.co/checkout/payment")
		})

		t.Run("can verify transactions", func(t *testing.T) {
			response, err := paymentProvider.Verify(request.TransactionRef) // a paid txn
			assert.NoError(t, err)

			assert.Equal(t, "success", response.Status)
			assert.Equal(t, "Payment details fetched successfully", response.Message)
			//assert.NotZero(t, response.Data.TransactionFee)   // uncomment this for live mode
		})

		t.Run("can get transactions", func(t *testing.T) {
			response, err := paymentProvider.getTransactions()
			assert.NoError(t, err)

			assert.Equal(t, "success", response.Status)
			assert.Equal(t, "Transactions retrieved successfully", response.Message)
		})

		t.Run("can get banks", func(t *testing.T) {
			response, err := paymentProvider.getBanks()
			assert.NoError(t, err)

			assert.Equal(t, "Banks retrieved", response.Message)
		})

		t.Run("cannot verify unpaid transaction", func(t *testing.T) {
			request := &PaymentRequest{
				Amount:         10,
				Currency:       "ETB",
				FirstName:      "chap",
				LastName:       "ET",
				Email:          "chap@et.io",
				CallbackURL:    "",
				TransactionRef: RandomString(20),
				Customization: map[string]interface{}{
					"title":       "A Unique Title",
					"description": "This a perfect description",
					"logo":        "https://your.logo",
				},
			}

			response, err := paymentProvider.Verify(request.TransactionRef)
			assert.NoError(t, err)
			assert.Equal(t, "Invalid transaction reference", response.Message)
		})

		t.Run("successful bank transfer", func(t *testing.T) {
			request := &BankTransfer{
				AccountName:   "Leul Abay Ejigu",
				AccountNumber: "1000212482106",
				Amount:        10,
				Currency:      "ETB",
				Reference:     "3241342142sfdd",
				BankCode:      "946",
			}

			response, err := paymentProvider.TransferToBank(request)
			assert.NoError(t, err)

			assert.Equal(t, "success", response.Status)
			// update below assertion on live mode
			assert.Equal(t, "Transfer queued successfully in Test Mode.", response.Message)
			assert.NotEmpty(t, response.Data)
		})

		t.Run("invalid input bank transfer", func(t *testing.T) {
			request := &BankTransfer{
				AccountNumber: "34264263",
				Amount:        10,
				Currency:      "ETB",
				Reference:     "3264063st01",
				BankCode:      "32735b19-bb36-4cd7-b226-fb7451cd98f0",
			}

			response, err := paymentProvider.TransferToBank(request)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "invalid input")
			assert.Nil(t, response)
		})

		t.Run("successful bulk transfer", func(t *testing.T) {
			bulkData := BulkData{
				AccountName:   "Leul Abay Ejigu",
				AccountNumber: "1000212482106",
				Amount:        10,
				Reference:     "3241342142sfdd",
				BankCode:      "946",
			}

			request := &BulkTransferRequest{
				Title:    "Transfer to leul",
				Currency: "ETB",
				BulkData: []BulkData{bulkData},
			}

			response, err := paymentProvider.bulkTransfer(request)
			assert.NoError(t, err)

			assert.Equal(t, "success", response.Status)
			// update below assertion on live mode
			// assert.Equal(t, "Bulk transfer queued successfully", response.Message)
			assert.Equal(t, "Bulk Transfer queued successfully in Test Mode.", response.Message)
			assert.NotEmpty(t, response.Data)
		})

		t.Run("invalid input for bulk transfer", func(t *testing.T) {
			request := &BulkTransferRequest{
				Title:    "",
				Currency: "ETB",
				BulkData: []BulkData{},
			}

			response, err := paymentProvider.bulkTransfer(request)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "invalid input")
			assert.Nil(t, response)
		})
	})
}
