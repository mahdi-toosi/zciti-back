package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go-fiber-starter/utils/config"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
)

type SepGateway struct {
	*PaymentService
}

//nolint:tagliatelle // bank need to use PascalCase
type PaymentRequest struct {
	TerminalID       string `json:"TerminalId"`
	ResNum           string `json:"ResNum"`
	RedirectURL      string `json:"RedirectUrl"`
	CellNumber       string `json:"CellNumber"`
	Action           string `json:"action"`
	HashedCardNumber string `json:"HashedCardNumber,omitempty"`
	TokenExpiryInMin int    `json:"TokenExpiryInMin,omitempty"`
	Amount           int    `json:"Amount"`
}

//nolint:tagliatelle // bank need to use PascalCase
type CallbackResponse struct {
	MID              string `json:"MID"`
	State            string `json:"State"`
	Status           string `json:"Status"`
	RRN              string `json:"RRN"`
	RefNum           string `json:"RefNum"`
	ResNum           string `json:"ResNum"`
	TerminalID       string `json:"TerminalId"`
	TraceNo          string `json:"TraceNo"`
	SecurePan        string `json:"SecurePan"`
	HashedCardNumber string `json:"HashedCardNumber"`
	Amount           int    `json:"Amount"`
	Wage             int    `json:"Wage,omitempty"`
}

//nolint:tagliatelle // bank need to use PascalCase
type VerifyRequest struct {
	RefNum         string `json:"RefNum"`
	TerminalNumber int64  `json:"TerminalNumber"`
}

//nolint:tagliatelle // bank need to use PascalCase
type VerifyResponse struct {
	ResultDescription string `json:"ResultDescription"`
	TransactionDetail struct {
		RRN             string `json:"RRN"`
		RefNum          string `json:"RefNum"`
		MaskedPan       string `json:"MaskedPan"`
		HashedPan       string `json:"HashedPan"`
		StraceDate      string `json:"StraceDate"`
		StraceNo        string `json:"StraceNo"`
		TerminalNumber  int32  `json:"TerminalNumber"`
		OriginalAmount  int64  `json:"OrginalAmount"`
		AffectiveAmount int64  `json:"AffectiveAmount"`
	} `json:"TransactionDetail"`
	Success    bool `json:"Success"`
	ResultCode int  `json:"ResultCode"`
}

type PaymentService struct {
	Logger     *zerolog.Logger
	TerminalID string
}

func NewSepGateway(cfg *config.Config, logger zerolog.Logger) *SepGateway {
	terminalID := cfg.Services.Saman.TerminalID
	if terminalID == "" {
		logger.Panic().Msg("Saman TerminalID is not configured")
	}

	service, err := NewPaymentService(terminalID, &logger)
	if err != nil {
		logger.Panic().Err(err).Msg("Failed to initialize Saman payment service")
	}

	return &SepGateway{
		PaymentService: service,
	}
}

func NewPaymentService(terminalID string, logger *zerolog.Logger) (*PaymentService, error) {
	return &PaymentService{
		TerminalID: terminalID,
		Logger:     logger,
	}, nil
}

const (
	requestURL = "https://sep.shaparak.ir/OnlinePG/OnlinePG"
	payURL     = "https://sep.shaparak.ir/OnlinePG/SendToken"
	verifyURL  = "https://sep.shaparak.ir/verifyTxnRandomSessionkey/ipg/VerifyTransaction"
	reverseURL = "https://sep.shaparak.ir/verifyTxnRandomSessionkey/ipg/ReverseTransaction"
)

