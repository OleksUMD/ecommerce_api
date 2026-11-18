package providers

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/OleksUMD/ecommerce_api/internal/config"
	"github.com/OleksUMD/ecommerce_api/internal/dto"
	"github.com/stripe/stripe-go/v87"
	"github.com/stripe/stripe-go/v87/paymentintent"
	"github.com/stripe/stripe-go/v87/webhook"
)

type PaymentProvider struct {
	Config *config.PaymentConfig
}

func NewPaymentProvider(cfg *config.PaymentConfig) *PaymentProvider {
	return &PaymentProvider{Config: cfg}
}

func (p *PaymentProvider) CreatePaymentIntentResponse(amount float64, idempotentKey string) (*dto.PaymentIntentResponse, error) {
	stripe.Key = p.Config.SecretKey
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(math.Round(amount * 100))),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
	}

	params.SetIdempotencyKey(idempotentKey)
	pi, err := paymentintent.New(params)
	if err != nil {
		return nil, err
	}
	return &dto.PaymentIntentResponse{ID: pi.ID, ClientSecret: pi.ClientSecret}, nil
}

func (p *PaymentProvider) ProcessWebHook(sigHeader string, reqBody []byte) (string, string, error) {
	event, err := webhook.ConstructEvent(reqBody, sigHeader, p.Config.WebhookSecret)
	if err != nil {
		return "", "", err
	}
	var paymentIntent stripe.PaymentIntent
	if event.Type == "payment_intent.succeeded" {
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			return "", "", err
		}
	} else {
		return "", "", fmt.Errorf("payment failed: %s", event.Type)
	}
	return event.ID, paymentIntent.ID, nil
}
