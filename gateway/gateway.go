package gateway

import "context"

// CustomerID and ChargeID are opaque identifiers minted by the payment provider.
type (
	CustomerID string
	ChargeID   string
)

// PaymentGateway is the port to a payment provider. Adapters: Fake (default) and Stripe.
type PaymentGateway interface {
	EnsureCustomer(ctx context.Context, email string) (CustomerID, error)
	Charge(ctx context.Context, c CustomerID, amount int) (ChargeID, error)
}
