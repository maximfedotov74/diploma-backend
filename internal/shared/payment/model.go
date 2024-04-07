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

type PaymentMethod struct {
	Type  string            `json:"type"`
	ID    string            `json:"id"`
	Saved bool              `json:"saved"`
	Card  PaymentMethodCard `json:"card"`
	Title string            `json:"title"`
}

type PaymentMethodCard struct {
	First6                   string                   `json:"first6"`
	Last4                    string                   `json:"last4"`
	ExpiryMonth              string                   `json:"expiry_month"`
	ExpiryYear               string                   `json:"expiry_year"`
	CardType                 string                   `json:"card_type"`
	IssuerCountry            string                   `json:"issuer_country"`
	IssuerName               string                   `json:"issuer_name"`
	PaymentMethodCardProduct PaymentMethodCardProduct `json:"card_product"`
}

type PaymentMethodCardProduct struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type OrderPayment struct {
	ID            string        `json:"id"`
	Status        string        `json:"status"`
	Paid          bool          `json:"paid"`
	Amount        Amount        `json:"amount"`
	CreatedAt     time.Time     `json:"created_at"`
	Description   string        `json:"description"`
	ExpiresAt     time.Time     `json:"expires_at"`
	PaymentMethod PaymentMethod `json:"payment_method"`
	Recipient     Recipient     `json:"recipient"`
	Refundable    bool          `json:"refundable"`
	Test          bool          `json:"test"`
}

type RefundDto struct {
	Amount    Amount `json:"amount"`
	PaymentId string `json:"payment_id"`
}

type RefundRespomse struct {
	Id        string    `json:"id"`
	Status    string    `json:"status"`
	Amount    Amount    `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
	PaymentId string    `json:"payment_id"`
}
