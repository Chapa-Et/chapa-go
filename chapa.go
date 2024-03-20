package chapa

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	acceptPaymentV1APIURL = "https://api.chapa.co/v1/transaction/initialize"
	verifyPaymentV1APIURL = "https://api.chapa.co/v1/transaction/verify/%v"
)

type API interface {
	PaymentRequest(request *PaymentRequest) (*PaymentResponse, error)
	Verify(txnRef string) (*VerifyResponse, error)
}

type chap struct {
	apiKey string
	client *http.Client
}

func New() API {
	return &chap{
		apiKey: viper.GetString("API_KEY"),
		client: &http.Client{
			Timeout: 1 * time.Minute,
		},
	}
}

func (c *chap) PaymentRequest(request *PaymentRequest) (*PaymentResponse, error) {
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

func (c *chap) Verify(txnRef string) (*VerifyResponse, error) {

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

	body, err := ioutil.ReadAll(resp.Body)
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
