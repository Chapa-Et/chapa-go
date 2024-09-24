package chapa

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

// Placeholder data
var (
	firstName1 = "Jon"
	firstName2 = "Do"
	lastName1  = "Mary"
	lastName2  = "Josef"
	email1     = RandomString(5) + "@gmail.com"
	email2     = RandomString(5) + "@gmail.com"

	Customers = []Customer{
		{
			ID:        1002,
			FirstName: &firstName1,
			LastName:  &lastName1,
			Email:     &email1,
		},
		{
			ID:        1032,
			FirstName: &firstName2,
			LastName:  &lastName2,
			Email:     &email2,
		},
	}

	transactions = []Transaction{
		{
			TransID:   RandomString(10),
			Amount:    "10.00",
			Charge:    "0.35",
			Currency:  "ETB",
			CreatedAt: time.Now().String(),
			Customer:  &Customers[0],
		},
		{
			TransID:   RandomString(10),
			Amount:    "20.00",
			Charge:    "0.40",
			Currency:  "ETB",
			CreatedAt: time.Now().String(),
			Customer:  &Customers[0],
		},
	}
)

type (
	ExamplePaymentService interface {
		Checkout(ctx context.Context, CustomerID int64, form *CheckoutForm) (*PaymentResponse, error)
		ListTransactions(ctx context.Context) (*TransactionList, error)
	}

	AppExamplePaymentService struct {
		mu                     *sync.Mutex
		paymentGatewayProvider API
	}
)

func NewExamplePaymentService(
	paymentGatewayProvider API,
) *AppExamplePaymentService {
	return &AppExamplePaymentService{
		mu:                     &sync.Mutex{},
		paymentGatewayProvider: paymentGatewayProvider,
	}
}

func (s *AppExamplePaymentService) Checkout(ctx context.Context, CustomerID int64, form *CheckoutForm) (*Transaction, error) {

	Customer, err := s.CustomerByID(ctx, CustomerID)
	if err != nil {
		return &Transaction{}, err
	}

	invoice := &PaymentRequest{
		Amount:         form.Amount,
		Currency:       form.Currency,
		Email:          *Customer.Email,
		FirstName:      *Customer.FirstName,
		LastName:       *Customer.LastName,
		CallbackURL:    "https://webhook.site/077164d6-29cb-40df-ba29-8a00e59a7e60",
		TransactionRef: RandomString(10),
	}

	response, err := s.paymentGatewayProvider.PaymentRequest(invoice)
	if err != nil {
		return &Transaction{}, err
	}

	if response.Status != "success" {

		// log the response
		log.Printf("[ERROR] Failed to checkout Customer request response = [%+v]", response)

		return &Transaction{}, fmt.Errorf("failed to checkout err = %v", response.Message)
	}

	transaction := Transaction{
		TransID:   invoice.TransactionRef,
		Amount:    fmt.Sprintf("%.2f", form.Amount),
		Currency:  form.Currency,
		Customer:  Customer,
		Status:    PendingTransactionStatus,
		CreatedAt: time.Now().String(),
	}

	err = s.saveTransaction(ctx, transaction)
	if err != nil {
		return &Transaction{}, nil
	}

	return &transaction, nil
}

func (s *AppExamplePaymentService) ListTransactions(ctx context.Context) (*TransactionList, error) {

	// validations here

	transactionList := &TransactionList{
		Transactions: transactions,
	}

	return transactionList, nil
}

func (s *AppExamplePaymentService) saveTransaction(ctx context.Context, transaction Transaction) error {

	s.mu.Lock()
	defer s.mu.Unlock()

	transactions = append([]Transaction{transaction}, transactions...)

	return nil
}

// CustomerByID - normally you'd fetch Customer from the db
func (s *AppExamplePaymentService) CustomerByID(ctx context.Context, CustomerID int64) (*Customer, error) {

	for index := range Customers {
		if Customers[index].ID == CustomerID {
			return &Customers[index], nil
		}
	}

	return &Customer{}, errors.New("Customer not found")
}
