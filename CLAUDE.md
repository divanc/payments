# payments

Small HTTP payments API over a payment gateway, backed by SQLite.

## Commands

- Run: `go run .`
- Test: `go test ./...`
- Build / vet / format: `go build ./...`, `go vet ./...`, `gofmt -w .`

## Architecture

Hexagonal. Dependencies point inward:

```
api  ->  purchases.Service  ->  gateway.PaymentGateway (port)
                             ->  store.Repository       (port)
```

- Ports are declared at their consumer: `Repository` in `purchases`, `Service` in `api`.
- Domain types (`Customer`, `Purchase`) and sentinel errors live in `purchases`.
- Concrete adapters (`gateway.Fake`, `store` SQLite) are wired in `main`.
- The default gateway is the in-memory `Fake`; the Stripe adapter is selected via env.

## Conventions

- One file per resource within a package (`api/customers.go`, `store/customers.go`); keep
  package plumbing separate (`api/router.go`, `api/respond.go`, `store/sqlite.go`).
- HTTP concerns stay in `api`; business rules stay in `purchases`. The domain never imports
  `net/http`. Domain errors map to status codes only in `api/respond.go` (via `errors.Is`).
- Wrap long function signatures: one parameter per line, same-type params grouped, trailing
  comma, return types after the closing paren.
- Money is integer minor units (cents). Timestamps are RFC3339 UTC, truncated to seconds.
