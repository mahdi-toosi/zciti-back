package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"go-fiber-starter/utils/config"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// ZarinPal is the base struct for zarinPal payment
// gateway, one shall not create or manipulate instances
// if this struct manually and just use provided methods
// to woke with it.
type ZarinPal struct {
	MerchantID      string
	Sandbox         bool
	APIEndpoint     string
	PaymentEndpoint string
}

type paymentRequestReqBody struct {
	MerchantID  string
	Amount      int
	CallbackURL string
	Description string
	Email       string
	Mobile      string
}

type paymentRequestResp struct {
	Status    int
	Authority string
}

type paymentVerificationReqBody struct {
	MerchantID string
	Authority  string
	Amount     int
}

type paymentVerificationResp struct {
	Status int
	RefID  json.Number
}

type unverifiedTransactionsReqBody struct {
	MerchantID string
}

// UnverifiedAuthority is the base struct for Authorities in unverifiedTransactionsResp
type UnverifiedAuthority struct {
	Authority   string
	Amount      int
	Channel     string
	CallbackURL string
	Referer     string
	Email       string
	CellPhone   string
	Date        string // ToDo Check type to be date
}

type unverifiedTransactionsResp struct {
	Status      int
	Authorities []UnverifiedAuthority
}

type refreshAuthorityReqBody struct {
	MerchantID string
	Authority  string
	ExpireIn   int
}

type refreshAuthorityResp struct {
	Status int
}

// NewZarinPal creates a new instance of zarinPal payment
// gateway with provided configs. It also tries to validate
// provided configs.
func NewZarinPal(cfg *config.Config) *ZarinPal {
	merchantID := cfg.Services.ZarinPal.MerchantID
	sandbox := cfg.Services.ZarinPal.Sandbox

	if len(merchantID) != 36 {
		panic("MerchantID must be 36 characters")
	}
	apiEndPoint := "https://www.zarinpal.com/pg/rest/WebGate/"
	paymentEndpoint := "https://www.zarinpal.com/pg/StartPay/"
	if sandbox {
		apiEndPoint = "https://sandbox.zarinpal.com/pg/rest/WebGate/"
		paymentEndpoint = "https://sandbox.zarinpal.com/pg/StartPay/"
	}
	return &ZarinPal{
		Sandbox:         sandbox,
		MerchantID:      merchantID,
		APIEndpoint:     apiEndPoint,
		PaymentEndpoint: paymentEndpoint,
	}
}

// NewPaymentRequest gets a payment url from ZarinPal.
// amount is in Tomans (not Rials) format.
// email and mobile are optional.
//
// If error is not nil, you can check statusCode for
// specific error handling based on ZarinPal error codes.
// If statusCode is not 100, it means ZarinPal raised an error
// on their end and you can check the error code and its reason
// based on their documentation placed in
// https://github.com/ZarinPal-Lab/Documentation-PaymentGateway/archive/master.zip
func (zarinPal *ZarinPal) NewPaymentRequest(amount int, callbackURL, description, email, mobile string) (paymentURL, authority string, statusCode int, err error) {
	if amount < 1 {
		err = errors.New("amount must be a positive number")
		return
	}
	if callbackURL == "" {
		err = errors.New("callbackURL should not be empty")
		return
	}
	if description == "" {
		err = errors.New("description should not be empty")
		return
	}
	paymentRequest := paymentRequestReqBody{
		MerchantID:  zarinPal.MerchantID,
		Amount:      amount,
		CallbackURL: callbackURL,
		Description: description,
		Email:       email,
		Mobile:      mobile,
	}
	var resp paymentRequestResp
	err = zarinPal.request("PaymentRequest.json", &paymentRequest, &resp)
	if err != nil {
		return
	}
	statusCode = resp.Status
	if resp.Status == 100 {
		authority = resp.Authority
		paymentURL = zarinPal.PaymentEndpoint + resp.Authority
	} else {
		err = errors.New(strconv.Itoa(resp.Status))
	}
	return
}

