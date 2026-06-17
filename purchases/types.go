package purchases

import "time"

type Customer struct {
	ID                int64     `json:"id"`
	Email             string    `json:"email"`
	GatewayCustomerID string    `json:"-"`
	CreatedAt         time.Time `json:"created_at"`
}

type Purchase struct {
	ID              int64     `json:"id"`
	CustomerID      int64     `json:"customer_id"`
	Amount          int       `json:"amount"`
	Currency        string    `json:"currency"`
	GatewayChargeID string    `json:"charge_id"`
	CreatedAt       time.Time `json:"created_at"`
}
