package purchases

import (
	"context"
	"fmt"
	"strings"

	"github.com/divanc/payments/gateway"
)

// Repository persists customers and purchases.
type Repository interface {
	CreateCustomer(ctx context.Context, email, gatewayCustomerID string) (Customer, error)
	GetCustomerByID(ctx context.Context, id int64) (Customer, error)
	InsertPurchase(
		ctx context.Context,
		customerID int64,
		amount int,
		currency, gatewayChargeID string,
	) (Purchase, error)
	ListPurchases(ctx context.Context, customerID int64) ([]Purchase, error)
}

// Service holds the payment business logic. It depends only on ports.
type Service struct {
	repo Repository
	gw   gateway.PaymentGateway
}

func NewService(repo Repository, gw gateway.PaymentGateway) *Service {
	return &Service{repo: repo, gw: gw}
}

func (s *Service) CreateCustomer(ctx context.Context, email string) (Customer, error) {
	email = strings.TrimSpace(email)
	if email == "" || !strings.Contains(email, "@") {
		return Customer{}, ErrInvalidEmail
	}

	gid, err := s.gw.EnsureCustomer(ctx, email)
	if err != nil {
		return Customer{}, fmt.Errorf("%w: %v", ErrGateway, err)
	}

	return s.repo.CreateCustomer(ctx, email, string(gid))
}

// CreatePurchase charges an existing customer and records the purchase. The flow is
// charge-then-insert: the customer is verified, the gateway is charged, and the result
// is persisted. There is no rollback if the insert fails after a successful charge.
func (s *Service) CreatePurchase(
	ctx context.Context,
	customerID int64,
	amount int,
) (Purchase, error) {
	if amount <= 0 {
		return Purchase{}, ErrInvalidAmount
	}

	c, err := s.repo.GetCustomerByID(ctx, customerID)
	if err != nil {
		return Purchase{}, err
	}

	chargeID, err := s.gw.Charge(ctx, gateway.CustomerID(c.GatewayCustomerID), amount)
	if err != nil {
		return Purchase{}, fmt.Errorf("%w: %v", ErrGateway, err)
	}

	return s.repo.InsertPurchase(ctx, customerID, amount, "usd", string(chargeID))
}
