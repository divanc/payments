package store

import (
	"context"
	"time"

	"github.com/divanc/payments/purchases"
)

func (r *sqliteRepo) InsertPurchase(
	ctx context.Context,
	customerID int64,
	amount int,
	currency, gatewayChargeID string,
) (purchases.Purchase, error) {
	now := time.Now().UTC().Truncate(time.Second)
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO purchases(customer_id, amount, currency, gateway_charge_id, created_at) VALUES(?, ?, ?, ?, ?)`,
		customerID, amount, currency, gatewayChargeID, now.Format(time.RFC3339))
	if err != nil {
		return purchases.Purchase{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return purchases.Purchase{}, err
	}
	return purchases.Purchase{
		ID:              id,
		CustomerID:      customerID,
		Amount:          amount,
		Currency:        currency,
		GatewayChargeID: gatewayChargeID,
		CreatedAt:       now,
	}, nil
}

func (r *sqliteRepo) ListPurchases(ctx context.Context, customerID int64) ([]purchases.Purchase, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, customer_id, amount, currency, gateway_charge_id, created_at
		   FROM purchases WHERE customer_id = ? ORDER BY id`, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []purchases.Purchase{}
	for rows.Next() {
		var (
			p       purchases.Purchase
			created string
		)
		if err := rows.Scan(&p.ID, &p.CustomerID, &p.Amount, &p.Currency, &p.GatewayChargeID, &created); err != nil {
			return nil, err
		}
		p.CreatedAt = parseTime(created)
		out = append(out, p)
	}
	return out, rows.Err()
}
