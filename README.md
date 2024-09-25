<h1 align="center">
<div align="center">
  <a href="https://chapa.co/" target="_blank">
    <img src="./docs/logo.png" width="320" alt="Nest Logo"/>
  </a>
  <p align="center">Go SDK for chapa</p>
</div>
</h1>

![build-workflow](https://github.com/Chapa-Et/chapa-go/actions/workflows/test.yml/badge.svg)

Unofficial Golang SDK for Chapa ET API

### Todo:
- [ ] We could add nice validations on demand.
- [ ] Add implementation for the remaining API endpoints. 


### Usage
##### 1. Installation
```
    go get github.com/Chapa-Et/chapa-go
```

###### API_KEY
Add your `API_KEY: CHASECK_xxxxxxxxxxxxxxxx` inside `config.yaml` file.
If you want to run the githb action on your forked repository, you have to create a secrete key named `API_KEY`.

##### 2. Setup

```go
    package main

    import (
        chapa "github.com/Chapa-Et/chapa-go"
    )

    func main(){
        chapaAPI := chapa.New()
    }
```

##### 3. Accept Payments
```go
    request := &chapaAPI.PaymentRequest{
        Amount:         10,
        Currency:       "ETB",
        FirstName:      "Chapa",
        LastName:       "ET",
        Email:          "chapa@et.io",
        CallbackURL:    "https://posthere.io/e631-44fe-a19e",
        TransactionRef: RandomString(20),
        Customization: map[string]interface{}{
            "title":       "A Unique Title",
            "description": "This a perfect description",
            "logo":        "https://your.logo",
        },
    }

    response, err := chapaAPI.PaymentRequest(request)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    fmt.Printf("payment response: %+v\n", response)
```

##### 4. Verify Payment Transactions
```go
    response, err := chapaAPI.Verify("your-txn-ref")
    if err != nil {
         fmt.Println(err)
         os.Exit(1)
    }

    fmt.Printf("verification response: %+v\n", response)
```

##### 5. Transfer to bank
```go
    request := &BankTransfer{
	    AccountName:     "Yinebeb Tariku", 
	    AccountNumber:   "34264263", 
	    Amount:          10,
	    BeneficiaryName: "Yinebeb Tariku",
	    Currency:        "ETB",
	    Reference:       "3264063st01",
	    BankCode:        "32735b19-bb36-4cd7-b226-fb7451cd98f0",
	}
		
	response, err := chapaAPI.TransferToBank(request)
	fmt.Printf("transfer response: %+v\n", response)
```

##### 6. Get transactions
```go	
	response, err := chapaAPI.getTransactions()
	fmt.Printf("transactions response: %+v\n", response)
```

##### 7. Get banks
```go	
	response, err := chapaAPI.getBanks()
	fmt.Printf("banks response: %+v\n", response)
```

##### 8. Bulk transfer
```go	
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

    response, err := chapaAPI.bulkTransfer(request)
    fmt.Printf("bulk transfer response: %+v\n", response)
```

### Resources
- https://developer.chapa.co/docs/overview/

### Quirks
Suggestions on how to improve the API:
- Introduction of `status codes` would be a nice to have in the future. Status codes are better than the `message` in a way considering there are so many reasons a transaction could fail.
e.g 
```shell
    1001: Success
    4001: DuplicateTransaction
    4002: InvalidCurrency
    4003: InvalidAmount 
    4005: InsufficientBalance  
    5000: InternalServerError
    5001: GatewayError
    5002: RejectedByGateway
```
Just an example!
### Contributions
- Highly welcome
