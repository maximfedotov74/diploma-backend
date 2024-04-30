package payment

//https://yookassa.ru/developers/payment-acceptance/getting-started/quick-start
import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/maximfedotov74/diploma-backend/internal/shared/fall"
)

type PaymentService struct {
	shopId     string
	secretKey  string
	appLink    string
	httpClient *http.Client
}

const paymentUrl = "https://api.yookassa.ru/v3/payments"
const refundUrl = "https://api.yookassa.ru/v3/refunds"

func NewPaymentService(shopId string, secretKey string, appLink string) *PaymentService {
	return &PaymentService{shopId: shopId, secretKey: secretKey, appLink: appLink, httpClient: &http.Client{}}
}

func (ps *PaymentService) CheckPayment(paymentId string) (*OrderPayment, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", paymentUrl, paymentId), nil)
	if err != nil {
		return nil, err
	}

	authString := fmt.Sprintf("%s:%s", ps.shopId, ps.secretKey)
	authBase64 := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(authString)))
	req.Header.Set("Authorization", authBase64)

	response, err := ps.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("ошибка при получении платежа, http код: %d", response.StatusCode)
	}

	var p OrderPayment
	err = json.Unmarshal(bytes, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (ps *PaymentService) RefundPayment(paymentId string, totalPrice float64) (*RefundResponse, error) {

	dto := RefundDto{
		Amount:    Amount{Value: fmt.Sprintf("%.2f", totalPrice), Currency: "RUB"},
		PaymentId: paymentId,
	}

	dtoBytes, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", refundUrl, bytes.NewReader(dtoBytes))
	if err != nil {
		return nil, err
	}
	idempotenceKey := uuid.New().String()
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotence-Key", idempotenceKey)
	authString := fmt.Sprintf("%s:%s", ps.shopId, ps.secretKey)
	authBase64 := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(authString)))
	req.Header.Set("Authorization", authBase64)

	response, err := ps.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	log.Println(response.StatusCode)

	if response.StatusCode != fall.STATUS_OK {
		return nil, errors.New("ошибка при оформлении возврата")
	}

	var p RefundResponse
	err = json.Unmarshal(bytes, &p)
	if err != nil {
		return nil, err
	}
	log.Println(p)
	return &p, nil
}

func (ps *PaymentService) CreatePayment(orderId string, totalPrice float64) (*Payment, error) {
	dto := PaymentDto{
		Amount:      Amount{Value: fmt.Sprintf("%.2f", totalPrice), Currency: "RUB"},
		Capture:     true,
		Description: fmt.Sprintf("Оплата заказа №%s в магазине FamilyModa", orderId),
		Confirmation: Confirmation{
			Type:      "redirect",
			ReturnURL: ps.appLink + "/api/order/confirm-online-payment/" + orderId,
		},
	}
	dtoBytes, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", paymentUrl, bytes.NewReader(dtoBytes))
	if err != nil {
		return nil, err
	}
	idempotenceKey := uuid.New().String()
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotence-Key", idempotenceKey)
	authString := fmt.Sprintf("%s:%s", ps.shopId, ps.secretKey)
	authBase64 := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(authString)))
	req.Header.Set("Authorization", authBase64)
	response, err := ps.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, errors.New("ошибка при обработке платежа заказа №" + orderId)
	}
	var p Payment
	err = json.Unmarshal(bytes, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
