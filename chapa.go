package chapa

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	chapaAcceptPaymentV1APIURL = "https://api.chapa.co/v1/transaction/initialize"
	chapaVerifyPaymentV1APIURL = "https://api.chapa.co/v1/transaction/verify/%v"
)

type (
	ChapaAPI interface {
		PaymentRequest(request *PaymentRequest) (*PaymentResponse, error)
		Verify(txnRef string) (*VerifyResponse, error)
	}

	Chapa struct {
		apiKey string
		client *http.Client
	}
)

func New(apiKey string) *Chapa {
	return &Chapa{
		apiKey: apiKey,
		client: &http.Client{
			Timeout: 1 * time.Minute,
		},
	}
}

func (c *Chapa) PaymentRequest(request *PaymentRequest) (*PaymentResponse, error) {
	var err error
	if err = request.Validate(); err != nil {
		err := fmt.Errorf("invalid input %v", err)
		log.Printf("error %v input %v", err.Error(), request)
		return &PaymentResponse{}, err
	}

	data, err := json.Marshal(request)
	if err != nil {
		return &PaymentResponse{}, err
	}

	req, err := http.NewRequest(http.MethodPost, chapaAcceptPaymentV1APIURL, bytes.NewBuffer(data))
	if err != nil {
		return &PaymentResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Close = true

	resp, err := c.client.Do(req)
	if err != nil {
		return &PaymentResponse{}, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &PaymentResponse{}, err
	}

	var chapaPaymentResponse PaymentResponse

	err = json.Unmarshal(body, &chapaPaymentResponse)
	if err != nil {
		return &PaymentResponse{}, err
	}

	return &chapaPaymentResponse, nil
}

func (c *Chapa) Verify(txnRef string) (*VerifyResponse, error) {

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(chapaVerifyPaymentV1APIURL, txnRef), nil)
	if err != nil {
		return &VerifyResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Close = true

	resp, err := c.client.Do(req)
	if err != nil {
		return &VerifyResponse{}, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &VerifyResponse{}, err
	}

	var chapaVerifyResponse VerifyResponse

	err = json.Unmarshal(body, &chapaVerifyResponse)
	if err != nil {
		return &VerifyResponse{}, err
	}

	return &chapaVerifyResponse, nil
}