// PaymentVerification verifies if a payment was done successfully, Authority of the
// payment request should be passed to this method alongside its Amount in Tomans.
//
// If error is not nil, you can check statusCode for
// specific error handling based on ZarinPal error codes.
// If statusCode is not 100, it means ZarinPal raised an error
// on their end and you can check the error code and its reason
// based on their documentation placed in
// https://github.com/ZarinPal-Lab/Documentation-PaymentGateway/archive/master.zip
func (zarinPal *ZarinPal) PaymentVerification(amount int, authority string) (verified bool, refID string, statusCode int, err error) {
	if amount <= 0 {
		err = errors.New("amount must be a positive number")
		return
	}
	if authority == "" {
		err = errors.New("authority should not be empty")
		return
	}
	paymentVerification := paymentVerificationReqBody{
		MerchantID: zarinPal.MerchantID,
		Amount:     amount,
		Authority:  authority,
	}
	var resp paymentVerificationResp
	err = zarinPal.request("PaymentVerification.json", &paymentVerification, &resp)
	if err != nil {
		return
	}
	statusCode = resp.Status
	if resp.Status == 100 {
		verified = true
		refID = string(resp.RefID)
	} else {
		err = errors.New(strconv.Itoa(resp.Status))
	}
	return
}

// UnverifiedTransactions gets unverified transactions.
//
// If error is not nil, you can check statusCode for
// specific error handling based on ZarinPal error codes.
// If statusCode is not 100, it means ZarinPal raised an error
// on their end and you can check the error code and its reason
// based on their documentation placed in
// https://github.com/ZarinPal-Lab/Documentation-PaymentGateway/archive/master.zip
func (zarinPal *ZarinPal) UnverifiedTransactions() (authorities []UnverifiedAuthority, statusCode int, err error) {
	unverifiedTransactions := unverifiedTransactionsReqBody{
		MerchantID: zarinPal.MerchantID,
	}

	var resp unverifiedTransactionsResp
	err = zarinPal.request("UnverifiedTransactions.json", &unverifiedTransactions, &resp)
	if err != nil {
		return
	}

	if resp.Status == 100 {
		statusCode = resp.Status
		authorities = resp.Authorities
	} else {
		err = errors.New(strconv.Itoa(resp.Status))
	}
	return
}

// RefreshAuthority update authority expiration time.\n
// expire should be number between [1800,3888000] seconds.
//
// If error is not nil, you can check statusCode for
// specific error handling based on ZarinPal error codes.
// If statusCode is not 100, it means ZarinPal raised an error
// on their end and you can check the error code and its reason
// based on their documentation placed in
// https://github.com/ZarinPal-Lab/Documentation-PaymentGateway/archive/master.zip
func (zarinPal *ZarinPal) RefreshAuthority(authority string, expire int) (statusCode int, err error) {
	if authority == "" {
		err = errors.New("authority should not be empty")
		return
	}
	if expire < 1800 {
		err = errors.New("expire must be at least 1800")
		return
	} else if expire > 3888000 {
		err = errors.New("expire must not be greater than 3888000")
		return
	}

	refreshAuthority := refreshAuthorityReqBody{
		MerchantID: zarinPal.MerchantID,
		Authority:  authority,
		ExpireIn:   expire,
	}
	var resp refreshAuthorityResp
	err = zarinPal.request("RefreshAuthority.json", &refreshAuthority, &resp)
	if err != nil {
		return
	}
	if resp.Status == 100 {
		statusCode = resp.Status
	} else {
		err = errors.New(strconv.Itoa(resp.Status))
	}
	return
}

func (zarinPal *ZarinPal) request(method string, data interface{}, res interface{}) error {
	reqBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", zarinPal.APIEndpoint+method, bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Println(string(body))
	err = json.Unmarshal(body, res)
	if err != nil {
		err = errors.New("zarinPal invalid json response")
		return err
	}
	return nil
}
