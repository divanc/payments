package api

import (
	"context"
	"net/http"

	"github.com/divanc/payments/purchases"
)

// Service is the business-logic port the HTTP layer depends on.
type Service interface {
	CreateCustomer(ctx context.Context, email string) (purchases.Customer, error)
}

type Handler struct {
	svc Service
	mux *http.ServeMux
}

func NewHandler(svc Service) *Handler {
	h := &Handler{svc: svc, mux: http.NewServeMux()}
	h.mux.HandleFunc("POST /v1/customers", h.createCustomer)
	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}
