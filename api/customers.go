package api

import (
	"encoding/json"
	"net/http"
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
