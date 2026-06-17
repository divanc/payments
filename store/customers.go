package store

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/divanc/payments/purchases"
)

func (r *sqliteRepo) CreateCustomer(
	ctx context.Context,
	email, gatewayCustomerID string,
) (purchases.Customer, error) {
	now := time.Now().UTC().Truncate(time.Second)
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO customers(email, gateway_customer_id, created_at) VALUES(?, ?, ?)`,
		email, gatewayCustomerID, now.Format(time.RFC3339))
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return purchases.Customer{}, purchases.ErrCustomerExists
		}
		return purchases.Customer{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return purchases.Customer{}, err
	}
	return purchases.Customer{
		ID:                id,
		Email:             email,
		GatewayCustomerID: gatewayCustomerID,
		CreatedAt:         now,
	}, nil
}

func (r *sqliteRepo) GetCustomerByID(ctx context.Context, id int64) (purchases.Customer, error) {
	var (
		c       purchases.Customer
		created string
	)
	err := r.db.QueryRowContext(ctx,
		`SELECT id, email, gateway_customer_id, created_at FROM customers WHERE id = ?`, id).
		Scan(&c.ID, &c.Email, &c.GatewayCustomerID, &created)
	if errors.Is(err, sql.ErrNoRows) {
		return purchases.Customer{}, purchases.ErrCustomerNotFound
	}
	if err != nil {
		return purchases.Customer{}, err
	}
	c.CreatedAt = parseTime(created)
	return c, nil
}
