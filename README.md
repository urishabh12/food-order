# Order Food (Go + Chi)

A tiny REST API for listing products and placing orders — implemented in Go with the Chi router.
It uses **in-memory data** for products and order responses, plus **promo code validation** loaded from three local files.

---

## Endpoints

* `GET {BASE_PATH}/healthz` – liveness check
* `GET {BASE_PATH}/product` – list products
* `GET {BASE_PATH}/product/{productId}` – get product by ID (`int64` path param)
* `POST {BASE_PATH}/order` – place order (**requires** header `api_key: apitest`)

> Default `BASE_PATH` is `/api`, so URLs are `/api/product`, etc.

---

## Quick Start (Docker Compose)

1. Build & run:

```bash
docker compose up --build
```

2. Tail logs:

```bash
docker compose logs -f api
# listening on :8080 (basePath /api)
```

---

## Run with Plain Docker

```bash
# Build
docker build -t order-food-api:latest .

# Run (mount coupons, set envs)
docker run --rm -p 8080:8080 \
  -e ADDR=":8080" \
  -e BASE_PATH="/api" \
  -e API_KEY="apitest" \
  -e COUPONS_DIR="/app/coupons" \
  -v "$(pwd)/coupons:/app/coupons:ro" \
  order-food-api:latest
```

---

**Health**

```bash
curl -i http://localhost:8080/api/healthz
```

**List products**

```bash
curl -s http://localhost:8080/api/product | jq
```

**Get product by ID**

```bash
curl -s http://localhost:8080/api/product/10 | jq
# 400 if ID is not int-like; 404 if not found
```

**Place order (no coupon)**

```bash
curl -s -X POST http://localhost:8080/api/order \
  -H "Content-Type: application/json" \
  -H "api_key: apitest" \
  -d '{
        "items": [
          {"productId": "10", "quantity": 2},
          {"productId": "31", "quantity": 1}
        ]
      }' | jq
```

**Place order (valid coupon)**

> Coupon must be **8–10 chars** and present in **≥ 2** coupon files.

```bash
curl -s -X POST http://localhost:8080/api/order \
  -H "Content-Type: application/json" \
  -H "api_key: apitest" \
  -d '{
        "couponCode": "SUPER100",
        "items": [{"productId": "10", "quantity": 1}]
      }' | jq
```

**Place order (invalid coupon)**

```bash
curl -i -X POST http://localhost:8080/api/order \
  -H "Content-Type: application/json" \
  -H "api_key: apitest" \
  -d '{
        "couponCode": "SUPER1000",
        "items": [{"productId": "10", "quantity": 1}]
      }'
# => 422 Validation exception
```

---

## Configuration

| Variable      | Default     | Description                             |
| ------------- | ----------- | --------------------------------------- |
| `ADDR`        | `:8080`     | Server listen address                   |
| `BASE_PATH`   | `/api`      | API base path (use `/` for root)        |
| `API_KEY`     | `apitest`   | Required header value for `POST /order` |
| `COUPONS_DIR` | `./coupons` | Directory containing coupon files       |

---

## Development Tips

```bash
# Rebuild & recreate with latest changes
docker compose up --build -d

# Force clean rebuild if cache is stale
docker compose build --no-cache api && docker compose up -d --force-recreate

# Logs
docker compose logs -f api
```

---

## Testing (Go)

```bash
go test ./...
```