func (ps *PaymentService) SendRequest(amount int, resNum, cellNumber, redirectURL string) (string, error) {
	paymentRequest := &PaymentRequest{
		Action:      "token",
		Amount:      amount,
		ResNum:      resNum,
		CellNumber:  cellNumber,
		RedirectURL: redirectURL,
		TerminalID:  ps.TerminalID,
	}

	jsonData, err := json.Marshal(paymentRequest)
	if err != nil {
		return "", fmt.Errorf("خطا در تبدیل درخواست به JSON: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("خطا در ساخت درخواست HTTP: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("خطا در ارسال درخواست HTTP: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			ps.Logger.Error().Err(err).Msg("Error closing response body")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("خطا: %s", GetSamanError(resp.StatusCode))
	}

	var responseMap map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&responseMap); err != nil {
		return "", fmt.Errorf("خطا در پردازش پاسخ: %w", err)
	}

	if responseMap["errorCode"] != nil {
		errorCode, err := strconv.Atoi(responseMap["errorCode"].(string))
		if err != nil {
			return "", fmt.Errorf("خطا در پردازش پاسخ: %w", err)
		}
		return "", fmt.Errorf("خطا: %s", GetSamanError(errorCode))
	}

	token, ok := responseMap["token"].(string)
	if !ok {
		return "", fmt.Errorf("خطا: %s", GetSamanError(10))
	}

	paymentURL := fmt.Sprintf("%s?token=%s", payURL, token)

	return paymentURL, nil
}

func (ps *PaymentService) Verify(ctx context.Context, refNum string) (*VerifyResponse, error) {
	data := url.Values{}
	data.Set("TerminalNumber", ps.TerminalID)
	data.Set("RefNum", refNum)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, verifyURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			ps.Logger.Error().Err(closeErr).Msg("failed to close response body")
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if code, err := strconv.Atoi(string(body)); err == nil {
		return nil, fmt.Errorf("خطا در تایید: %s", GetSamanVerifyAndReverseError(code))
	}

	var result VerifyResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if result.ResultCode != 0 {
		return nil, fmt.Errorf("خطا در تایید: %s", GetSamanVerifyAndReverseError(result.ResultCode))
	}

	return &result, nil
}

var errorMessages = map[int]string{
	2:  "پرداخت با موفقیت انجام شد",
	1:  "کاربر از پرداخت انصراف داده است",
	3:  "پذیرنده فروشگاهی نامعتبر است",
	4:  "کاربر در بازه زمانی تعیین شده پاسخی ارسال نکرده است",
	5:  "پارامترهای ارسالی نامعتبر است",
	8:  "آدرس سرور پذیرنده نامعتبر است",
	10: "توکن ارسال شده یافت نشد",
	11: "با این شماره ترمینال فقط تراکنش های توکنی قابل پرداخت هستند.",
	12: "شماره ترمینال ارسال شده یافت نشد",
	21: "محدودیت های مدل چند حسابی رعایت نشده",
}

var verifyAndReverseErrorMessages = map[int]string{
	-9999: "دریافت خطای استثنا",
	-9998: "دریافت تایم اوت 65 ثانیه ای",
	-106:  "آدرس آی پی درخواستی غیرمجاز می باشد.",
	-105:  "ترمینال ارسالی در سیستم موجود نمی باشد.",
	-104:  "ترمینال ارسالی غیرفعال می باشد.",
	-18:   "IP Address فروشنده نامعتبر است.",
	-14:   "چنین تراکنشی تعریف نشده است.",
	-11:   "طول ورودی ها کمتر از حد مجاز است.",
	-10:   "رسید دیجیتالی به صورت Base64 نیست.(حاوی کاراکترهای غیرمجاز است)",
	-8:    "طول ورودی ها بیشتر از حد مجاز است.",
	-7:    "رسید دیجیتال تهی است.",
	-6:    "بیش از نیم ساعت از زمان اجرای تراکنش گذشته است.",
	-2:    "تراکنش یافت نشد.",
	0:     "موفق",
	1:     "کاربر از پرداخت انصراف داده است",
	2:     "درخواست تکراری می باشد.",
	5:     "تراکنش برگشت خورده می باشد.",
}

func GetSamanError(code int) string {
	if errMsg, ok := errorMessages[code]; ok {
		return errMsg
	}
	return fmt.Sprintf("کد خطا: %d. خطای ناشناخته.", code)
}

func GetSamanVerifyAndReverseError(code int) string {
	if errMsg, ok := verifyAndReverseErrorMessages[code]; ok {
		return errMsg
	}
	return fmt.Sprintf("کد خطا: %d. خطای ناشناخته.", code)
}
