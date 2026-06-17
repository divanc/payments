package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/divanc/payments/api"
	"github.com/divanc/payments/gateway"
	"github.com/divanc/payments/purchases"
	"github.com/divanc/payments/store"
)

func main() {
	port := envOr("PORT", "8080")
	dbPath := envOr("DB_PATH", "sqlite/payments.db")

	if dir := filepath.Dir(dbPath); dir != "." && dbPath != ":memory:" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			log.Fatalf("create db dir: %v", err)
		}
	}

	repo, err := store.NewSQLite(dbPath)
	if err != nil {
		log.Fatalf("open store: %v", err)
	}

	gw, err := newGateway(os.Getenv("GATEWAY"), os.Getenv("STRIPE_SECRET_KEY"))
	if err != nil {
		log.Fatalf("select gateway: %v", err)
	}
	if f, ok := gw.(*gateway.Fake); ok {
		if d := os.Getenv("FAKE_CHARGE_DELAY"); d != "" {
			delay, err := time.ParseDuration(d)
			if err != nil {
				log.Fatalf("parse FAKE_CHARGE_DELAY: %v", err)
			}
			f.ChargeLatency = delay
		}
	}
	svc := purchases.NewService(repo, gw)
	handler := api.NewHandler(svc)

	log.Printf("listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

// newGateway selects a payment gateway adapter. GATEWAY defaults to "fake";
// "stripe" requires a test-mode STRIPE_SECRET_KEY.
func newGateway(name, stripeKey string) (gateway.PaymentGateway, error) {
	switch name {
	case "", "fake":
		return gateway.NewFake(), nil
	case "stripe":
		return gateway.NewStripe(stripeKey)
	default:
		return nil, fmt.Errorf("unknown GATEWAY %q (want fake or stripe)", name)
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
