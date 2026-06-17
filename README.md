# payments

A small HTTP payments API over a payment gateway, backed by SQLite.

## Running

```sh
go run .
```

The server listens on `:8080` by default. It uses an in-memory fake payment gateway
by default, so no API key or network access is required to run locally.

## Configuration

All configuration is via environment variables.

| Variable            | Default              | Description                                          |
| ------------------- | -------------------- | ---------------------------------------------------- |
| `PORT`              | `8080`               | TCP port the server listens on.                      |
| `DB_PATH`           | `sqlite/payments.db` | Path to the SQLite database file.                    |
| `GATEWAY`           | `fake`               | Payment gateway to use: `fake` or `stripe`.          |
| `STRIPE_SECRET_KEY` | —                    | Stripe secret key, required when `GATEWAY=stripe`.   |

## Endpoints

### Create a customer

```sh
curl -s -X POST localhost:8080/v1/customers \
  -H 'Content-Type: application/json' \
  -d '{"email":"a@b.com"}'
```

```json
{"id":1,"email":"a@b.com","created_at":"2026-06-17T00:00:00Z"}
```

### Create a purchase

Charges the customer through the gateway and records the purchase.

```sh
curl -s -X POST localhost:8080/v1/purchases \
  -H 'Content-Type: application/json' \
  -d '{"customer_id":1,"amount":1999}'
```

```json
{"id":1,"customer_id":1,"amount":1999,"currency":"usd","charge_id":"ch_fake_1","created_at":"2026-06-17T00:00:00Z"}
```

Amounts are in integer minor units (cents).

### Get a customer with purchases

```sh
curl -s localhost:8080/v1/customers/1
```

```json
{"customer":{"id":1,"email":"a@b.com","created_at":"2026-06-17T00:00:00Z"},"purchases":[{"id":1,"customer_id":1,"amount":1999,"currency":"usd","charge_id":"ch_fake_1","created_at":"2026-06-17T00:00:00Z"}]}
```

## Testing

```sh
go test ./...
```
