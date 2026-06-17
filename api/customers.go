package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/divanc/payments/purchases"
)

func (h *Handler) createCustomer(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	c, err := h.svc.CreateCustomer(r.Context(), req.Email)
	if err != nil {
		code, msg := mapError(err)
		writeError(w, code, msg)
		return
	}
	writeJSON(w, http.StatusCreated, c)
}

func (h *Handler) getCustomer(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		code, msg := mapError(purchases.ErrCustomerNotFound)
		writeError(w, code, msg)
		return
	}

	c, ps, err := h.svc.GetCustomer(r.Context(), id)
	if err != nil {
		code, msg := mapError(err)
		writeError(w, code, msg)
		return
	}

	writeJSON(w, http.StatusOK, struct {
		Customer  purchases.Customer   `json:"customer"`
		Purchases []purchases.Purchase `json:"purchases"`
	}{Customer: c, Purchases: ps})
}
