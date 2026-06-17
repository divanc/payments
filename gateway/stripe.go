package gateway

import (
	"context"
	"errors"
	"strings"

	"github.com/stripe/stripe-go/v82"
)

const testKeyPrefix = "sk_test_"

// Stripe is a PaymentGateway backed by the Stripe API in test mode.
type Stripe struct {
	client *stripe.Client
}

// NewStripe constructs a Stripe gateway. It accepts only test-mode secret keys
// (prefix "sk_test_") so a live key can never place a real charge.
func NewStripe(key string) (*Stripe, error) {
	if !strings.HasPrefix(key, testKeyPrefix) {
		return nil, errors.New("gateway: STRIPE_SECRET_KEY must be a test-mode key (sk_test_…)")
	}

	return &Stripe{client: stripe.NewClient(key)}, nil
}

func (s *Stripe) EnsureCustomer(ctx context.Context, email string) (CustomerID, error) {
	c, err := s.client.V1Customers.Create(ctx, &stripe.CustomerCreateParams{
		Email: stripe.String(email),
	})
	if err != nil {
		return "", err
	}

	return CustomerID(c.ID), nil
}

func (s *Stripe) Charge(ctx context.Context, c CustomerID, amount int) (ChargeID, error) {
	pi, err := s.client.V1PaymentIntents.Create(ctx, &stripe.PaymentIntentCreateParams{
		Amount:        stripe.Int64(int64(amount)),
		Currency:      stripe.String(string(stripe.CurrencyUSD)),
		Customer:      stripe.String(string(c)),
		PaymentMethod: stripe.String("pm_card_visa"),
		AutomaticPaymentMethods: &stripe.PaymentIntentCreateAutomaticPaymentMethodsParams{
			Enabled:        stripe.Bool(true),
			AllowRedirects: stripe.String(string(stripe.PaymentIntentAutomaticPaymentMethodsAllowRedirectsNever)),
		},
	})
	if err != nil {
		return "", err
	}
	pi, err = s.client.V1PaymentIntents.Confirm(ctx, pi.ID, &stripe.PaymentIntentConfirmParams{})
	if err != nil {
		return "", err
	}
	return ChargeID(pi.ID), nil
}
