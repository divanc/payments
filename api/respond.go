package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/divanc/payments/purchases"
)

// mapError translates a domain error into an HTTP status and a safe message.
func mapError(err error) (int, string) {
	switch {
	case errors.Is(err, purchases.ErrCustomerExists):
		return http.StatusConflict, "customer already exists"
	case errors.Is(err, purchases.ErrCustomerNotFound):
		return http.StatusNotFound, "customer not found"
	case errors.Is(err, purchases.ErrInvalidAmount):
		return http.StatusBadRequest, "amount must be positive"
	case errors.Is(err, purchases.ErrInvalidEmail):
		return http.StatusBadRequest, "invalid email"
	case errors.Is(err, purchases.ErrGateway):
		return http.StatusBadGateway, "payment failed"
	default:
		return http.StatusInternalServerError, "internal error"
	}
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}
