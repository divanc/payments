//go:build e2e

package e2e

import (
	"fmt"
	"net/http"
	"testing"
)

func TestCreateCustomer(t *testing.T) {
	email := uniqueEmail()
	resp := postJSON(t, "/v1/customers", fmt.Sprintf(`{"email":%q}`, email))
	wantStatus(t, resp, http.StatusCreated)

	var c struct {
		ID        int64  `json:"id"`
		Email     string `json:"email"`
		CreatedAt string `json:"created_at"`
	}
	decode(t, resp, &c)
	if c.ID == 0 || c.Email != email || c.CreatedAt == "" {
		t.Fatalf("unexpected customer body: %+v", c)
	}
}

func TestCreateCustomer_DuplicateEmail(t *testing.T) {
	email := uniqueEmail()
	body := fmt.Sprintf(`{"email":%q}`, email)

	first := postJSON(t, "/v1/customers", body)
	wantStatus(t, first, http.StatusCreated)
	first.Body.Close()

	dup := postJSON(t, "/v1/customers", body)
	wantStatus(t, dup, http.StatusConflict)
	dup.Body.Close()
}

func TestCreateCustomer_InvalidEmail(t *testing.T) {
	resp := postJSON(t, "/v1/customers", `{"email":"not-an-email"}`)
	wantStatus(t, resp, http.StatusBadRequest)
	resp.Body.Close()
}

func TestCreateCustomer_MalformedJSON(t *testing.T) {
	resp := postJSON(t, "/v1/customers", `{`)
	wantStatus(t, resp, http.StatusBadRequest)
	resp.Body.Close()
}

func TestCreatePurchase(t *testing.T) {
	id := createCustomer(t)
	resp := postJSON(t, "/v1/purchases", fmt.Sprintf(`{"customer_id":%d,"amount":500}`, id))
	wantStatus(t, resp, http.StatusCreated)

	var p struct {
		ID         int64  `json:"id"`
		CustomerID int64  `json:"customer_id"`
		Amount     int    `json:"amount"`
		Currency   string `json:"currency"`
		ChargeID   string `json:"charge_id"`
	}
	decode(t, resp, &p)
	if p.CustomerID != id || p.Amount != 500 || p.ChargeID == "" {
		t.Fatalf("unexpected purchase body: %+v", p)
	}
}

func TestCreatePurchase_UnknownCustomer(t *testing.T) {
	resp := postJSON(t, "/v1/purchases", `{"customer_id":999999,"amount":500}`)
	wantStatus(t, resp, http.StatusNotFound)
	resp.Body.Close()
}

func TestCreatePurchase_NonPositiveAmount(t *testing.T) {
	id := createCustomer(t)
	resp := postJSON(t, "/v1/purchases", fmt.Sprintf(`{"customer_id":%d,"amount":0}`, id))
	wantStatus(t, resp, http.StatusBadRequest)
	resp.Body.Close()
}

func TestGetCustomer_WithPurchases(t *testing.T) {
	id := createCustomer(t)

	mk := postJSON(t, "/v1/purchases", fmt.Sprintf(`{"customer_id":%d,"amount":750}`, id))
	wantStatus(t, mk, http.StatusCreated)
	mk.Body.Close()

	resp := getJSON(t, fmt.Sprintf("/v1/customers/%d", id))
	wantStatus(t, resp, http.StatusOK)

	var out struct {
		Customer struct {
			ID int64 `json:"id"`
		} `json:"customer"`
		Purchases []struct {
			Amount   int    `json:"amount"`
			ChargeID string `json:"charge_id"`
		} `json:"purchases"`
	}
	decode(t, resp, &out)
	if out.Customer.ID != id {
		t.Fatalf("customer id = %d, want %d", out.Customer.ID, id)
	}
	if len(out.Purchases) != 1 || out.Purchases[0].Amount != 750 || out.Purchases[0].ChargeID == "" {
		t.Fatalf("unexpected purchases: %+v", out.Purchases)
	}
}

func TestGetCustomer_UnknownID(t *testing.T) {
	resp := getJSON(t, "/v1/customers/999999")
	wantStatus(t, resp, http.StatusNotFound)
	resp.Body.Close()
}
