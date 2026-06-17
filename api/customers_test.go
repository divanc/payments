package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/divanc/payments/api"
	"github.com/divanc/payments/gateway"
	"github.com/divanc/payments/purchases"
	"github.com/divanc/payments/store"
)

func newHandler(t *testing.T) *api.Handler {
	t.Helper()
	repo, err := store.NewSQLite(":memory:")
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	svc := purchases.NewService(repo, gateway.NewFake())
	return api.NewHandler(svc)
}

func post(t *testing.T, h http.Handler, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewBufferString(body))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func get(t *testing.T, h http.Handler, path string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestGetCustomer_WithPurchases(t *testing.T) {
	h := newHandler(t)
	if rec := post(t, h, "/v1/customers", `{"email":"r@b.com"}`); rec.Code != http.StatusCreated {
		t.Fatalf("create customer status = %d, want 201", rec.Code)
	}
	if rec := post(t, h, "/v1/purchases", `{"customer_id":1,"amount":1999}`); rec.Code != http.StatusCreated {
		t.Fatalf("create purchase status = %d, want 201; body=%s", rec.Code, rec.Body.String())
	}

	rec := get(t, h, "/v1/customers/1")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}

	var got struct {
		Customer  map[string]any   `json:"customer"`
		Purchases []map[string]any `json:"purchases"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Customer["email"] != "r@b.com" {
		t.Errorf("customer.email = %v, want r@b.com", got.Customer["email"])
	}
	if len(got.Purchases) != 1 {
		t.Fatalf("purchases len = %d, want 1", len(got.Purchases))
	}
	if got.Purchases[0]["charge_id"] != "ch_fake_1" {
		t.Errorf("charge_id = %v, want ch_fake_1", got.Purchases[0]["charge_id"])
	}
}

func TestGetCustomer_EmptyPurchasesSerializeAsArray(t *testing.T) {
	h := newHandler(t)
	if rec := post(t, h, "/v1/customers", `{"email":"e@b.com"}`); rec.Code != http.StatusCreated {
		t.Fatalf("create customer status = %d, want 201", rec.Code)
	}

	rec := get(t, h, "/v1/customers/1")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	if !bytes.Contains(rec.Body.Bytes(), []byte(`"purchases":[]`)) {
		t.Errorf("purchases should serialize as [], got body=%s", rec.Body.String())
	}
}

func TestGetCustomer_UnknownID(t *testing.T) {
	h := newHandler(t)
	rec := get(t, h, "/v1/customers/999")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404; body=%s", rec.Code, rec.Body.String())
	}
}

func TestGetCustomer_NonIntegerID(t *testing.T) {
	h := newHandler(t)
	rec := get(t, h, "/v1/customers/abc")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404; body=%s", rec.Code, rec.Body.String())
	}
}

func TestCreateCustomer_Created(t *testing.T) {
	h := newHandler(t)
	rec := post(t, h, "/v1/customers", `{"email":"a@b.com"}`)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body=%s", rec.Code, rec.Body.String())
	}
	var got map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got["email"] != "a@b.com" {
		t.Errorf("email = %v, want a@b.com", got["email"])
	}
	if _, ok := got["id"]; !ok {
		t.Errorf("missing id in %v", got)
	}
	if _, ok := got["created_at"]; !ok {
		t.Errorf("missing created_at in %v", got)
	}
}

func TestCreateCustomer_DuplicateEmail(t *testing.T) {
	h := newHandler(t)
	if rec := post(t, h, "/v1/customers", `{"email":"dup@b.com"}`); rec.Code != http.StatusCreated {
		t.Fatalf("first create status = %d, want 201", rec.Code)
	}
	rec := post(t, h, "/v1/customers", `{"email":"dup@b.com"}`)
	if rec.Code != http.StatusConflict {
		t.Fatalf("status = %d, want 409; body=%s", rec.Code, rec.Body.String())
	}
}

func TestCreateCustomer_InvalidEmail(t *testing.T) {
	h := newHandler(t)
	rec := post(t, h, "/v1/customers", `{"email":"not-an-email"}`)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body=%s", rec.Code, rec.Body.String())
	}
}

func TestCreateCustomer_MalformedJSON(t *testing.T) {
	h := newHandler(t)
	rec := post(t, h, "/v1/customers", `{`)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body=%s", rec.Code, rec.Body.String())
	}
}
