package chapa

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

const (
	acceptPaymentV1APIURL  = "https://api.chapa.co/v1/transaction/initialize"
	verifyPaymentV1APIURL  = "https://api.chapa.co/v1/transaction/verify/%v"
	transferToBankV1APIURL = "https://api.chapa.co/v1/transfers"
	transactionsV1APIURL   = "https://api.chapa.co/v1/transactions"
	banksV1APIURL          = "https://api.chapa.co/v1/banks"
	bulkTransferAPIURL     = "https://api.chapa.co/v1/bulk-transfers"
)

type API interface {
	PaymentRequest(request *PaymentRequest) (*PaymentResponse, error)
	Verify(txnRef string) (*VerifyResponse, error)
	TransferToBank(request *BankTransfer) (*BankTransferResponse, error)
}

type chapa struct {
	apiKey string
	client *http.Client
}

func New() API {
	return &chapa{
		apiKey: viper.GetString("API_KEY"),
		client: &http.Client{
			Timeout: 1 * time.Minute,
		},
	}
}

func (c *chapa) PaymentRequest(request *PaymentRequest) (*PaymentResponse, error) {
	var err error
	if err = request.Validate(); err != nil {
		err := fmt.Errorf("invalid input %v", err)
		log.Printf("error %v input %v", err.Error(), request)
		return &PaymentResponse{}, err
	}

	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, acceptPaymentV1APIURL, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Close = true

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var paymentResponse PaymentResponse

	err = json.Unmarshal(body, &paymentResponse)
	if err != nil {
		return nil, err
	}
	return &paymentResponse, nil
}

func (c *chapa) Verify(txnRef string) (*VerifyResponse, error) {

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(verifyPaymentV1APIURL, txnRef), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Close = true

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var verifyResponse VerifyResponse

	err = json.Unmarshal(body, &verifyResponse)
	if err != nil {
		return nil, err
	}

	return &verifyResponse, nil
}

func (c *chapa) TransferToBank(request *BankTransfer) (*BankTransferResponse, error) {
	var err error
	if err = request.Validate(); err != nil {
		err := fmt.Errorf("invalid input %v", err)
		log.Printf("error %v input %v", err, request)
		return nil, err
	}

	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, transferToBankV1APIURL, bytes.NewBuffer(data))
	if err != nil {
		log.Printf("error %v", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Close = true

	resp, err := c.client.Do(req)
	if err != nil {
		log.Printf("error %v", err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error while reading resposne body %v", err)
		return nil, err
	}

	response := BankTransferResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error while unmarshaling  response %v", err)
		return nil, err
	}
	return &response, nil
}

func (c *chapa) getTransactions() (*TransactionsResponse, error) {
	req, err := http.NewRequest(http.MethodGet, transactionsV1APIURL, nil)
	if err == nil {
		log.Printf("error %v", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		log.Printf("error %v", err)
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error while reading response body %v", err)
		return nil, err
	}

	var response TransactionsResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error while unmarshaling  response %v", err)
		return nil, err
	}

	return &response, nil
}

func (c *chapa) getBanks() (*BanksResponse, error) {
	req, err := http.NewRequest(http.MethodGet, banksV1APIURL, nil)
	if err == nil {
		log.Printf("error %v", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		log.Printf("error %v", err)
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error while reading response body %v", err)
		return nil, err
	}

	var response BanksResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error while unmarshaling  response %v", err)
		return nil, err
	}

	return &response, nil
}

func (c *chapa) BulkTransfer(request *BulkTransferRequest) (*BulkTransferResponse, error) {
	var err error
	if err = request.Validate(); err != nil {
		err := fmt.Errorf("invalid input %v", err)
		log.Printf("error %v input %v", err, request)
		return nil, err
	}

	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, bulkTransferAPIURL, bytes.NewBuffer(data))
	if err != nil {
		log.Printf("error %v", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Close = true

	resp, err := c.client.Do(req)
	if err != nil {
		log.Printf("error %v", err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error while reading resposne body %v", err)
		return nil, err
	}

	response := BulkTransferResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error while unmarshaling  response %v", err)
		return nil, err
	}
	return &response, nil
}
