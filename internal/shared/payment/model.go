package payment

import "time"

type Payment struct {
	ID           string               `json:"id"`
	Status       string               `json:"status"`
	Paid         bool                 `json:"paid"`
	Amount       Amount               `json:"amount"`
	Confirmation ConfirmationResponse `json:"confirmation"`
	CreatedAt    time.Time            `json:"created_at"`
	Description  string               `json:"description"`
	Recipient    Recipient            `json:"recipient"`
	Refundable   bool                 `json:"refundable"`
	Test         bool                 `json:"test"`
}

type Recipient struct {
	AccountID string `json:"account_id"`
	GatewayID string `json:"gateway_id"`
}

type Amount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type ConfirmationResponse struct {
	Type            string `json:"type"`
	ConfirmationURL string `json:"confirmation_url"`
}

type Confirmation struct {
	Type      string `json:"type"`
	ReturnURL string `json:"return_url"`
}

type PaymentDto struct {
	Capture      bool         `json:"capture"`
	Description  string       `json:"description"`
	Amount       Amount       `json:"amount"`
	Confirmation Confirmation `json:"confirmation"`
}
