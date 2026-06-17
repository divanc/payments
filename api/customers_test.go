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
