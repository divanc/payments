package api_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"testing"

	"github.com/divanc/payments/api"
	"github.com/divanc/payments/gateway"
	"github.com/divanc/payments/purchases"
	"github.com/divanc/payments/store"
)

// failingGateway charges fail to exercise the 502 path.
type failingGateway struct{}

func (failingGateway) EnsureCustomer(ctx context.Context, email string) (gateway.CustomerID, error) {
	return "cus_test", nil
}

func (failingGateway) Charge(ctx context.Context, c gateway.CustomerID, amount int) (gateway.ChargeID, error) {
	return "", errors.New("provider down")
}

func itoa(n int64) string { return strconv.FormatInt(n, 10) }

func newHandlerWith(t *testing.T, gw gateway.PaymentGateway) *api.Handler {
	t.Helper()
	repo, err := store.NewSQLite(":memory:")
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	return api.NewHandler(purchases.NewService(repo, gw))
}

// createCustomer posts a customer and returns its id.
func createCustomer(t *testing.T, h http.Handler, email string) int64 {
	t.Helper()
	rec := post(t, h, "/v1/customers", `{"email":"`+email+`"}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create customer status = %d, want 201; body=%s", rec.Code, rec.Body.String())
	}
	var got struct {
		ID int64 `json:"id"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode customer: %v", err)
	}
	return got.ID
}

func TestCreatePurchase_Created(t *testing.T) {
	h := newHandler(t)
	id := createCustomer(t, h, "buyer@b.com")

	rec := post(t, h, "/v1/purchases", `{"customer_id":`+itoa(id)+`,"amount":1999}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body=%s", rec.Code, rec.Body.String())
	}
	var got map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got["amount"] != float64(1999) {
		t.Errorf("amount = %v, want 1999", got["amount"])
	}
	if got["currency"] != "usd" {
		t.Errorf("currency = %v, want usd", got["currency"])
	}
	if got["charge_id"] != "ch_fake_1" {
		t.Errorf("charge_id = %v, want ch_fake_1", got["charge_id"])
	}
	if _, ok := got["id"]; !ok {
		t.Errorf("missing id in %v", got)
	}
	if _, ok := got["created_at"]; !ok {
		t.Errorf("missing created_at in %v", got)
	}
}

func TestCreatePurchase_NonPositiveAmount(t *testing.T) {
	h := newHandler(t)
	id := createCustomer(t, h, "buyer@b.com")

	rec := post(t, h, "/v1/purchases", `{"customer_id":`+itoa(id)+`,"amount":0}`)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body=%s", rec.Code, rec.Body.String())
	}
}

func TestCreatePurchase_MalformedJSON(t *testing.T) {
	h := newHandler(t)
	rec := post(t, h, "/v1/purchases", `{`)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body=%s", rec.Code, rec.Body.String())
	}
}

func TestCreatePurchase_UnknownCustomer(t *testing.T) {
	h := newHandler(t)
	rec := post(t, h, "/v1/purchases", `{"customer_id":999,"amount":1999}`)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404; body=%s", rec.Code, rec.Body.String())
	}
}

func TestCreatePurchase_GatewayError(t *testing.T) {
	h := newHandlerWith(t, failingGateway{})
	id := createCustomer(t, h, "buyer@b.com")

	rec := post(t, h, "/v1/purchases", `{"customer_id":`+itoa(id)+`,"amount":1999}`)
	if rec.Code != http.StatusBadGateway {
		t.Fatalf("status = %d, want 502; body=%s", rec.Code, rec.Body.String())
	}
}
