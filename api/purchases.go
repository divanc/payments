package api

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) createPurchase(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CustomerID int64 `json:"customer_id"`
		Amount     int   `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	p, err := h.svc.CreatePurchase(r.Context(), req.CustomerID, req.Amount)
	if err != nil {
		code, msg := mapError(err)
		writeError(w, code, msg)
		return
	}
	writeJSON(w, http.StatusCreated, p)
}
