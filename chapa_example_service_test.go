package chapa

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChapaExampleService(t *testing.T) {
	InitConfig()
	ctx := context.Background()
	exampleService := NewExamplePaymentService(
		New(),
	)

	t.Run("can list payment transactions", func(t *testing.T) {
		transactionList, err := exampleService.ListPaymentTransactions(ctx)
		assert.NoError(t, err)

		assert.Equal(t, 2, len(transactionList.Transactions))
	})

	t.Run("can successfully checkout", func(t *testing.T) {
		form := &CheckoutForm{
			Amount:   12.30,
			Currency: "ETB",
		}

		paymentTxn, err := exampleService.Checkout(ctx, 1032, form)
		assert.NoError(t, err)

		assert.Equal(t, form.Amount, paymentTxn.Amount)
		assert.Equal(t, form.Currency, paymentTxn.Currency)
		assert.Equal(t, PendingTransactionStatus, paymentTxn.Status)
		assert.Zero(t, paymentTxn.MerchantFee)
		assert.NotZero(t, paymentTxn.TransactionID)

		assert.Equal(t, 3, len(transactions))
	})

	t.Run("cannot checkout if user is unavailable", func(t *testing.T) {
		form := &CheckoutForm{
			Amount:   12.30,
			Currency: "ETB",
		}

		_, err := exampleService.Checkout(ctx, 0, form)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
	})
}
