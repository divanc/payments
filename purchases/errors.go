package purchases

import "errors"

var (
	ErrCustomerExists   = errors.New("customer already exists")
	ErrCustomerNotFound = errors.New("customer not found")
	ErrInvalidAmount    = errors.New("amount must be positive")
	ErrInvalidEmail     = errors.New("invalid email")
	ErrGateway          = errors.New("payment gateway error")
)
