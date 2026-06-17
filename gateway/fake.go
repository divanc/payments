package gateway

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"
)

// Fake is an in-memory PaymentGateway for local development and tests. It mints
// deterministic identifiers and tracks how many charges it has placed.
type Fake struct {
	customers atomic.Int64
	charges   atomic.Int64

	// ChargeLatency optionally delays each Charge to simulate provider round-trips.
	ChargeLatency time.Duration
}

func NewFake() *Fake { return &Fake{} }

func (f *Fake) EnsureCustomer(ctx context.Context, email string) (CustomerID, error) {
	n := f.customers.Add(1)
	return CustomerID(fmt.Sprintf("cus_fake_%d", n)), nil
}

func (f *Fake) Charge(ctx context.Context, c CustomerID, amount int) (ChargeID, error) {
	if f.ChargeLatency > 0 {
		time.Sleep(f.ChargeLatency)
	}
	n := f.charges.Add(1)
	return ChargeID(fmt.Sprintf("ch_fake_%d", n)), nil
}

// Charges reports the number of charges placed so far.
func (f *Fake) Charges() int { return int(f.charges.Load()) }
