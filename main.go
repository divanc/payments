package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

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

	gw := gateway.NewFake()
	svc := purchases.NewService(repo, gw)
	handler := api.NewHandler(svc)

	log.Printf("listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
