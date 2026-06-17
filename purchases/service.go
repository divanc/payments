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
