CREATE TABLE IF NOT EXISTS customers (
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    email               TEXT    NOT NULL UNIQUE,
    gateway_customer_id TEXT    NOT NULL,
    created_at          TEXT    NOT NULL
);

CREATE TABLE IF NOT EXISTS purchases (
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    customer_id       INTEGER NOT NULL REFERENCES customers(id),
    amount            INTEGER NOT NULL,
    currency          TEXT    NOT NULL,
    gateway_charge_id TEXT    NOT NULL,
    created_at        TEXT    NOT NULL
);
